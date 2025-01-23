package dto

type JSONSubSection struct {
	Id          string           `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Type        string           `json:"type"`
	Content     string           `json:"content"`
	Url         string           `json:"url"`
	Items       []JSONSubSection `json:"items"`
	Video       string           `json:"video"`
	Images      string           `json:"image"`
	Price       int64            `json:"price"`
	Tutorial    string           `json:"tutorial"`
}
