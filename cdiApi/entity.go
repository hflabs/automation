package apiCdi

type cdiApi struct {
	url      string
	username string
	password string
}

type SearchPartyRequest struct {
	Hid                 int32   `json:"hid,omitempty"`
	Query               string  `json:"query,omitempty"`
	PartyType           string  `json:"partyType,omitempty"`
	MaxCount            int     `json:"maxCount,omitempty"`
	Include             Include `json:"include,omitempty"`
	LastChangeTimeStamp int64   `json:"lastChangeTimestamp,omitempty"`
}

type SearchRelatedPartiesRequest struct {
	FirstPartySearch    SearchPartyRequest `json:"firstPartySearch,omitempty"`
	SecondPartySearch   SearchPartyRequest `json:"secondPartySearch,omitempty"`
	RelationTypes       RelationTypes      `json:"relationTypes,omitempty"`
	Include             Include            `json:"include,omitempty"`
	ReturnSourceParties bool               `json:"lastChangeTimestamp,omitempty"`
}

type FuzzySearchPartyRequest struct {
	Party              Party   `json:"party,omitempty"`
	IncludePartyFields bool    `json:"includePartyFields,omitempty"`
	IncludePartyInfo   Include `json:"include,omitempty"`
}

type FuzzySearchPartyResponse struct {
	MatchedParties []MatchParty `json:"matchedParty,omitempty"`
}

type MatchParty struct {
	MatchRule  int   `json:"matchRule"`
	MatchScope int   `json:"matchScope"`
	Party      Party `json:"party"`
}

type Include struct {
	PartyInfo []string `json:"partyInfo,omitempty"`
}

type Field struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Attribute struct {
	Hid     int32   `json:"hid,omitempty"`
	RawId   string  `json:"rawId,omitempty"`
	Type    string  `json:"type,omitempty"`
	Fields  []Field `json:"field,omitempty"`
	Deleted bool    `json:"deleted,omitempty"`
}

type PartiesResponse struct {
	Party        []Party `json:"party,omitempty"`
	ErrorType    string  `json:"errorType,omitempty"`
	ErrorMessage string  `json:"errorMessage,omitempty"`
}

type PartyResponse struct {
	Party        Party  `json:"party,omitempty"`
	NotModified  bool   `json:"notModified,omitempty"`
	ErrorType    string `json:"errorType,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type Party struct {
	SourceSystem  string      `json:"sourceSystem,omitempty"`
	RawId         string      `json:"rawId,omitempty"`
	Hid           int32       `json:"hid,omitempty"`
	Type          string      `json:"type,omitempty"`
	Fields        []Field     `json:"field,omitempty"`
	Attributes    []Attribute `json:"attribute,omitempty"`
	Relations     []Relation  `json:"relation,omitempty"`
	SourceParties []Source    `json:"source,omitempty"`
}

type Source struct {
	SourceSystem string `json:"sourceSystem,omitempty"`
	RawId        string `json:"rawId,omitempty"`
	Hid          int32  `json:"hid,omitempty"`
}

type RelationTypes struct {
	RelationType []string `json:"relationType,omitempty"`
}

type Relation struct {
	Type       int           `json:"type,omitempty"`
	First      *RelationEdge `json:"first,omitempty"`
	Second     *RelationEdge `json:"second,omitempty"`
	Attributes []Attribute   `json:"attribute,omitempty"`
	Deleted    bool          `json:"deleted,omitempty"`
}

type RelationEdge struct {
	Type         string `json:"type,omitempty"`
	Hid          int32  `json:"hid,omitempty"`
	SourceSystem string `json:"sourceSystem,omitempty"`
	RawId        string `json:"rawId,omitempty"`
}

type RelatedParty struct {
	FirstParty  Party      `json:"firstParty,omitempty"`
	SecondParty Party      `json:"secondParty,omitempty"`
	Relation    []Relation `json:"relation,omitempty"`
}

type RelatedPartyResponse struct {
	Relations []RelatedParty `json:"relatedPartiesEntry,omitempty"`
}

type SaveAndMergeRequest struct {
	Party   []Party `json:"party"`
	Include Include `json:"include"`
}

type SaveRequest struct {
	Party   Party   `json:"party"`
	Include Include `json:"include"`
}

type RelationRequest struct {
	Relation []Relation `json:"relation"`
}

const (
	PhysicalType = "PHYSICAL"
	LegalType    = "LEGAL"

	RawSourceField = "rawSource"

	FullNameField   = "fullNameRawSource"
	SurnameField    = "surname"
	NameField       = "name"
	PatronymicField = "patronymic"
	GenderField     = "gender"

	EmailTypeField  = "type"
	EmailValueField = "email"
	EmailTypeWork   = "WORK"

	PhoneTypeField        = "type"
	PhoneCountryCodeField = "countryCode"
	PhoneCityCodeField    = "cityCode"
	PhoneNumberField      = "number"
	PhoneTypeMobile       = "MOBILE"

	AttributeTypePhone = "PHONE"
	AttributeTypeEmail = "EMAIL"

	LastChangeField = "lastChangeTimestamp"

	RelationAttributeFieldValue = "value"
)
