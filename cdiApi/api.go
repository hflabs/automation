package apiCdi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/carlmjohnson/requests"
	"io"
	"net/http"
)

func NewCdiApi(url, username, password string) CdiApi {
	return &cdiApi{
		url:      url,
		username: username,
		password: password,
	}
}

var partyInfo []string = []string{"REQUISITE", "ATTRIBUTE", "RELATION", "RELATION_ATTRIBUTE"}

func (c *cdiApi) SearchParty(query, partyType string) ([]Party, error) {
	req := SearchPartyRequest{
		Query:     query,
		PartyType: partyType,
		Include:   Include{partyInfo},
	}
	var parties PartiesResponse
	err := requests.New().Post().
		BaseURL(fmt.Sprintf("%s/soap/services/2_13/PartyRA/search", c.url)).
		BasicAuth(c.username, c.password).
		BodyJSON(req).
		ToJSON(&parties).
		AddValidator(c.validateStatus).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return parties.Party, nil
}

func (c *cdiApi) SearchRelatedParties(firstPartyQuery, firstPartyType, secondPartyQuery, secondPartyType string, relationTypes []string, returnSourceParties bool) ([]RelatedParty, error) {
	req := SearchRelatedPartiesRequest{
		FirstPartySearch: SearchPartyRequest{
			Query:     firstPartyQuery,
			PartyType: firstPartyType,
		},
		SecondPartySearch: SearchPartyRequest{
			Query:     secondPartyQuery,
			PartyType: secondPartyType,
		},
		RelationTypes: RelationTypes{
			RelationType: relationTypes,
		},
		Include:             Include{partyInfo},
		ReturnSourceParties: returnSourceParties,
	}
	var result RelatedPartyResponse
	err := requests.New().Post().
		BaseURL(fmt.Sprintf("%s/soap/services/2_13/PartyRA/searchRelatedParties", c.url)).
		BasicAuth(c.username, c.password).
		BodyJSON(req).
		ToJSON(&result).
		AddValidator(c.validateStatus).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return result.Relations, nil
}

func (c *cdiApi) FuzzySearch(party Party) ([]MatchParty, error) {
	req := FuzzySearchPartyRequest{
		Party:              party,
		IncludePartyFields: true,
		IncludePartyInfo:   Include{partyInfo},
	}
	var result FuzzySearchPartyResponse
	err := requests.New().Post().
		BaseURL(fmt.Sprintf("%s/soap/services/2_13/PartyRA/fuzzySearch", c.url)).
		BasicAuth(c.username, c.password).
		BodyJSON(req).
		ToJSON(&result).
		AddValidator(c.validateStatus).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return result.MatchedParties, nil
}

// GetPartyByHid — метод поиска карточки по её HID. Если указать на входе lastChangeTimestamp,
// то при отсутствии изменений с этой даты отдаст ответ мгновенно, в ответе будет пустой Party и NotModified=true
func (c *cdiApi) GetPartyByHid(hid int32, lastChangeTimestamp int64, partyType string) (Party, bool, error) {
	req := SearchPartyRequest{
		Hid:                 hid,
		PartyType:           partyType,
		LastChangeTimeStamp: lastChangeTimestamp,
		Include:             Include{partyInfo}}
	var party PartyResponse
	err := requests.New().Post().
		BaseURL(fmt.Sprintf("%s/soap/services/2_13/PartyRA/getByHID", c.url)).
		BasicAuth(c.username, c.password).
		BodyJSON(req).
		ToJSON(&party).
		AddValidator(c.validateStatus).
		Fetch(context.Background())
	if err != nil {
		return Party{}, false, err
	}
	if party.NotModified {
		return Party{}, true, nil
	}
	if party.Party.Hid == 0 {
		return Party{}, false, fmt.Errorf("HID not found")
	}
	return party.Party, false, nil
}

func (c *cdiApi) SaveAndMerge(parties []Party) ([]Party, error) {
	var result PartiesResponse
	smr := SaveAndMergeRequest{Party: parties}
	body, err := json.Marshal(smr)
	err = requests.New().Post().
		BaseURL(fmt.Sprintf("%s/soap/services/2_13/PartyRA/saveAndMerge", c.url)).
		BasicAuth(c.username, c.password).
		BodyJSON(SaveAndMergeRequest{Party: parties}).
		ToJSON(&result).
		AddValidator(c.validateStatus).
		Fetch(context.Background())
	if err != nil {
		println(string(body))
		if result.ErrorType != "" {
			return nil, fmt.Errorf("%w:%s:%s", err, result.ErrorType, result.ErrorMessage)
		}
		return nil, err
	}
	return result.Party, nil
}

func (c *cdiApi) Save(party Party) (Party, error) {
	var result PartyResponse
	err := requests.New().Post().
		BaseURL(fmt.Sprintf("%s/soap/services/2_13/PartyRA/save", c.url)).
		BasicAuth(c.username, c.password).
		BodyJSON(SaveRequest{Party: party}).
		ToJSON(&result).
		AddValidator(c.validateStatus).
		Fetch(context.Background())
	if err != nil {
		return Party{}, err
	}
	return result.Party, nil
}

func (c *cdiApi) SaveRelations(relations []Relation) error {
	var result PartyResponse
	err := requests.New().Post().
		BaseURL(fmt.Sprintf("%s/soap/services/2_13/PartyRA/saveRelations", c.url)).
		BasicAuth(c.username, c.password).
		BodyJSON(RelationRequest{Relation: relations}).
		ToJSON(&result).
		AddValidator(c.validateStatus).
		Fetch(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (c *cdiApi) validateStatus(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("status code %v.\nBody:%s", resp.StatusCode, string(b))
}
