package main

import (
	"net/http"

	logging "github.com/dscottboggs/middleware-logging"
)

// DemonstrationHandler -- This is a typical request handler
func DemonstrationHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world\n"))
}

// Middleware -- functions that intercept the request and act upon it, either
// by making some changes to it before passing it on, or like this case,
// asynchronously taking note of the condition under which it was invoked and
// taking some separate action based on that.
func Middleware(
	w http.ResponseWriter,
	r *http.Request,
) (
	http.ResponseWriter,
	*http.Request,
) {
	go logging.LogRequest(w, r)
	return w, r
}

func main() {
	logging.InitializeLogger(
		logging.DefaultLoggingConfig,
		"datetime",
		"method",
		"endpoint",
	)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		DemonstrationHandler(Middleware(w, r))
	})
	http.ListenAndServe(":8080", mux)
}
