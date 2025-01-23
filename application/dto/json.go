package dto

type JSONSubSection struct {
	Id          string           `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Type        string           `json:"type"`
	Content     string           `json:"content,omitempty"`
	Url         string           `json:"url"`
	Items       []JSONSubSection `json:"items,omitempty"`
	IsPremium   bool             `json:"isPremium"`
	Video       string           `json:"video"`
	Images      string           `json:"image"`
	Price       int64            `json:"price"`
	Tutorial    string           `json:"tutorial"`
}
