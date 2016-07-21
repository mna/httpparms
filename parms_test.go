package httpparms

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type parmTest struct {
	S string `schema:"s"`
	I int    `schema:"i"`
	Q string `schema:":q" json:":q"`
}

func (pt *parmTest) Validate() error {
	if pt.I > 2 {
		return errors.New("too big")
	}
	return nil
}

func TestParseQueryForm(t *testing.T) {
	cases := []struct {
		parms url.Values
		want  parmTest
		err   bool
	}{
		{url.Values{}, parmTest{}, false},
		{url.Values{"s": {"a"}}, parmTest{S: "a"}, false},
		{url.Values{"i": {"9"}}, parmTest{I: 9}, true},
		{url.Values{"i": {"1"}}, parmTest{I: 1}, false},
		{url.Values{"i": {"X"}}, parmTest{}, true},
		{url.Values{"s": {"X"}, "i": {"-1"}, ":q": {"q"}}, parmTest{S: "X", I: -1, Q: "q"}, false},
	}
	for i, c := range cases {
		var pt parmTest
		r, err := http.NewRequest("GET", "/a", nil)
		require.NoError(t, err)
		r.Form = c.parms

		got := ParseQueryForm(r, &pt)
		assert.Equal(t, c.err, got != nil, "%d: error expected?", i)
	}
}

func TestParseJSON(t *testing.T) {
	cases := []struct {
		body string
		want parmTest
		err  bool
	}{
		{``, parmTest{}, false},
		{`"abc"`, parmTest{}, true},
		{`{"s": "X"}`, parmTest{S: "X"}, false},
		{`{"i": "X"}`, parmTest{}, true},
		{`{"i": 9}`, parmTest{I: 9}, true},
		{`{"s": "X", "i": 1, ":q": "Q"}`, parmTest{I: 1, S: "X", Q: "Q"}, false},
	}
	for i, c := range cases {
		var pt parmTest
		r, err := http.NewRequest("GET", "/a", strings.NewReader(c.body))
		require.NoError(t, err)

		got := ParseJSON(r, &pt)
		if !assert.Equal(t, c.err, got != nil, "%d: error expected?", i) {
			t.Logf("%d: unexpected error: %v", i, got)
		}
	}
}
