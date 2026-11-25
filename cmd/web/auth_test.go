package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/lshevick/internal/models"
	"github.com/lshevick/internal/store"
)

// newTestApp builds an application wired exactly like main.run(), but against a
// throwaway database so tests never touch the real data/app.db.
func newTestApp(t *testing.T) *application {
	t.Helper()
	os.Setenv("DB_PATH", t.TempDir()+"/test_auth.db")

	s, err := store.New(os.Getenv("DB_PATH"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	if err := s.MigrateUp(); err != nil {
		t.Fatalf("MigrateUp: %v", err)
	}

	cache, err := newTemplateCache()
	if err != nil {
		t.Fatalf("newTemplateCache: %v", err)
	}

	sm := scs.New()
	sm.Store = sqlite3store.New(s.DB)
	sm.Lifetime = 24 * time.Hour
	sm.Cookie.Secure = false

	return &application{
		logger:         slog.New(slog.NewTextHandler(os.Stdout, nil)),
		templateCache:  cache,
		sessionManager: sm,
		formDecoder:    form.NewDecoder(),
		store:          s,
		users:          &models.UserModel{DB: s.DB},
	}
}

// extractCSRF pulls the hidden csrf_token value out of a rendered page.
func extractCSRF(t *testing.T, body []byte) string {
	t.Helper()
	re := regexp.MustCompile(`name="csrf_token" value="([^"]+)"`)
	m := re.FindSubmatch(body)
	if m == nil {
		t.Fatalf("no csrf_token found in page")
	}
	return string(m[1])
}

func mustCookieJar() *cookiejar.Jar {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	return jar
}

func mustRead(resp *http.Response) []byte {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return b
}

// postForm builds a browser-like POST (browsers send Origin + Referer + Sec-Fetch-Site for same-origin form posts,
// which nosurf's same-origin check requires). It returns the raw response without following redirects.
func postForm(t *testing.T, hc *http.Client, target string, values url.Values) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, target, strings.NewReader(values.Encode()))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Emulate a real browser submitting a same-origin form.
	req.Header.Set("Origin", target)
	req.Header.Set("Referer", target)
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.ContentLength = int64(len(values.Encode()))
	resp, err := hc.Do(req)
	if err != nil {
		t.Fatalf("%s: %v", target, err)
	}
	return resp
}

