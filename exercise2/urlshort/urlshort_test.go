package urlshort

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testString string = "Test string"

func testMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, testString)
		})
	return mux
}

var pathsToUrls map[string]string = map[string]string{"/foo": "/bar", "/foo2": "/bar2"}

var yamlMap string = `
- path: /foo
  url: /bar
- path: /foo2
  url: /bar2
- path:
  url: /bar3
- path: /foo3
  url:
`

func checkPathRedirects(t *testing.T, handler http.HandlerFunc, pathsToUrls map[string]string) {
	t.Helper()
	for oldPath, newPath := range pathsToUrls {
		request, _ := http.NewRequest(http.MethodGet, oldPath, nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		response := recorder.Result()
		// Check HTTP Redirect: Response status code and location header.
		assertResponseCode(t, response.StatusCode, http.StatusMovedPermanently)
		assertResponseLocation(t, response.Header.Get("location"), newPath)
	}
}

func TestMapHandler(t *testing.T) {
	handler := MapHandler(pathsToUrls, testMux())

	t.Run("Test URL Redirects", func(t *testing.T) {
		checkPathRedirects(t, handler, pathsToUrls)
	})

	t.Run("Test Fallback", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/inexisting-path", nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		response := recorder.Result()
		assertResponseCode(t, response.StatusCode, 200)
		assertResponseLocation(t, response.Header.Get("location"), "")
	})
}

func TestYAMLHandler(t *testing.T) {
	handler, err := YAMLHandler([]byte(yamlMap), testMux())
	assertNoError(t, err)

	t.Run("Test YAML URL Redirects", func(t *testing.T) {
		checkPathRedirects(t, handler, pathsToUrls)
	})

	t.Run("Test YAML Fallback", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/foo3", nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		response := recorder.Result()
		assertResponseCode(t, response.StatusCode, 200)
		assertResponseLocation(t, response.Header.Get("location"), "")
	})

	t.Run("Test erroneous YAML", func(t *testing.T) {
		_, err := YAMLHandler([]byte(`erlkjerjl`), testMux())
		assertError(t, err)
	})
}

func assertError(t *testing.T, got error) {
	t.Helper()
	if got == nil {
		t.Fatalf("expected error, got '%v'", got)
	}
}

func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("expected no error, got '%v'", got)
	}
}

func assertResponseCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("want '%d' response code, got '%d'", want, got)
	}
}

func assertResponseLocation(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("want '%s' response location, got '%s'", want, got)
	}
}
