package payment

import (
	"crypto/md5"
	"encoding/hex"
	"sort"
	"strings"
)

// GenerateSign creates a signature for a given map of parameters and a secret key.
func GenerateSign(params map[string]string, key string) string {
	// Step 1: Sort keys
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Step 2: Concatenate key-value pairs
	var queryParts []string
	for _, k := range keys {
		value := params[k]
		if value != "" {
			queryParts = append(queryParts, k+"="+value)
		}
	}
	queryString := strings.Join(queryParts, "&")

	// Step 3: Append key and hash
	stringToSign := queryString + key
	hasher := md5.New()
	hasher.Write([]byte(stringToSign))
	return hex.EncodeToString(hasher.Sum(nil))
}
