package urlshort

import (
	yaml "gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get request's URL, find it in paths to url, if not, use fallback.
		path := r.URL.Path
		if newPath, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, newPath, http.StatusMovedPermanently)
			return
		}
		if fallback != nil {
			fallback.ServeHTTP(w, r)
		}
	}
}

type pathUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func parseYAML(input []byte) ([]pathUrl, error) {
	var p []pathUrl
	err := yaml.UnmarshalStrict(input, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func buildMap(parsedYaml []pathUrl) map[string]string {
	pathMap := make(map[string]string, len(parsedYaml))
	for _, p := range parsedYaml {
		if p.Path == "" || p.Url == "" {
			continue
		}
		pathMap[p.Path] = p.Url
	}
	return pathMap
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(input []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(input)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}
