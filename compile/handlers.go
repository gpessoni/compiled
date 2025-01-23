package compile

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	ElementalPersistence "github.com/gpessoni/compiled/adapters/persistence/elemental"
	ListPersistence "github.com/gpessoni/compiled/adapters/persistence/list"
	"github.com/gpessoni/compiled/application/constants"
	"github.com/gpessoni/compiled/application/dto"
	"github.com/gpessoni/compiled/application/utils"
)

const maxHeaderLevel = 6

func PrepareListResponse(db *sql.DB, listId int64, authUserId, token, format, groupBy, fields string) (interface{}, error) {
	ListPersistence := ListPersistence.NewListPersistence(db)
	list, err := ListPersistence.FindByID(listId)
	if err != nil {
		return "List not found", err
	}

	canProceed, err := checkIfListIsPremiumBought(list, authUserId, token)
	if err != nil || !canProceed {
		return "You need to buy this list to access it", err
	}

	childs, err := ListPersistence.GetAllCompiledTextList(listId, "")
	if err != nil {
		return "Failed to get prompts from list", err
	}

	for i := 0; i < len(childs); i++ {
		for j := i + 1; j < len(childs); j++ {
			if childs[i].Level < childs[j].Level {
				childs[i], childs[j] = childs[j], childs[i]
			}
		}

		if childs[i].ElementalTypeId == constants.ElementalConstants.Table.ID {
			cells, err := ListPersistence.GetAllCompiledTextList(childs[i].Id, groupBy)
			if err != nil {
				return "Failed to get prompts from list", err
			}

			for k := 0; k < len(cells); k++ {
				for j := k + 1; j < len(cells); j++ {
					if cells[k].Level < cells[j].Level {
						cells[k], cells[j] = cells[j], cells[k]
					}
				}
			}

			fillItems(&childs[i], cells, authUserId, token)

		}
	}

	root, ok := utils.Find(childs, func(c dto.ListChild) bool {
		return c.LId == listId
	})
	if !ok {
		return "Failed to get prompts from list", err
	}

	fillItems(&root, childs, authUserId, token)

	if format == constants.Formats.JSON {
		response := map[string]interface{}{
			"compiled_items": parseListResponseAsJSON(root, authUserId, token, fields),
		}
		return response, nil
	} else {
		compiledText := parseListResponse(root, "", authUserId, token, new(int), new(int), format, fields)
		return dto.CompiledList{CompiledItems: compiledText}, nil
	}
}

func PrepareResponseElemental(db *sql.DB, elementalId, authUserId, token, format, groupBy, fields string) (interface{}, error) {
	elementalPersistence := ElementalPersistence.NewElementalPersistence(db)
	ListPersistence := ListPersistence.NewListPersistence(db)

	elemental, err := elementalPersistence.FindById(elementalId)
	if err != nil {
		return "Failed to get elemental", err
	}

	elementalData := map[string]interface{}{
		"id":          elemental.Id,
		"title":       elemental.Title,
		"type":        strings.Title(constants.ElementalConstants.ElementalsArray[elemental.ElementalTypeId].Name),
		"description": elemental.Description,
		"content":     utils.RemoveHTMLTags(elemental.Template),
		"url":         elemental.Url,
		"is_premium":  elemental.IsPremium,
		"video":       elemental.Video,
		"images":      elemental.Images,
		"price":       elemental.Price,
		"tutorial":    elemental.Tutorial,
	}

	if elemental.ElementalTypeId == constants.ElementalConstants.Table.ID {
		childs, err := ListPersistence.GetAllCompiledTextList(elementalId, groupBy)
		if err != nil {
			return "Failed to get prompts from list", err
		}

		for i := 0; i < len(childs); i++ {
			for j := i + 1; j < len(childs); j++ {
				if childs[i].Level < childs[j].Level || (childs[i].Level == childs[j].Level && childs[i].TableIndex > childs[j].TableIndex) {
					childs[i], childs[j] = childs[j], childs[i]
				}
			}
		}

		root := dto.ListChild{
			Id:              elemental.Id,
			Title:           elemental.Title,
			ElementalTypeId: elemental.ElementalTypeId,
			Description:     elemental.Description,
			Level:           0,
			Template:        elemental.Template,
		}

		fillItems(&root, childs, authUserId, token)

		if format == constants.Formats.JSON {
			response := map[string]interface{}{
				"compiled_items": parseListResponseAsJSON(root, authUserId, token, fields),
			}
			return response, nil
		} else {
			compiledText := parseListResponse(root, "", authUserId, token, new(int), new(int), format, fields)
			return dto.CompiledList{CompiledItems: compiledText}, nil
		}
	}

	canProceed, err := checkIfElementalIsPremiumBought(elemental, authUserId, token)
	if err != nil {
		return nil, err
	}

	selectedFields := strings.Split(fields, ",")
	compiledItems := buildDynamicResponse(selectedFields, elementalData, format, canProceed)

	if format == constants.Formats.JSON {
		parsedCompiledItems := map[string]string{}
		for _, line := range strings.Split(compiledItems, "\n") {
			parts := strings.SplitN(line, ": ", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				parsedCompiledItems[key] = value
			}
		}

		return map[string]interface{}{
			"compiled_items": parsedCompiledItems,
		}, nil
	} else if format == constants.Formats.Markdown {
		return dto.CompiledList{
			CompiledItems: compiledItems,
		}, nil
	} else {
		return dto.CompiledList{
			CompiledItems: compiledItems,
		}, nil
	}
}

