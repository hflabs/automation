package confluence

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestExtractHashcodeFromContent(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		hash     string
	}{
		{
			name:     "01. hashcode at page start",
			filePath: "test_data/content_with_hashcode_start.txt",
			hash:     "7a8b9ecc2a62cc944862a408a913e913",
		},
		{
			name:     "02. hashcode after similar macros",
			filePath: "test_data/content_with_hashcode.txt",
			hash:     "7a8b9ecc2a62cc944862a408a913e914",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := os.ReadFile(tt.filePath)
			require.NoError(t, err)
			resultHash := extractHashcodeFromContent(string(content))
			require.Equal(t, tt.hash, resultHash)
		})
	}

}
