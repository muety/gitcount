package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// ReadMailMap reads a .mailmap file and returns the relation mapped amongst the emails.
func ReadMailMap(reader io.Reader) (map[string]*MailMapEntry, error) {
	if reader == nil {
		return nil, nil
	}

	mailmap := make(map[string]*MailMapEntry)
	scanner := bufio.NewScanner(reader)
	utf8bom := []byte{0xEF, 0xBB, 0xBF} //UTF8 Byte order mark

	// Regular Expression to extract the fields required. Only uses proper email and commit email
	pattern := regexp.MustCompile(`^(?P<pname>[^<]+)?(?:\\s+)?(?:<(?P<pmail>.*?)>)?(?:\\s+)?(?P<cname>[^<]+)?(?:\\s+)?(?:<(?P<cmail>.*?)>)?$`)

	currentLine := 0
	for scanner.Scan() {
		scannedBytes := scanner.Bytes()
		if currentLine == 0 {
			scannedBytes = bytes.TrimPrefix(scannedBytes, utf8bom)
		}
		currentLine++
		row := string(scannedBytes)
		if strings.HasPrefix(row, "#") {
			continue
		}
		if row := strings.TrimSpace(row); row == "" {
			continue
		}

		groups := pattern.FindStringSubmatch(row)
		entry := MailMapEntry{groups[2], groups[4]}

		if _, hasKey := mailmap[entry.CommitEmail]; !hasKey {
			mailmap[entry.CommitEmail] = &entry
		}
		if _, hasKey := mailmap[entry.ProperEmail]; !hasKey {
			mailmap[entry.ProperEmail] = &entry
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading .mailmap: %v", err)
	}
	return mailmap, nil
}
