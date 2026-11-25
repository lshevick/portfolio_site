package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/justinas/nosurf"
)

// Replicates routes.go's dynamic chain (LoadAndSave -> noSurf -> authenticate) but
// with a logging failure handler so we can see exactly which nosurf check fails.
func TestDebugRealChain(t *testing.T) {
	app := newTestApp(t)

	var capturedReason string
	auth := app.authenticate(http.HandlerFunc(app.loginPost)) // authenticate wraps loginPost (alice order)
	csrf := nosurf.New(auth)                                   // noSurf wraps it, with logging failure handler
	csrf.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedReason = nosurf.Reason(r).Error()
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Bad Request")
	}))
	csrf.SetIsTLSFunc(func(_ *http.Request) bool { return false }) // plain HTTP test server
	sessioned := app.sessionManager.LoadAndSave(csrf)          // LoadAndSave outermost

	srv := httptest.NewServer(sessioned)
	defer srv.Close()

	jar, _ := cookiejar.New(nil)
	hc := &http.Client{Jar: jar}

	resp, err := hc.Get(srv.URL + "/user/login")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	re := regexp.MustCompile(`name="csrf_token" value="([^"]+)"`)
	m := re.FindSubmatch(body)

	fmt.Printf("=== GET reason=%q ===\n", capturedReason)

	form := url.Values{}
	if m != nil {
		form.Set("csrf_token", string(m[1]))
	}
	form.Set("email", "x@y.com")
	form.Set("password", "z")

	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/user/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", srv.URL)
	req.Header.Set("Referer", srv.URL+"/user/login")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	resp2, err := hc.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("=== POST status=%d reason=%q ===\n", resp2.StatusCode, capturedReason)
	b2, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()
	fmt.Printf("  body=%q\n", string(b2))

	_ = os.Stdout
}
