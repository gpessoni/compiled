package persistence

import "time"

type Elemental struct {
	Id              string
	UserId          string
	Template        string
	Description     string
	Title           string
	IsPremium       bool
	ElementalTypeId int64
}

type List struct {
	CreatedAt        time.Time `json:"-"`
	Id               int64     `json:"id"`
	UserID           string    `json:"user_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Url              string    `json:"url"`
	IsPrivate        bool      `json:"is_private"`
	IsPremium        bool      `json:"is_premium"`
	IsTutorialHidden bool      `json:"is_tutorial_hidden"`
	Avatar           string    `json:"avatar,omitempty"`
	Video            string    `json:"video,omitempty"`
	Price            int64     `json:"price,omitempty"`
	PriceOriginal    *int64    `json:"price_original,omitempty"`
	PriceTypeId      int64     `json:"price_type_id,omitempty"`
	StripeIsProduct  bool      `json:"-"`
	ElementalTypeId  int64     `json:"elemental_type_id"`
	IsHidden         bool      `json:"-"`
	CompaniesSection bool      `json:"-"`
	CompanyId        int64     `json:"company_id,omitempty"`
	TableId          *string   `json:"table_id,omitempty"`
	TableOrientation *string   `json:"table_orientation,omitempty"`
	TableIndex       *int64    `json:"table_index,omitempty"`
	AuxId            *string   `json:"-"`
	IsNew            bool      `json:"-"`
}

type ListChild struct {
	Id              string      `json:"id"`
	UserId          string      `json:"user_id"`
	LId             int64       `json:"lId"`
	IsList          bool        `json:"isList"`
	ListId          int64       `json:"listId"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	Template        string      `json:"template"`
	IsPremium       bool        `json:"is_premium"`
	Items           []ListChild `json:"items"`
	Level           int64       `json:"level"`
	ElementalTypeId int64       `json:"elemental_type_id"`
}

type JSONSubSection struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Type        string           `json:"type"`
	Body        string           `json:"body,omitempty"`
	Items       []JSONSubSection `json:"items,omitempty"`
}

type ElementalJSONResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Body        string `json:"body"`
}

type CompiledList struct {
	CompiledItems string `json:"compiled_items"`
}

type ElementalPersistence interface {
	FindById(id string) (Elemental, error)
}

type ListMarketplaceInfo struct {
	IsBought     bool           `json:"isBought"`
	IsTipable    bool           `json:"isTipable"`
	NSales       int64          `json:"nSales"`
	TutorialStep []TutorialStep `json:"tutorialStep,omitempty"`
}

type TutorialStep struct {
	Id          int64  `json:"id,omitempty"`
	PromptId    string `json:"promptId,omitempty"`
	ListId      string `json:"listId,omitempty"`
	Title       string `json:"title"`
	VideoUrl    string `json:"videoUrl,omitempty"`
	Description string `json:"description"`
	OrderIndex  int64  `json:"orderIndex"`
}

type ElementalMarketplaceInfo struct {
	IsBought     bool           `json:"isBought"`
	IsTipable    bool           `json:"isTipable"`
	Template     string         `json:"template,omitempty"`
	TutorialStep []TutorialStep `json:"tutorialStep,omitempty"`
	NSales       int64          `json:"nSales"`
}
