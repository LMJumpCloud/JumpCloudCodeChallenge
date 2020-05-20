package routing

import (
	"regexp"
	"strings"
)

var pathParamRegex = regexp.MustCompile(`\{(.+?)\}`)

type ParameterizedPath struct {
	Path string
	Length int
	Route map[int]string // map of path index -> path segment
	Subs map[int]string // map of path index -> parameter name
}

func IsParameterizedPath(path string) bool {
	return pathParamRegex.MatchString(path)
}

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

func SplitPath(path string) []string {
	return strings.Split(path, "/")
}
