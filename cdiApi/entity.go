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

type GetByHidRequest struct {
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
	SourceSystem string      `json:"sourceSystem,omitempty"`
	RawId        string      `json:"rawId,omitempty"`
	Hid          int32       `json:"hid,omitempty"`
	Type         string      `json:"type,omitempty"`
	Fields       []Field     `json:"field,omitempty"`
	Attributes   []Attribute `json:"attribute,omitempty"`
	Relations    []Relation  `json:"relation,omitempty"`
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

type SaveAndMergeRequest struct {
	Party []Party `json:"party"`
}

type SaveRequest struct {
	Party Party `json:"party"`
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