func parseListResponse(list dto.ListChild, level string, authUserId string, token string, sectionCounter *int, subSectionCounter *int, format, fields string) string {
	selectedFields := make(map[string]bool)
	for _, field := range strings.Split(fields, ",") {
		selectedFields[strings.TrimSpace(field)] = true
	}

	var result strings.Builder
	typeName := strings.Title(constants.ElementalConstants.ElementalsArray[list.ElementalTypeId].Name)
	if typeName == "" {
		typeName = strings.Title(constants.ElementalConstants.List.Name)
	}

	var id interface{}
	if list.LId != 0 {
		id = list.LId
	} else {
		id = list.Id
	}

	listData := map[string]interface{}{
		"id":          id,
		"title":       list.Title,
		"type":        typeName,
		"description": utils.RemoveHTMLTags(list.Description),
		"content":     utils.RemoveHTMLTags(list.Template),
		"url":         list.Url,
		"video":       list.Video,
		"is_premium":  list.IsPremium,
		"images":      list.Images,
		"price":       list.Price,
		"tutorial":    list.Tutorial,
	}

	var keys []string
	for key := range listData {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		fieldOrder := strings.Split(fields, ",")
		fieldMap := make(map[string]int)
		for idx, field := range fieldOrder {
			fieldMap[strings.TrimSpace(field)] = idx
		}
		return fieldMap[keys[i]] < fieldMap[keys[j]]
	})

	levelParts := strings.Split(level, ".")
	indentation := strings.Repeat("  ", len(levelParts))
	headerDepth := len(levelParts) + 1
	if headerDepth > maxHeaderLevel {
		headerDepth = maxHeaderLevel
	}
	headerLevel := strings.Repeat("#", headerDepth)

	if isFieldSelected("title", selectedFields) {
		if format == constants.Formats.Markdown {
			result.WriteString(fmt.Sprintf("%s%s %s\n", indentation, headerLevel, listData["title"]))
		} else {
			result.WriteString(fmt.Sprintf("%s%s\n", indentation, listData["title"]))
		}
	}

	for _, key := range keys {
		if key != "title" && isFieldSelected(key, selectedFields) {
			result.WriteString(fmt.Sprintf("%s- %s: %v\n", indentation, strings.Title(key), listData[key]))
		}
	}

	for i, item := range list.Items {
		itemLevel := fmt.Sprintf("%s.%d", level, i+1)
		result.WriteString(parseListResponse(item, itemLevel, authUserId, token, sectionCounter, subSectionCounter, format, fields))
	}

	return result.String()
}

