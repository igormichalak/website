package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"
	"unicode"
)

var UntrackedPaths = []string{
	"/css",
	"/fonts",
	"/js",
	"/images",
}

const DayInSeconds = 24 * 60 * 60

const ETagCharset = "0123456789abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var RuntimeETag = randomETag()

func randomETag() string {
	buf := make([]byte, 16)
	charsetLen := uint(len(ETagCharset))
	for i := range buf {
		buf[i] = ETagCharset[rand.UintN(charsetLen)]
	}
	return string(buf)
}

func (app *application) cacheHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ccVal := fmt.Sprintf(
			"public, max-age=1800, stale-while-revalidate=%d, stale-if-error=%d",
			365 * DayInSeconds, DayInSeconds,
		)
		w.Header().Set("Cache-Control", ccVal)
		w.Header().Set("ETag", fmt.Sprintf(`"%s"`, RuntimeETag))

		if strings.Contains(r.Header.Get("If-None-Match"), RuntimeETag) {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) securityHeaders(next http.Handler) http.Handler {
	rules := []string{
		"default-src 'self' misc.igormichalak.com",
		"script-src 'self' 'unsafe-eval' cdn.usefathom.com",
		"img-src 'self' cdn.usefathom.com",
		"frame-src 'self'",
		"style-src 'self' 'unsafe-inline'",
		"connect-src 'self' cdn.usefathom.com",
	}
	csp := strings.Join(rules, "; ")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", csp)
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Host, Origin, Referer, Accept, Content-Type, User-Agent, Cookie, X-Csrf-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Add("Vary", "Origin")

		// Preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.error(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) wwwRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Host, "localhost") || unicode.IsDigit([]rune(r.Host)[0]) {
			next.ServeHTTP(w, r)
			return
		}
		if !strings.HasPrefix(r.Host, "www.") {
			dst := "https://www." + r.Host + r.URL.RequestURI()
			http.Redirect(w, r, dst, http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func formatDuration(d time.Duration) string {
	if d.Seconds() >= 1 {
		return fmt.Sprintf("%.2fs", float64(d.Milliseconds())/1e3)
	}
	if d.Milliseconds() >= 1 {
		return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1e3)
	}
	return fmt.Sprintf("%dμs", d.Microseconds())
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range UntrackedPaths {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		start := time.Now()

		next.ServeHTTP(w, r)

		var (
			method  = r.Method
			uri     = r.URL.RequestURI()
			took    = formatDuration(time.Now().Sub(start))
			referer = r.Referer()
		)

		app.Logger.Info("Request", "method", method, "uri", uri, "referer", referer, "took", took)
	})
}
