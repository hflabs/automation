package artifactory

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
)

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

func buildFindQuery(repo, folderPath string, options ListOptions) (string, error) {
	query := findQuery{
		Repo: repo,
		Path: path.Clean(folderPath),
		Type: cmp.Or(options.Type, TypeAny),
	}
	if options.Name != "" {
		query.Name = &name{Value: options.Name}
	}

	jsQuery, err := json.Marshal(query)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}
	result := fmt.Sprintf(findTemplate, jsQuery) + includeTemplate
	if options.SortValue != "" {
		result += fmt.Sprintf(sortTemplate, options.SortValue, cmp.Or(options.SortDirection, SortTypeDesc))
	}
	if options.Limit > 0 {
		result += fmt.Sprintf(limitTemplate, options.Limit)
	}
	return result, nil
}
