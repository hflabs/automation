package confluence

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
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

func extractHashcodeFromContent(content string) string {
	match := regexp.MustCompile(hashcode_pattern).FindStringSubmatch(content)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
