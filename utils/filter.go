package utils

import (
	log "github.com/sirupsen/logrus"
	"regexp"
)

func BuildPatternsFromStrings(names []string) []*regexp.Regexp {
	var lists []*regexp.Regexp
	for _, v := range names {
		re, err := regexp.Compile(v)
		if err != nil {
			log.Fatal("Invalid patthern found in bucker filter param")
		}
		lists = append(lists, re)
	}
	return lists
}

func MatchExclude(filters []*regexp.Regexp, name string) bool {
	for _, filter := range filters {
		if filter.MatchString(name) {
			return true
		}
	}
	return false
}