func addFieldToSubSection(field string, subSection *dto.JSONSubSection, item dto.ListChild) {
	id := ""
	if item.IsList {
		id = fmt.Sprintf("%d", item.LId)
	} else {
		id = item.Id
	}
	switch field {
	case "id":
		subSection.Id = id
	case "title":
		subSection.Title = item.Title
	case "description":
		subSection.Description = utils.RemoveHTMLTags(item.Description)
	case "type":
		subSection.Type = strings.Title(constants.ElementalConstants.ElementalsArray[item.ElementalTypeId].Name)
	case "content":
		subSection.Content = utils.RemoveHTMLTags(item.Template)
	case "url":
		subSection.Url = item.Url
	case "video":
		subSection.Video = item.Video
	case "images":
		subSection.Images = item.Images
	case "price":
		subSection.Price = item.Price
	case "tutorial":
		subSection.Tutorial = item.Tutorial
	}
}
func parseListResponseAsJSON(list dto.ListChild, authUserId string, token, fields string) dto.JSONSubSection {
	// Criando um mapa para os campos selecionados
	selectedFields := make(map[string]bool)
	for _, field := range strings.Split(fields, ",") {
		selectedFields[strings.TrimSpace(field)] = true
	}

	// Inicializando a lista de sub-seções
	subSections := []dto.JSONSubSection{}

	// Iterando pelos itens da lista
	for _, item := range list.Items {
		// Se for uma lista, processa recursivamente
		if item.IsList {
			childJSON := parseListResponseAsJSON(item, authUserId, token, fields)
			subSections = append(subSections, dto.JSONSubSection{
				Items: childJSON.Items,
			})

			// Adiciona os campos selecionados no item
			for field := range selectedFields {
				addFieldToSubSection(field, &subSections[len(subSections)-1], item)
			}
		} else {
			// Criação da estrutura base para o item
			list := dto.List{
				Id:          item.ListId,
				IsPremium:   item.IsPremium,
				UserID:      item.UserId,
				Title:       item.Title,
				Description: utils.RemoveHTMLTags(item.Description),
				Url:         item.Url,
				Video:       item.Video,
				Images:      item.Images,
				Price:       item.Price,
				Tutorial:    item.Tutorial,
			}

			// Verificando se o usuário pode acessar a lista
			canProceed, err := checkIfListIsPremiumBought(list, authUserId, token)
			if err != nil || !canProceed {
				continue
			}

			// Processamento recursivo do item
			childJSON := parseListResponseAsJSON(item, authUserId, token, fields)
			subSection := dto.JSONSubSection{Items: childJSON.Items}

			// Adicionando os campos selecionados ao item
			for field := range selectedFields {
				addFieldToSubSection(field, &subSection, item)
			}
			subSections = append(subSections, subSection)
		}
	}

	// Criando o retorno final com apenas os campos selecionados
	result := dto.JSONSubSection{}

	// Adicionando os campos selecionados dinamicamente
	for field := range selectedFields {
		addFieldToSubSection(field, &result, list)
	}

	// Adicionando o campo 'items' no final
	result.Items = subSections

	return result
}

func getFieldValue(field string, subSection dto.JSONSubSection) interface{} {
	switch field {
	case "id":
		return subSection.Id
	case "title":
		return subSection.Title
	case "description":
		return subSection.Description
	case "type":
		return subSection.Type
	case "content":
		return subSection.Content
	case "url":
		return subSection.Url
	case "video":
		return subSection.Video
	case "images":
		return subSection.Images
	case "price":
		return subSection.Price
	case "tutorial":
		return subSection.Tutorial
	}
	return nil
}
func buildDynamicResponse(fields []string, data map[string]interface{}, format string, canProceed bool) string {
	var compiledItems string
	var isFirstField bool

	isFirstField = true

	if _, exists := data["content"]; exists && !canProceed {
		data["content"] = ""
	}

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if value, exists := data[field]; exists {
			if format == constants.Formats.Markdown && isFirstField {
				compiledItems += fmt.Sprintf("# %s: %v\n", field, value)
				isFirstField = false
			} else {
				compiledItems += fmt.Sprintf("%s: %v\n", field, value)
			}
		}
	}

	return compiledItems
}

func isFieldSelected(field string, selectedFields map[string]bool) bool {
	return selectedFields[field]
}
