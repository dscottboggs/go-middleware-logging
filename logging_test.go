package logging

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dscottboggs/attest"
)

//
// func TestMain(m *testing.M) {
// }

func TestLogger(t *testing.T) {
	cfg := DefaultLoggingConfig
	logrec, logwriter := io.Pipe()
	defer func() {
		logrec.Close()
		logwriter.Close()
	}()
	cfg.Writer = logwriter
	InitializeLogger(cfg, "datetime", "method", "endpoint", "identifier")
	test := attest.Test{t}
	patt := `\d+-\d+-\d+ \d+:\d+:\d+.\d+ - GET - /test - testid`
	query := "/test?id=testid"
	urlv := "http://example.com" + query
	request := httptest.NewRequest("GET", urlv, nil)
	writer := httptest.NewRecorder()
	go func() {
		LogRequest(writer, request)
	}()
	logged := make([]byte, len(patt)*2)
	numread, err := logrec.Read(logged)
	test.Handle(err)
	test.GreaterThan(len(patt), numread)
	test.Matches(patt, strings.Trim(string(logged), "\x00\n"))
}