func TestLoginPageRenders(t *testing.T) {
	app := newTestApp(t)
	srv := httptest.NewServer(app.routes())
	defer srv.Close()

	hc := &http.Client{Jar: mustCookieJar()}
	resp, err := hc.Get(srv.URL + "/user/login")
	if err != nil {
		t.Fatalf("GET /user/login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	body := mustRead(resp)
	if !strings.Contains(string(body), `name="csrf_token"`) {
		t.Fatal("login page missing csrf_token field")
	}
}

func TestSignupAndLoginFlow(t *testing.T) {
	app := newTestApp(t)
	srv := httptest.NewServer(app.routes())
	defer srv.Close()

	hc := &http.Client{Jar: mustCookieJar()}

	// GET login to obtain a valid CSRF token.
	loginResp, err := hc.Get(srv.URL + "/user/login")
	if err != nil {
		t.Fatalf("GET /user/login: %v", err)
	}
	csrf := extractCSRF(t, mustRead(loginResp))
	loginResp.Body.Close()

	// POST signup with matching passwords -> should redirect to dashboard.
	form := url.Values{}
	form.Set("csrf_token", csrf)
	form.Set("username", "testuser")
	form.Set("email", "test@example.com")
	form.Set("password", "pw123456")
	form.Set("confirm", "pw123456")

	resp := postForm(t, hc, srv.URL+"/user/signup", form)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("signup status = %d, want 303 redirect", resp.StatusCode)
	}

	// Follow to the authenticated dashboard; it must show our username.
	dashResp, err := hc.Get(srv.URL + "/user/dashboard")
	if err != nil {
		t.Fatalf("GET /user/dashboard: %v", err)
	}
	defer dashResp.Body.Close()
	body, _ := io.ReadAll(dashResp.Body)
	if !strings.Contains(string(body), "testuser") {
		t.Fatalf("dashboard did not show username; body=%s", string(body))
	}

	// The same session must now pass the protected route directly (no redirect).
	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/user/dashboard", nil)
	hc2 := &http.Client{Jar: mustCookieJar()}
	hc2.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp2, err := hc2.Do(req)
	if err != nil {
		t.Fatalf("protected GET: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != 200 {
		t.Fatalf("authenticated dashboard status = %d, want 200", resp2.StatusCode)
	}
}

func TestLoginRejectsBadCredentials(t *testing.T) {
	app := newTestApp(t)
	srv := httptest.NewServer(app.routes())
	defer srv.Close()

	hc := &http.Client{Jar: mustCookieJar()}
	loginResp, err := hc.Get(srv.URL + "/user/login")
	if err != nil {
		t.Fatalf("GET /user/login: %v", err)
	}
	csrf := extractCSRF(t, mustRead(loginResp))
	loginResp.Body.Close()

	form := url.Values{}
	form.Set("csrf_token", csrf)
	form.Set("email", "nobody@example.com")
	form.Set("password", "wrong")

	resp := postForm(t, hc, srv.URL+"/user/login", form)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("bad-login status = %d, want 200 (re-render)", resp.StatusCode)
	}
	body := mustRead(resp)
	if !strings.Contains(string(body), "invalid credentials") {
		t.Fatalf("expected invalid-credentials flash in body; got=%s", string(body))
	}
}

func TestSignupDuplicateEmailRejected(t *testing.T) {
	app := newTestApp(t)
	srv := httptest.NewServer(app.routes())
	defer srv.Close()

	hc := &http.Client{Jar: mustCookieJar()}
	loginResp, err := hc.Get(srv.URL + "/user/login")
	if err != nil {
		t.Fatalf("GET /user/login: %v", err)
	}
	csrf := extractCSRF(t, mustRead(loginResp))
	loginResp.Body.Close()

	form := url.Values{}
	form.Set("csrf_token", csrf)
	form.Set("username", "dup")
	form.Set("email", "dupe@example.com")
	form.Set("password", "pw123456")
	form.Set("confirm", "pw123456")

	// First signup succeeds.
	if resp := postForm(t, hc, srv.URL+"/user/signup", form); resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("first signup status = %d, want 303 redirect; body=%s", resp.StatusCode, mustRead(resp))
	}

	// Second with same email must be rejected, not create a duplicate.
	resp := postForm(t, hc, srv.URL+"/user/signup", form)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("duplicate signup status = %d, want 200 (re-render)", resp.StatusCode)
	}
	body := mustRead(resp)
	if !strings.Contains(string(body), "already registered") {
		t.Fatalf("expected duplicate-email flash; got=%s", string(body))
	}
}

func TestSignupMismatchedPasswordsRejected(t *testing.T) {
	app := newTestApp(t)
	srv := httptest.NewServer(app.routes())
	defer srv.Close()

	hc := &http.Client{Jar: mustCookieJar()}
	loginResp, err := hc.Get(srv.URL + "/user/login")
	if err != nil {
		t.Fatalf("GET /user/login: %v", err)
	}
	csrf := extractCSRF(t, mustRead(loginResp))
	loginResp.Body.Close()

	form := url.Values{}
	form.Set("csrf_token", csrf)
	form.Set("username", "mismatch")
	form.Set("email", "mm@example.com")
	form.Set("password", "abc12345")
	form.Set("confirm", "totally-different")

	resp := postForm(t, hc, srv.URL+"/user/signup", form)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("mismatch status = %d, want 200 (re-render)", resp.StatusCode)
	}
	body := mustRead(resp)
	if !strings.Contains(string(body), "Passwords do not match") {
		t.Fatalf("expected password-mismatch flash; got=%s", string(body))
	}
}

func TestDashboardRequiresAuthentication(t *testing.T) {
	app := newTestApp(t)
	srv := httptest.NewServer(app.routes())
	defer srv.Close()

	// A client with no session cookie must be redirected to login (not followed).
	hc := &http.Client{Jar: mustCookieJar()}
	hc.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := hc.Get(srv.URL + "/user/dashboard")
	if err != nil {
		t.Fatalf("GET /user/dashboard (unauth): %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("unauthenticated dashboard status = %d, want 303 redirect to login", resp.StatusCode)
	}
	loc := resp.Header.Get("Location")
	if !strings.HasSuffix(loc, "/user/login") {
		t.Fatalf("redirect Location = %q, want /user/login suffix", loc)
	}
}
