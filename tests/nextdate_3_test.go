package tests

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type nextDate struct {
	date   string
	repeat string
	want   string
}

func TestNextDate(t *testing.T) {
	tbl := []nextDate{
		{"20240126", "", ""},
		{"20240126", "k 34", ""},
		{"20240126", "ooops", ""},
		{"15000156", "y", ""},
		{"ooops", "y", ""},
		{"16890220", "y", `20240220`},
		{"20250701", "y", `20260701`},
		{"20240101", "y", `20250101`},
		{"20231231", "y", `20241231`},
		{"20240229", "y", `20250301`},
		{"20240301", "y", `20250301`},
		{"20240113", "d", ""},
		{"20240113", "d 7", `20240127`},
		{"20240120", "d 20", `20240209`},
		{"20240202", "d 30", `20240303`},
		{"20240320", "d 401", ""},
		{"20231225", "d 12", `20240130`},
		{"20240228", "d 1", "20240229"},
	}
	check := func() {
		for _, v := range tbl {
			urlPath := fmt.Sprintf("api/nextdate?now=20240126&date=%s&repeat=%s",
				url.QueryEscape(v.date), url.QueryEscape(v.repeat))
			get, err := getBody(urlPath)
			assert.NoError(t, err)
			next := strings.TrimSpace(string(get))
			_, err = time.Parse("20060102", next)
			if err != nil && len(v.want) == 0 {
				continue
			}
			assert.Equal(t, v.want, next, `{%q, %q, %q}`,
				v.date, v.repeat, v.want)
		}
	}
	check()
	if !FullNextDate {
		return
	}
}
