package dto

type TutorialStep struct {
	Id          int64  `json:"id,omitempty"`
	PromptId    string `json:"promptId,omitempty"`
	ListId      string `json:"listId,omitempty"`
	Title       string `json:"title"`
	VideoUrl    string `json:"videoUrl,omitempty"`
	Description string `json:"description"`
	OrderIndex  int64  `json:"orderIndex"`
}

type ListMarketplaceInfo struct {
	IsBought     bool           `json:"isBought"`
	IsTipable    bool           `json:"isTipable"`
	NSales       int64          `json:"nSales"`
	TutorialStep []TutorialStep `json:"tutorialStep,omitempty"`
}
