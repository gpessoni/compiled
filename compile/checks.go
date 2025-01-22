package compile

import (
	"github.com/gpessoni/compiled/adapters/marketplace"
	"github.com/gpessoni/compiled/application/dto"
)

func checkIfListIsPremiumBought(list dto.List, authUserId, token string) (bool, error) {

	if !list.IsPremium || list.UserID == authUserId {
		return true, nil
	}

	infos, err := marketplace.UserHasBoughtList(list.Id, token)
	if err != nil || !infos.IsBought {
		return false, err
	}

	return true, nil
}

func checkIfElementalIsPremiumBought(element dto.Elemental, authUserId, token string) (bool, error) {
	if element.UserId == authUserId {
		return true, nil
	}

	if !element.IsPremium {
		return true, nil
	}

	infos, err := marketplace.UserHasBoughtElemental(element.Id, token)
	if err != nil || !infos.IsBought {
		return false, nil
	}

	return true, nil
}
