package compilation

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gpessoni/libgo/constants"
	"github.com/gpessoni/libgo/persistence"
	"github.com/gpessoni/libgo/utils"
	_ "github.com/lib/pq"
)

func GetAllCompiledTextList(db *sql.DB, listId int64, authUserId, token, format, groupBy string) (persistence.CompiledList, error) {
	response, err := prepareListResponse(db, listId, authUserId, token, format, groupBy)
	if err != nil {
		return persistence.CompiledList{}, err
	}
	return response.(persistence.CompiledList), nil
}

func GetAllCompiledTextElemental(db *sql.DB, elementalId string, authUserId, token, format, groupBy string) (persistence.CompiledList, error) {
	response, err := prepareResponseElemental(db, elementalId, authUserId, token, format, groupBy)
	if err != nil {
		return persistence.CompiledList{}, err
	}
	return response.(persistence.CompiledList), nil
}

func GetAllCompiledJsonElemental(db *sql.DB, elementalId string, authUserId, token, format, groupBy string) (map[string]interface{}, error) {
	response, err := prepareResponseElemental(db, elementalId, authUserId, token, format, groupBy)
	if err != nil {
		return nil, err
	}
	return response.(map[string]interface{}), nil
}

func prepareListResponse(db *sql.DB, listId int64, authUserId, token, format, groupBy string) (interface{}, error) {
	adapter := persistence.NewPostgresAdapter(db)
	// list, err := adapter.ListFindByID(listId)
	// if err != nil {
	// 	return "List not found", err
	// }

	// canProceed, err := checkIfListIsPremiumBought(list, authUserId, token)
	// if err != nil || !canProceed {
	// 	return nil, errListPremiumNotBought
	// }

	childs, err := adapter.GetAllCompiledTextList(listId, "")
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
			cells, err := adapter.GetAllCompiledTextList(childs[i].Id, groupBy)
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

	root, ok := utils.Find(childs, func(c persistence.ListChild) bool {
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
		return persistence.CompiledList{CompiledItems: compiledText}, nil
	}
}

func prepareResponseElemental(db *sql.DB, elementalId, authUserId, token, format, groupBy string) (interface{}, error) {
	adapter := persistence.NewPostgresAdapter(db)

	elemental, err := adapter.ElementalFindById(elementalId)
	if err != nil {
		return "Failed to get prompts from list", err
	}

	if elemental.ElementalTypeId == constants.ElementalConstants.Table.ID {
		childs, err := adapter.GetAllCompiledTextList(elementalId, groupBy)
		if err != nil {
			return "Failed to get prompts from list", err
		}

		for i := 0; i < len(childs); i++ {
			for j := i + 1; j < len(childs); j++ {
				if childs[i].Level < childs[j].Level {
					childs[i], childs[j] = childs[j], childs[i]
				}
			}
		}

		root := persistence.ListChild{
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
			return persistence.CompiledList{CompiledItems: compiledText}, nil
		}
	}

	// canProceed, err := checkIfElementalIsPremiumBought(elemental, authUserId, token)
	// if err != nil {
	// 	return nil, err
	// }

	// body := ""
	// if canProceed {
	// 	body = elemental.Template
	// }

	if format == constants.Formats.JSON {
		response := persistence.ElementalJSONResponse{
			Title:       elemental.Title,
			Description: elemental.Description,
			Type:        strings.Title(constants.ElementalConstants.ElementalsArray[elemental.ElementalTypeId].Name),
			Body:        elemental.Template,
		}
		return map[string]interface{}{
			"compiled_items": response,
		}, nil
	} else if format == constants.Formats.Markdown {
		return persistence.CompiledList{
			CompiledItems: fmt.Sprintf("# %s\nType: %s\nDescription: %s\nBody: %s\n",
				elemental.Title,
				strings.Title(constants.ElementalConstants.ElementalsArray[elemental.ElementalTypeId].Name),
				elemental.Description,
				elemental.Template),
		}, nil
	} else {
		return persistence.CompiledList{
			CompiledItems: fmt.Sprintf("Title: %s\nType: %s\nDescription: %s\nBody: %s\n",
				elemental.Title,
				strings.Title(constants.ElementalConstants.ElementalsArray[elemental.ElementalTypeId].Name),
				elemental.Description,
				elemental.Template),
		}, nil
	}
}

func fillItems(parent *persistence.ListChild, childs []persistence.ListChild, authUserId string, token string) {
	fmt.Print(parent, "\n")
	fmt.Print("\n", len(childs), "\n", "\n")
	if parent.Items == nil {
		parent.Items = []persistence.ListChild{}
	}

	visited := make(map[int64]bool)
	fillItemsRecursive(parent, childs, authUserId, token, visited)
}

func fillItemsRecursive(parent *persistence.ListChild, childs []persistence.ListChild, authUserId string, token string, visited map[int64]bool) {
	if visited[parent.LId] {
		return
	}
	visited[parent.LId] = true

	for i := range childs {
		if parent.LId == childs[i].ListId {
			if childs[i].IsList {
				fillItemsRecursive(&childs[i], childs, authUserId, token, visited)
			}

		}
	}
}

const maxHeaderLevel = 6

func parseListResponse(list persistence.ListChild, level string, authUserId string, token string, sectionCounter *int, subSectionCounter *int, format string) string {
	fmt.Print("Aqui lista \n")
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
		fmt.Print(item.ElementalTypeId, "\n")
		fmt.Print(item.IsList, "\n")

		if !item.IsList && item.ElementalTypeId != constants.ElementalConstants.Table.ID {
			fmt.Print("Aqui2")

			if format == constants.Formats.Markdown {
				result.WriteString(fmt.Sprintf("%s%s %s %s\n", childIndentation, childHeaderLevel, itemLevel, item.Title))
			} else {
				result.WriteString(fmt.Sprintf("%s %s %s\n", childIndentation, itemLevel, item.Title))
			}
			result.WriteString(fmt.Sprintf("%s- Type: %s\n", childIndentation, strings.Title(constants.ElementalConstants.ElementalsArray[item.ElementalTypeId].Name)))
			result.WriteString(fmt.Sprintf("%s- Description: %s\n", childIndentation, utils.RemoveHTMLTags(item.Description)))
			result.WriteString(fmt.Sprintf("%s- Body: %s\n\n", childIndentation, utils.RemoveHTMLTags(item.Template)))
		} else {
			result.WriteString(parseListResponse(item, itemLevel, authUserId, token, sectionCounter, subSectionCounter, format))
		}
	}

	return result.String()
}

