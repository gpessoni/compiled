package marketplace

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gpessoni/compiled/persistence"
)

var marketplaceApi string

func init() {
	marketplaceApi = os.Getenv("MARKETPLACE_API_URL")
}

func UserHasBoughtList(listId int64, firebaseToken string) (infos persistence.ListMarketplaceInfo, err error) {
	url := marketplaceApi + "/lists/" + fmt.Sprint(listId) + "/information"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return infos, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", firebaseToken))
	req.Header.Add("Sp-App", "Magic")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return infos, err
	}

	if res.StatusCode != http.StatusOK {
		return infos, fmt.Errorf("user has not bought this list yet")
	}

	var data struct {
		Response persistence.ListMarketplaceInfo `json:"response"`
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return infos, err
	}

	if !data.Response.IsBought {
		return data.Response, fmt.Errorf("user has not bought this list yet")
	}

	return data.Response, nil
}

func UserHasBoughtElemental(elementalId string, firebaseToken string) (info persistence.ElementalMarketplaceInfo, err error) {
	url := marketplaceApi + "/prompts/" + elementalId + "/information"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return info, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", firebaseToken))
	req.Header.Add("Sp-App", "Magic")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return info, err
	}

	if res.StatusCode != http.StatusOK {
		return info, fmt.Errorf("user has not bought this elemental yet")
	}

	var data struct {
		Response persistence.ElementalMarketplaceInfo `json:"response"`
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return data.Response, err
	}

	if !data.Response.IsBought {
		return data.Response, fmt.Errorf("user has not bought this elemental yet")
	}

	return data.Response, nil
}
