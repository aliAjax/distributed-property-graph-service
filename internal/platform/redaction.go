package platform

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func Digest(value string) string {
	h := sha256.Sum256([]byte(value))
	return hex.EncodeToString(h[:12])
}
func Redact(values map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range values {
		lower := strings.ToLower(k)
		if strings.Contains(lower, "secret") || strings.Contains(lower, "token") || strings.Contains(lower, "password") {
			out[k] = "[REDACTED]"
		} else {
			out[k] = v
		}
	}
	return out
}
func MaskID(value string) string {
	if len(value) <= 4 {
		return "****"
	}
	return value[:2] + "****" + value[len(value)-2:]
}
