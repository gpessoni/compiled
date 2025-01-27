package dto

import "time"

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
	ElementalTypeId  int64     `json:"elemental_type_id"`
	IsHidden         bool      `json:"-"`
	CompanyId        int64     `json:"company_id,omitempty"`
	TableId          *string   `json:"table_id,omitempty"`
	TableOrientation *string   `json:"table_orientation,omitempty"`
	TableIndex       *int64    `json:"table_index,omitempty"`
	AuxId            *string   `json:"-"`
	IsNew            bool      `json:"-"`
	Video            string    `json:"video"`
	Images           string    `json:"images"`
	Price            int64     `json:"price"`
	Tutorial         string    `json:"tutorial"`
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
	TableIndex      int64       `json:"table_index"`
	Url             string      `json:"url"`
	Video           string      `json:"video"`
	Images          string      `json:"images"`
	Price           int64       `json:"price"`
	Tutorial        string      `json:"tutorial"`
}

type CompiledList struct {
	CompiledItems string `json:"compiled_items"`
}
