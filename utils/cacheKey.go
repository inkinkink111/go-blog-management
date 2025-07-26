package utils

import (
	"fmt"
	"sort"
	"strings"
)

func GenerateCacheKey(page, limit string, tags []string) string {
	var keyParts []string
	// Base key
	keyParts = append(keyParts, "blog:list")
	// Add pagination
	keyParts = append(keyParts, fmt.Sprintf("page:%s", page))
	keyParts = append(keyParts, fmt.Sprintf("limit:%s", limit))
	// Add tags
	if len(tags) > 0 {
		// Sort tags for consistent cache keys
		sortedTags := make([]string, len(tags))
		copy(sortedTags, tags)
		sort.Strings(sortedTags)

		// Join tags with a delimiter
		tagsString := strings.Join(sortedTags, ",")
		keyParts = append(keyParts, fmt.Sprintf("tags:%s", tagsString))
	}

	// Join all parts
	return strings.Join(keyParts, ":")
}
