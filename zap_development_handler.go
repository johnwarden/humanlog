package humanlog

import (
	"regexp"
	"strings"
	"time"
)

// Zap Development logs are made up of the following separated by whitespace
// 1. timestamp in ISO-8601 (??)
// 2. Log Level (one of DEBUG ERROR INFO WARN FATAL)
// 3. Caller Location in the source
// 4. The main logged message
// 5. a JSON object containing the structured k/v pairs
// 6. optional context lines - but since they are on a separate line the main
//    scanner loop will never capture them
var zapDevLogsPrefixRe = regexp.MustCompile("^(?P<timestamp>\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}-\\d{4})\\s+(?P<level>\\w{4,5})\\s+(?P<location>\\S+)\\s+(?P<message>[^{]+?)\\s+(?P<jsonbody>{.+})$")

// Zap Development Logs when run in Docker-Compose are nearly identical to before
// Fields are tab separated instead of whitespace
// Timestamp is now in ...
// Everything else remains the same
var zapDevDCLogsPrefixRe = regexp.MustCompile("^(?P<timestamp>\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}Z)\t(?P<level>\\w{4,5})\t(?P<location>\\S+)\t(?P<message>[^{]+?)\t(?P<jsonbody>{.+})$")

// This is not obviously an RFC-compliant format and is not a constant in the
// time package which is worrisome but this pattern does work.
const someRFC = "2006-01-02T15:04:05.000-0700"

func tryZapDevPrefix(d []byte, handler *JSONHandler) bool {
	if matches := zapDevLogsPrefixRe.FindSubmatch(d); matches != nil {
		if handler.TryHandle(matches[5]) {
			t, err := time.Parse(someRFC, string(matches[1]))
			if err != nil {
				return false
			}
			handler.Time = t

			handler.Level = strings.ToLower(string(matches[2]))
			handler.setField([]byte("caller"), matches[3])
			handler.Message = string(matches[4])
			return true
		}
	}
	return false
}
