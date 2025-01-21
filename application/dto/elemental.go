package dto

type Elemental struct {
	Id              string
	UserId          string
	Template        string
	Description     string
	Title           string
	IsPremium       bool
	ElementalTypeId int64
}

type ElementalJSONResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Content     string `json:"content"`
}

type ElementalMarketplaceInfo struct {
	IsBought     bool           `json:"isBought"`
	IsTipable    bool           `json:"isTipable"`
	Template     string         `json:"template,omitempty"`
	TutorialStep []TutorialStep `json:"tutorialStep,omitempty"`
	NSales       int64          `json:"nSales"`
}

type ElementalPersistence interface {
	FindById(id string) (Elemental, error)
}
