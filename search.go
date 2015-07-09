// Searching hops
package main

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// Returned when there's no hop for given name
	HopNotFound = errors.New("Can't find hop for given pattern")
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
		return nil, HopNotFound
	}
	return &fHops, nil
}

// Tests if name matches given pattern
func match(name, pattern string) bool {
	if strings.ContainsRune(pattern, '*') {
		rPattern := strings.Replace(pattern, "*", ".*", -1)
		r := regexp.MustCompile("^" + rPattern + "$")
		return r.MatchString(name)
	} else {
		return name == pattern
	}
}
