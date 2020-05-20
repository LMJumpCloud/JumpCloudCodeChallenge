package routing

import (
	"net/http"
	"regexp"
	"strings"
)

var pathParamRegex = regexp.MustCompile(`\{(.+?)\}`)

// ParameterizedPath holds information regarding how to pull data from a parameterized URL path
type ParameterizedPath struct {
	Path string
	Length int
	Route map[int]string // map of path index -> path segment
	Subs map[int]string // map of path index -> parameter name
}

func (p *ParameterizedPath) matches(path []string) bool {
	if p.Length != len(path) {
		return false
	}

	for i, segment := range p.Route {
		if path[i] != segment {
			return false
		}
	}
	return true
}

// ParseRequest first checks of the given request matches this parameterized path. If it
// matches, it will parse the appropriate path parameters into the request form and returns
// true. Otherwise, does nothing and returns false.
func (p *ParameterizedPath) ParseRequest(req *http.Request) bool {
	pathSplits := SplitPath(req.URL.Path)

	if !p.matches(pathSplits) {
		return false
	}

	for i, s := range p.Subs {
		req.Form.Add(s, pathSplits[i])
	}
	req.URL.Path = p.Path
	return true
}

// IsParameterizedPath returns true if the given string contains parameters, false otherwise
func IsParameterizedPath(path string) bool {
	return pathParamRegex.MatchString(path)
}

// ParseParameterizedPath returns a new ParameterizedPath pointer based on the provided URL path
func ParseParameterizedPath(path string) *ParameterizedPath {
	segments := SplitPath(path)

	paramPath := &ParameterizedPath{
		Path: path,
		Subs: make(map[int]string),
		Route: make(map[int]string),
		Length: len(segments),
	}

	for i, segment := range segments {
		match := pathParamRegex.FindStringSubmatch(segment)
		if len(match) > 1 {
			paramPath.Subs[i] = match[1]
		} else {
			paramPath.Route[i] = segment
		}
	}
	return paramPath
}

// SplitPath will split the given path string on the forward slash character
func SplitPath(path string) []string {
	return strings.Split(path, "/")
}
