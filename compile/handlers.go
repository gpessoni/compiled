package compile

import (
	"database/sql"
	"fmt"
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
			"compiled_items": parseListResponseAsJSON(root, authUserId, token),
		}
		return response, nil
	} else {
		compiledText := parseListResponse(root, "", authUserId, token, new(int), new(int), format)
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
		"content":     elemental.Template,
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
				"compiled_items": parseListResponseAsJSON(root, authUserId, token),
			}
			return response, nil
		} else {
			compiledText := parseListResponse(root, "", authUserId, token, new(int), new(int), format)
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
		return map[string]interface{}{
			"compiled_items": compiledItems,
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

func parseListResponse(list dto.ListChild, level string, authUserId string, token string, sectionCounter *int, subSectionCounter *int, format string) string {
	var result strings.Builder
	typeName := strings.Title(constants.ElementalConstants.ElementalsArray[list.ElementalTypeId].Name)
	if typeName == "" {
		typeName = strings.Title(constants.ElementalConstants.List.Name)
	}

	levelParts := strings.Split(level, ".")
	indentation := strings.Repeat("  ", len(levelParts))

	headerDepth := len(levelParts) + 1
	if headerDepth > maxHeaderLevel {
		headerDepth = maxHeaderLevel
	}
	headerLevel := strings.Repeat("#", headerDepth)

	if level == "" {
		*sectionCounter++
		level = fmt.Sprintf("%d", *sectionCounter)
		if format == constants.Formats.Markdown {
			result.WriteString(fmt.Sprintf("# %s\n", list.Title))
		} else {
			result.WriteString(fmt.Sprintf("%s\n", list.Title))
		}
		result.WriteString(fmt.Sprintf("%s- Type: %s\n", indentation, typeName))
		result.WriteString(fmt.Sprintf("%s- Description: %s\n\n", indentation, utils.RemoveHTMLTags(list.Description)))
	} else {
		if format == constants.Formats.Markdown {
			result.WriteString(fmt.Sprintf("%s%s %s %s\n", indentation, headerLevel, level, list.Title))
		} else {
			result.WriteString(fmt.Sprintf("%s %s %s\n", indentation, level, list.Title))
		}
		result.WriteString(fmt.Sprintf("%s- Type: %s\n", indentation, typeName))
		result.WriteString(fmt.Sprintf("%s- Description: %s\n\n", indentation, utils.RemoveHTMLTags(list.Description)))
	}

	for i, item := range list.Items {
		itemLevel := fmt.Sprintf("%s.%d", level, i+1)
		childIndentation := strings.Repeat("  ", len(strings.Split(itemLevel, ".")))
		childDepth := len(strings.Split(itemLevel, ".")) + 1
		if childDepth > maxHeaderLevel {
			childDepth = maxHeaderLevel
		}
		childHeaderLevel := strings.Repeat("#", childDepth)

		if !item.IsList && item.ElementalTypeId != constants.ElementalConstants.Table.ID {

			if format == constants.Formats.Markdown {
				result.WriteString(fmt.Sprintf("%s%s %s %s\n", childIndentation, childHeaderLevel, itemLevel, item.Title))
			} else {
				result.WriteString(fmt.Sprintf("%s %s %s\n", childIndentation, itemLevel, item.Title))
			}
			result.WriteString(fmt.Sprintf("%s- Type: %s\n", childIndentation, strings.Title(constants.ElementalConstants.ElementalsArray[item.ElementalTypeId].Name)))
			result.WriteString(fmt.Sprintf("%s- Description: %s\n", childIndentation, utils.RemoveHTMLTags(item.Description)))
			result.WriteString(fmt.Sprintf("%s- Content: %s\n\n", childIndentation, utils.RemoveHTMLTags(item.Template)))
		} else {
			result.WriteString(parseListResponse(item, itemLevel, authUserId, token, sectionCounter, subSectionCounter, format))
		}
	}

	return result.String()
}

func parseListResponseAsJSON(list dto.ListChild, authUserId string, token string) dto.JSONSubSection {
	subSections := []dto.JSONSubSection{}

	for _, item := range list.Items {
		if item.IsList {
			childJSON := parseListResponseAsJSON(item, authUserId, token)
			subSections = append(subSections, dto.JSONSubSection{
				Title:       item.Title,
				Description: utils.RemoveHTMLTags(item.Description),
				Type:        strings.Title(constants.ElementalConstants.List.Name),
				Items:       childJSON.Items,
			})
		} else {
			list := dto.List{
				Id:        item.ListId,
				IsPremium: item.IsPremium,
				UserID:    item.UserId,
			}
			canProceed, err := checkIfListIsPremiumBought(list, authUserId, token)
			if err != nil || !canProceed {
				continue
			}
			childJSON := parseListResponseAsJSON(item, authUserId, token)
			subSections = append(subSections, dto.JSONSubSection{
				Title:       item.Title,
				Description: utils.RemoveHTMLTags(item.Description),
				Type:        strings.Title(constants.ElementalConstants.ElementalsArray[item.ElementalTypeId].Name),
				Content:     utils.RemoveHTMLTags(item.Template),
				Items:       childJSON.Items,
			})
		}
	}

	typeName := strings.Title(constants.ElementalConstants.ElementalsArray[list.ElementalTypeId].Name)
	if typeName == "" {
		typeName = strings.Title(constants.ElementalConstants.List.Name)
	}

	return dto.JSONSubSection{
		Title:       list.Title,
		Type:        typeName,
		Description: utils.RemoveHTMLTags(list.Description),
		Items:       subSections,
	}
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
