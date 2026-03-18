package pbdbapi

import (
	"net/url"
	"strconv"
	"strings"
)

type pagination struct {
	Limit  int
	Offset int
}

func applyPagination(v url.Values, p pagination) {
	if p.Limit > 0 {
		v.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Offset > 0 {
		v.Set("offset", strconv.Itoa(p.Offset))
	}
}

func setString(v url.Values, key, val string) {
	if strings.TrimSpace(val) != "" {
		v.Set(key, val)
	}
}

func setInt(v url.Values, key string, val int) {
	if val > 0 {
		v.Set(key, strconv.Itoa(val))
	}
}

func setCSV(v url.Values, key string, vals []string) {
	if len(vals) == 0 {
		return
	}
	joined := strings.Join(vals, ",")
	if strings.TrimSpace(joined) != "" {
		v.Set(key, joined)
	}
}
