package confluence

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestExtractHashcodeFromContent(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		content, err := os.ReadFile("test_data/content_with_hashcode.txt")
		require.NoError(t, err)
		details := extractHashcodeFromContent(string(content))
		require.Equal(t, "304c51c4ef838285f89d5be131697258", details)
	})
}
