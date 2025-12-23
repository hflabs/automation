package confluence

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

var ErrNotFound = errors.New("confluence: data not found")

func validateStatus(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("status code %v.\nBody:%s", resp.StatusCode, string(b))
}

func extractHashcodeFromContent(content string) string {
	match := regexp.MustCompile(hashcodePattern).FindStringSubmatch(content)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