func parseListResponseAsJSON(list persistence.ListChild, authUserId string, token string) persistence.JSONSubSection {
	subSections := []persistence.JSONSubSection{}

	for _, item := range list.Items {
		if item.IsList {
			childJSON := parseListResponseAsJSON(item, authUserId, token)
			subSections = append(subSections, persistence.JSONSubSection{
				Title:       item.Title,
				Description: utils.RemoveHTMLTags(item.Description),
				Type:        strings.Title(constants.ElementalConstants.List.Name),
				Items:       childJSON.Items,
			})
		} else {
			// list := persistence.List{
			// 	Id:        item.ListId,
			// 	IsPremium: item.IsPremium,
			// 	UserID:    item.UserId,
			// 	canProceed, _ = checkIfListIsPremiumBought(list, authUserId, token)
			// }

			// // if canProceed {
			// // 	parent.Items = append(parent.Items, childs[i])
			// // }
			childJSON := parseListResponseAsJSON(item, authUserId, token)
			subSections = append(subSections, persistence.JSONSubSection{
				Title:       item.Title,
				Description: utils.RemoveHTMLTags(item.Description),
				Type:        strings.Title(constants.ElementalConstants.ElementalsArray[item.ElementalTypeId].Name),
				Body:        utils.RemoveHTMLTags(item.Template),
				Items:       childJSON.Items,
			})
		}
	}

	typeName := strings.Title(constants.ElementalConstants.ElementalsArray[list.ElementalTypeId].Name)
	if typeName == "" {
		typeName = strings.Title(constants.ElementalConstants.List.Name)
	}

	return persistence.JSONSubSection{
		Title:       list.Title,
		Type:        typeName,
		Description: utils.RemoveHTMLTags(list.Description),
		Items:       subSections,
	}
}
