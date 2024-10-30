package apiCdi

type CdiApi interface {
	SearchParty(query, partyType string) ([]Party, error)
	GetPartyByHid(hid int32, lastChangeTimestamp int64, partyType string) (Party, bool, error)
	SaveAndMerge(parties []Party) ([]Party, error)
	Save(party Party) (Party, error)
	SaveRelations(relations []Relation) error
}
