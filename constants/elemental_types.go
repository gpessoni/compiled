package constants

type elementalContants struct {
	List            ElementalContantStruct
	Prompt          ElementalContantStruct
	Snippet         ElementalContantStruct
	Document        ElementalContantStruct
	Output          ElementalContantStruct
	Automation      ElementalContantStruct
	Course          ElementalContantStruct
	Table           ElementalContantStruct
	Cell            ElementalContantStruct
	ElementalsArray []ElementalContantStruct
}

type Functionalities struct {
	IsTemplateHidden        bool `json:"isTemplateHidden"`
	IsTitleLowercase        bool `json:"isTitleLowercase"`
	HasTemplatePlaceholders bool `json:"hasTemplatePlaceholders"`
	HasTemplateStyles       bool `json:"hasTemplateStyles"`
	HasTemplateAllStyles    bool `json:"hasTemplateAllStyles"`
	HasTypeId               bool `json:"hasTypeId"`
	CanUse                  bool `json:"canUse"`
}

type Validation struct {
	CanDuplicatedTitle      bool     `json:"duplicateTitle"`
	LengthTitleMin          int      `json:"LengthTitleMin"`
	LengthTitleMax          int      `json:"LengthTitleMax"`
	LengthDescriptionMin    int      `json:"LengthDescriptionMin"`
	LengthDescriptionMax    int      `json:"LengthDescriptionMax"`
	LengthTemplateMin       int      `json:"LengthTemplateMin"`
	LengthTemplateMax       int      `json:"LengthTemplateMax"`
	LengthTopicsMin         int      `json:"LengthTopicsMin"`
	LengthTopicsMax         int      `json:"LengthTopicsMax"`
	PermittedTypeIds        []int    `json:"permittedTypeIds"`
	PermittedFileExtensions []string `json:"permittedFileExtensions"`
}

type Layout struct {
	Title                 string `json:"title"`
	Avatar                string `json:"avatar"`
	SaveToList            string `json:"saveToList"`
	Template              string `json:"template"`
	Attachment            string `json:"attachment"`
	IsPrivate             string `json:"isPrivate"`
	Comments              string `json:"comments"`
	Reviews               string `json:"reviews"`
	IsPremium             string `json:"isPremium"`
	Type                  string `json:"type"`
	TopicIds              string `json:"topic_ids"`
	CommandSlash          string `json:"command_slash"`
	Source                string `json:"source"`
	ImagesCover           string `json:"images_cover"`
	Description           string `json:"description"`
	CommandKey            string `json:"command_key"`
	Images                string `json:"images"`
	VideoUrl              string `json:"video_url"`
	HowtoSection          string `json:"howto_section"`
	FeaturedReview        string `json:"featured_review"`
	CompanySection        string `json:"company_section"`
	FeaturesSection       string `json:"features_section"`
	TestimonialSection    string `json:"testimonial_section"`
	TweetsSection         string `json:"tweets_section"`
	AdditionalTextSection string `json:"addicitional_text_section"`
}

type Metadata struct {
	IsActive        bool            `json:"isActive"`
	IsSearchable    bool            `json:"isSearchable"`
	Functionalities Functionalities `json:"functionalities"`
	Validation      Validation      `json:"validation"`
	Layout          Layout          `json:"layout"`
}
type ElementalContantStruct struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Url      string   `json:"url"`
	IsActive bool     `json:"-"`
	Metadata Metadata `json:"metadata"`
}
