package compile

import "github.com/gpessoni/compiled/application/dto"

func fillItems(parent *dto.ListChild, childs []dto.ListChild, authUserId string, token string) {
	if parent.Items == nil {
		parent.Items = []dto.ListChild{}
	}

	visited := make(map[int64]bool)
	fillItemsRecursive(parent, childs, authUserId, token, visited)
}

func fillItemsRecursive(parent *dto.ListChild, childs []dto.ListChild, authUserId string, token string, visited map[int64]bool) {
	if visited[parent.LId] {
		return
	}
	visited[parent.LId] = true

	for i := range childs {
		if parent.LId == childs[i].ListId {
			canProceed := true
			if childs[i].IsList {
				fillItemsRecursive(&childs[i], childs, authUserId, token, visited)
			}

			if childs[i].IsPremium {
				list := dto.List{
					Id:        childs[i].LId,
					IsPremium: childs[i].IsPremium,
					UserID:    childs[i].UserId,
				}
				canProceed, _ = checkIfListIsPremiumBought(list, authUserId, token)
			}

			if canProceed {
				parent.Items = append(parent.Items, childs[i])
			}
		}
	}
}
