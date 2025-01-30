package apiCdi

import (
	"fmt"
	"io"
	"net/http"
)

func ParseFields(fields []Field) map[string]string {
	result := make(map[string]string)
	for _, field := range fields {
		result[field.Name] = field.Value
	}
	return result
}

func GetFieldValue(fields map[string]string, field string) string {
	value, ok := fields[field]
	if ok {
		return value
	}
	return ""
}

func GetRelationHid(relation Relation) (int32, string) {
	var hid int32
	var partyType string
	if relation.First != nil {
		hid = relation.First.Hid
		partyType = relation.First.Type
	}
	if relation.Second != nil {
		hid = relation.Second.Hid
		partyType = relation.Second.Type
	}
	return hid, partyType
}

func validateStatus(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("status code %v.\nBody:%s", resp.StatusCode, string(b))
}
