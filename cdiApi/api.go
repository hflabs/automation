package apiCdi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/carlmjohnson/requests"
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
		Fetch(context.Background())
	if err != nil {
		if parties.ErrorType != "" {
			return nil, fmt.Errorf("%w:%s:%s", err, parties.ErrorType, parties.ErrorMessage)
		}
		return nil, err
	}
	return parties.Party, nil
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
	body, err := json.Marshal(req)
	if err != nil {
		return Party{}, false, err
	}
	err = requests.New().Post().
		BaseURL(fmt.Sprintf("%s/soap/services/2_13/PartyRA/getByHID", c.url)).
		BasicAuth(c.username, c.password).
		BodyJSON(req).
		ToJSON(&party).
		Fetch(context.Background())
	if err != nil {
		println(string(body))
		if party.ErrorType != "" {
			return Party{}, false, fmt.Errorf("%w:%s:%s", err, party.ErrorType, party.ErrorMessage)
		}
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
		Fetch(context.Background())
	if err != nil {
		if result.ErrorType != "" {
			return Party{}, fmt.Errorf("%w:%s:%s", err, result.ErrorType, result.ErrorMessage)
		}
		return Party{}, err
	}
	return result.Party, nil
}

func (c *cdiApi) SaveRelations(relations []Relation) error {
	body, err := json.Marshal(RelationRequest{Relation: relations})
	var result PartyResponse
	err = requests.New().Post().
		BaseURL(fmt.Sprintf("%s/soap/services/2_13/PartyRA/saveRelations", c.url)).
		BasicAuth(c.username, c.password).
		BodyJSON(RelationRequest{Relation: relations}).
		ToJSON(&result).
		Fetch(context.Background())
	if err != nil {
		println(string(body))
		if result.ErrorType != "" {
			return fmt.Errorf("%w:%s:%s", err, result.ErrorType, result.ErrorMessage)
		}
		return err
	}
	return nil
}
