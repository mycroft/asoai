package chat

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func PatchInput(input string) (string, error) {
	// Regular expression to match ![file <filename>]
	re := regexp.MustCompile(`!\[file\s+([^\]]+)\]`)

	// Keep track of replacement positions to avoid modifying the string
	// while we're iterating over matches
	type replacement struct {
		start, end int
		content    string
	}

	var replacements []replacement

	// Find all matches
	matches := re.FindAllStringSubmatchIndex(input, -1)
	for _, match := range matches {
		if len(match) >= 4 {
			// Extract the filename
			filenameStart, filenameEnd := match[2], match[3]
			filename := strings.TrimSpace(input[filenameStart:filenameEnd])

			// Read the file content
			content, err := os.ReadFile(filename)
			if err != nil {
				return "", fmt.Errorf("error reading file %s: %w", filename, err)
			}

			strContent := fmt.Sprintf("\nContent of file '%s' is:\n```%s\n```", filename, string(content))

			// Store the replacement information
			replacements = append(replacements, replacement{
				start:   match[0],
				end:     match[1],
				content: strContent,
			})
		}
	}

	// Apply replacements from last to first to avoid position shifts
	result := input
	for i := len(replacements) - 1; i >= 0; i-- {
		r := replacements[i]
		result = result[:r.start] + r.content + result[r.end:]
	}

	return result, nil
}
