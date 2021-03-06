// A Very Nice Web Framework that can be used as a standalone or a embeddable solution for existing Go project that require a Web Interface!
//
// It's built on top of the standard package 'net/http'!
package webby

import (
	"bufio"
	"fmt"
	html "html/template"
	"io"
	"net"
	"net/http"
	"net/http/cgi"
	"net/http/fcgi"
	"os"
	"runtime/debug"
	"time"
)

// Debug Mode
var DEBUG = false

var RootView RouteHandler = BootRoute{Boot, Route}

type webInterface interface {
	http.ResponseWriter
}

type webPrivate struct {
	path       string
	curpath    string
	reswrite   io.Writer
	cut        bool
	firstWrite bool
}

// The Framework Structure, it's implement the interfaces of 'net/http.ResponseWriter',
// 'net/http.Hijacker', 'net/http.Flusher' and 'net/http.Handler'
type Web struct {
	// Error Code
	Status int
	// Server Environment Variables
	Env http.Header
	// Request
	Req *http.Request
	// Meta, useful for storing login credentail
	Meta map[string]interface{}
	// Used by router for storing data of named group in RegExpRule
	Param Param
	// Function to load in html template system.
	HtmlFunc html.FuncMap
	// For holding session!
	Session interface{}
	Errors  *Errors
	webInterface
	pri *webPrivate
}

// HTTP Handler
func (_ Web) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	w := &Web{
		webInterface: res.(webInterface),
		Status:       http.StatusOK,
		Env:          req.Header,
		Req:          req,
		Meta:         map[string]interface{}{},
		Param:        Param{},
		HtmlFunc:     html.FuncMap{},
		Session:      nil,
		Errors: &Errors{
			E403: Error403,
			E404: Error404,
			E500: Error500,
		},
		pri: &webPrivate{
			path:       req.URL.Path,
			curpath:    "",
			cut:        false,
			firstWrite: true,
		},
	}

	w.pri.reswrite = w.webInterface

	w.initTrueHost()
	w.initTrueRemoteAddr()
	w.initTruePath()
	w.initSession()
	w.Header().Set("Content-Encoding", "plain")

	defer func() {
		if r := recover(); r != nil {
			DefaultPanicHandler.Panic(w, r, debug.Stack())
			if DEBUG {
				w.Status = 500
				w.Println("500 Internal Server Error")

				w.Printf("\r\n%s, %s, %s, %s, ?%s IP:%s\r\n",
					w.Req.Proto, w.Req.Method,
					w.Req.Host, w.Req.URL.Path,
					w.Req.URL.RawQuery, w.Req.RemoteAddr)

				w.Printf("\r\n%s\r\n\r\n%s", r, debug.Stack())

				w.Println("\r\nRequest Header:")
				w.Println(w.Req.Header)

				w.ParseForm()

				w.Println("\r\nForm Values:")
				w.Println(w.Req.Form)

				w.Println("\r\nForm Values (Multipart):")
				w.Println(w.Req.MultipartForm)

				w.Println("\r\nTime:")
				w.Println(time.Now())

				return
			}
			w.Error500()
		}
	}()

	defer w.closeCompression()

	w.debugStart()
	defer w.debugEnd()

	HtmlFuncBoot.Load(w)

	if w.CutOut() {
		return
	}

	MainBoot.Load(w)

	if w.CutOut() {
		return
	}

	RootView.View(w)

	if w.CutOut() {
		return
	}

	Error500(w)
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (w *Web) Write(data []byte) (int, error) {
	w.pri.cut = true

	if w.pri.firstWrite {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", http.DetectContentType(data))
		}

		w.pri.firstWrite = false
		w.WriteHeader(w.Status)
	}

	return w.pri.reswrite.Write(data)
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
// Note: Use Status properly to set error code! As this disable compression!
func (w *Web) WriteHeader(num int) {
	w.pri.cut = true

	if w.pri.firstWrite {
		w.pri.firstWrite = false
	}

	w.webInterface.WriteHeader(num)
}

// Hijack lets the caller take over the connection.
// After a call to Hijack(), the HTTP server library
// will not do anything else with the connection.
// It becomes the caller's responsibility to manage
// and close the connection.
func (w *Web) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	w.pri.cut = true

	switch t := w.webInterface.(type) {
	case http.Hijacker:
		return t.Hijack()
	}

	return nil, nil, ErrorStr("Connection is not Hijackable")
}

// Flush sends any buffered data to the client.
func (w *Web) Flush() {
	switch t := w.webInterface.(type) {
	case http.Flusher:
		t.Flush()
	}
}

// Print formats using the default formats for its operands and writes to client (http web server or browser).
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (w *Web) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(w, a...)
}

// Printf formats according to a format specifier and writes to client (http web server or browser).
// It returns the number of bytes written and any write error encountered.
func (w *Web) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(w, format, a...)
}

// Println formats using the default formats for its operands and writes to client (http web server or browser).
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (w *Web) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(w, a...)
}

// true if output was sent to client, otherwise false!
func (w *Web) CutOut() bool {
	return w.pri.cut
}

func (w *Web) debuginfo(a string) {
	if !DEBUG {
		return
	}
	ErrPrintf("--\r\n %s  %s, %s, %s, %s, ?%s IP:%s \r\n--\r\n",
		a, w.Req.Proto, w.Req.Method,
		w.Req.Host, w.Req.URL.Path,
		w.Req.URL.RawQuery, w.Req.RemoteAddr)
}

func (w *Web) debugStart() {
	w.debuginfo("START")
}

func (w *Web) debugEnd() {
	w.debuginfo("END  ")
}

// Start Http Server
func StartHttp(addr string) error {
	return http.ListenAndServe(addr, Web{})
}

// Start Http Server with TLS
func StartHttpTLS(addr string, certFile string, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, Web{})
}

// Start FastCGI Server
func StartFastCGI(l net.Listener) error {
	return fcgi.Serve(l, Web{})
}

// Start CGI, disables Stderr completely. (Due to the way how IIS handlers Stderr)
func StartCGI() error {
	os.Stderr = nil
	return cgi.Serve(Web{})
}
