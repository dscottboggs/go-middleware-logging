package logging
/*
# Middleware: Logging

A lightweight, configurable middleware module for async request logging. Simply
initialize the configuration:

    logging.InitializeLogger(logging.DefaultLoggingConfig, "datetime", "method")

and add this line to your middleware function:

	go logging.LogRequest(w, r)

You can also specify a custom LoggingConfig struct, with non-default format
strings or Writer, or override the default, like so

conf := logging.DefaultLoggingConfig
conf.Separator = "\t"

If you would like to add a formatter your contribution will be greatly
appreciated.
*/
import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var availableFormatters map[string]func(
	w http.ResponseWriter,
	r *http.Request,
) string
var enabledFormatters []string

// LoggingConfig contains configurable values for the logger
type LoggingConfig struct {
	DateTimeFormat string
	Separator      string
	Writer         io.WriteCloser
}

// Config is a global value for the LoggingConfig
var Config LoggingConfig

// DefaultLoggingConfig is an object that can be specified as the first argument to
// InitializeLogger in order to set the configuration to the default values.
var DefaultLoggingConfig LoggingConfig

func init() {
	// setup the available formatters.
	DefaultLoggingConfig.DateTimeFormat = "%02d-%02d-%04d %02d:%02d:%02d.%d"
	DefaultLoggingConfig.Separator = " - "
	DefaultLoggingConfig.Writer = os.Stdout
	availableFormatters = make(map[string]func(
		w http.ResponseWriter,
		r *http.Request,
	) string)
	enabledFormatters = make([]string, 4)
	availableFormatters["datetime"] = datetime
	availableFormatters["method"] = requestMethod
	availableFormatters["endpoint"] = endpoint
	availableFormatters["identifier"] = identifier
}

// InitializeLogger accepts an arbitrary number of strings, which correspond to
// the available formatters.
func InitializeLogger(configuration LoggingConfig, formattersToEnable ...string) error {
	Config = configuration
	enabledFormatters = formattersToEnable
	return nil
}

// PrettyListAvailableFormatters returns the list of available logging formatters
// in a human-readable format.
func PrettyListAvailableFormatters() string {
	var out string
	for key := range availableFormatters {
		out += "  " + key + ",\n"
	}
	return out[:len(out)-2] + "\n"
}

// ListAvailableFormatters returns the list of available logging formatters
// in a simple, space-separated list of values.
func ListAvailableFormatters() string {
	var out string
	for key := range availableFormatters {
		out += key + " "
	}
	return out[:len(out)-1]
}

// LogRequest to the configured io.Writer (stdout by default) for each
// enabledFormatter
func LogRequest(w http.ResponseWriter, r *http.Request) {
	var msg string
	for _, key := range enabledFormatters[:len(enabledFormatters)-1] {
		if thisMsg := availableFormatters[key](w, r); thisMsg != "" {
			msg += thisMsg + Config.Separator
		}
	}
	if lastMsg := availableFormatters[enabledFormatters[len(enabledFormatters)-1]](w, r); lastMsg != "" {
		msg += lastMsg + "\n"
	} else {
		msg = msg[:len(msg)-len(Config.Separator)] + "\n"
	}
	Config.Writer.Write(bytes.NewBufferString(msg).Bytes())
}

func datetime(w http.ResponseWriter, r *http.Request) string {
	now := time.Now()
	return fmt.Sprintf(
		Config.DateTimeFormat,
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond())
}

func requestMethod(w http.ResponseWriter, r *http.Request) string {
	return r.Method
}

func endpoint(w http.ResponseWriter, r *http.Request) string {
	return r.URL.Path
}

func identifier(w http.ResponseWriter, r *http.Request) string {
	return r.URL.Query().Get("id")
}
