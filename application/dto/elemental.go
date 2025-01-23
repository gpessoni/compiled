package dto

type Elemental struct {
	Id              string
	UserId          string
	Template        string
	Description     string
	Title           string
	IsPremium       bool
	ElementalTypeId int64
	Url             string
	Video           string
	Images          string
	Price           *int64
	Tutorial        string
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
