package artifactory

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestBuildQuery(t *testing.T) {
	repo := "my-repo"
	folderPath := "/my/path"
	tests := []struct {
		name     string
		options  ListOptions
		expected string
	}{
		{
			name:     "basic query",
			options:  ListOptions{},
			expected: `items.find({"repo":"my-repo","path":"/my/path","type":"any"})`,
		},
		{
			name: "query with name (exact match)",
			options: ListOptions{
				Name: "my-file.txt",
			},
			expected: `items.find({"repo":"my-repo","path":"/my/path","type":"any","name":"my-file.txt"})`,
		},
		{
			name: "query with name (wildcard match)",
			options: ListOptions{
				Name: "my-file-*.txt",
			},
			expected: `items.find({"repo":"my-repo","path":"/my/path","type":"any","name":{"$match":"my-file-*.txt"}})`,
		},
		{
			name: "query with sorting (asc)",
			options: ListOptions{
				SortValue:     SortByName,
				SortDirection: SortTypeAsc,
			},
			expected: `items.find({"repo":"my-repo","path":"/my/path","type":"any"})` + `.sort({"$name":["asc"]})`,
		},
		{
			name: "query with sorting (desc, default)",
			options: ListOptions{
				SortValue: SortByCreated,
			},
			expected: `items.find({"repo":"my-repo","path":"/my/path","type":"any"})` + `.sort({"$created":["desc"]})`,
		},
		{
			name: "query with limit",
			options: ListOptions{
				Limit: 10,
			},
			expected: `items.find({"repo":"my-repo","path":"/my/path","type":"any"})` + `.limit(10)`,
		},
		{
			name: "query with all options",
			options: ListOptions{
				Name:          "my-file.txt",
				SortValue:     SortBySize,
				SortDirection: SortTypeAsc,
				Limit:         5,
			},
			expected: `items.find({"repo":"my-repo","path":"/my/path","type":"any","name":"my-file.txt"})` +
				`.sort({"$size":["asc"]}).limit(5)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := buildFindQuery(repo, folderPath, tt.options)
			result = strings.Replace(result, includeTemplate, "", 1)
			result = strings.ReplaceAll(result, "\n", "")
			result = strings.ReplaceAll(result, " ", "")
			require.NoError(t, err, "unexpected error")
			require.Equal(t, tt.expected, result, "query mismatch")
		})
	}
}
