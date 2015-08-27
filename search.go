// Searching hops
package main

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrHopNotFound is returned when there's no hop for given name
	ErrHopNotFound = errors.New("Can't find hop for given pattern")
)

// Search hops with matchins names
func (hs *hops) searchHop(pattern string) (*hops, error) {
	fHops := make(hops)
	for name, h := range *hs {
		if match(name, pattern) {
			fHops[name] = h
		}
	}
	if len(fHops) == 0 {
		return nil, ErrHopNotFound
	}
	return &fHops, nil
}

// Tests if name matches given pattern
func match(name, pattern string) bool {
	if strings.ContainsRune(pattern, '*') {
		rPattern := strings.Replace(pattern, "*", ".*", -1)
		r := regexp.MustCompile("^" + rPattern + "$")
		return r.MatchString(name)
	}
	return name == pattern
}
