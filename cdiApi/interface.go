package apiCdi

type CdiApi interface {
	SearchParty(query, partyType string) ([]Party, error)
	SearchRelatedParties(firstPartyQuery, firstPartyType, secondPartyQuery, secondPartyType string, relationTypes []string, returnSourceParties bool) ([]RelatedParty, error)
	FuzzySearch(party Party) ([]MatchParty, error)

	GetPartyByHid(hid int32, lastChangeTimestamp int64, partyType string) (Party, bool, error)

	SaveAndMerge(parties []Party) ([]Party, error)
	Save(party Party) (Party, error)
	SaveRelations(relations []Relation) error

	CloseAttribute(partyType, attributeType string, attributeHid int32) error
}
