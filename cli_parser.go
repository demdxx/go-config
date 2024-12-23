package goconfig

import (
	"fmt"
	"strings"
)

// parseCommandFlags parses command line flags
// Examples:
//
//			--http-listen=addr:test --debug
//	   {"http-listen": "addr:test", "debug": "true"}
//
//			-v --http-listen=addr:test --debug --log-level debug
//	   {"v": "true", "http-listen": "addr:test", "debug": "true", "log-level": "debug"}
func parseCommandFlags(args []string) (map[string]string, error) {
	flags := make(map[string]string)

	for i := 0; i < len(args); i++ {
		var (
			arg        = args[i]
			key, value string
		)

		// Check if the argument starts with a flag prefix
		if strings.HasPrefix(arg, "--") {
			arg = strings.TrimPrefix(arg, "--")
		} else if strings.HasPrefix(arg, "-") {
			arg = strings.TrimPrefix(arg, "-")
		} else {
			return flags, fmt.Errorf("invalid flag: %s", arg)
		}

		parts := strings.SplitN(arg, "=", 2)
		key = parts[0]

		if len(parts) == 2 {
			value = parts[1]
		} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
			value = args[i+1]
			i++
		} else {
			value = "true"
		}
		flags[key] = value
	}

	return flags, nil
}
