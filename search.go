package main

import (
	"errors"
	"regexp"
	"strings"
)

var (
	HopNotFound = errors.New("Can't find hop for given pattern")
)

func (hs *hops) searchHop(pattern string) (*hops, error) {
	fHops := make(hops)
	for name, h := range *hs {
		if match(name, pattern) {
			fHops[name] = h
		}
	}
	if len(fHops) == 0 {
		return nil, HopNotFound
	}
	return &fHops, nil
}

func match(name, pattern string) bool {
	if strings.ContainsRune(pattern, '*') {
		rPattern := strings.Replace(pattern, "*", ".*", -1)
		r := regexp.MustCompile("^" + rPattern + "$")
		return r.MatchString(name)
	} else {
		return name == pattern
	}
}
