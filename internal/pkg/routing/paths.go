package routing

import (
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
