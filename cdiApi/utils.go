package apiCdi

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

func GetRelationHid(relation Relation) int32 {
	var hid int32
	if relation.First != nil {
		hid = relation.First.Hid
	}
	if relation.Second != nil {
		hid = relation.Second.Hid
	}
	return hid
}
