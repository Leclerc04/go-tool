package strs

// https://github.com/serenize/snaker/blob/master/snaker.go

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	acronyms = func() *regexp.Regexp {
		var ss []string
		for k := range commonInitialisms {
			ss = append(ss, k[:1]+strings.ToLower(k[1:]))
		}
		return regexp.MustCompile(fmt.Sprintf("(%s)$", strings.Join(ss, "|")))
	}()
	camelcase = regexp.MustCompile(`(?m)[-.$/:_{}\s]`)
)

// CamelToSnake converts a given string to snake case
func CamelToSnake(s string) string {
	var result string
	var words []string
	var lastPos int
	rs := []rune(s)

	for i := 0; i < len(rs); i++ {
		if i > 0 && (unicode.IsUpper(rs[i]) || unicode.IsDigit(rs[i])) {
			if initialism := startsWithInitialism(s[lastPos:]); initialism != "" {
				words = append(words, initialism)

				i += len(initialism) - 1
				lastPos = i
				continue
			}
			if unicode.IsDigit(rs[i]) {
				continue
			}

			words = append(words, s[lastPos:i])
			lastPos = i
		}
	}

	// append the last word
	if s[lastPos:] != "" {
		words = append(words, s[lastPos:])
	}

	for k, word := range words {
		if k > 0 {
			result += "_"
		}

		result += strings.ToLower(word)
	}

	return result
}

// SnakeToCamel returns a string converted from snake case to uppercase
func SnakeToCamel(s string, capFirst bool) string {
	return Depunct(s, capFirst)
}

// startsWithInitialism returns the initialism if the given string begins with it
func startsWithInitialism(s string) string {
	var initialism string
	// the longest initialism is 5 char, the shortest 2
	for i := 1; i <= 5; i++ {
		if len(s) > i-1 && commonInitialisms[s[:i]] {
			initialism = s[:i]
		}
	}
	return initialism
}

// Depunct normalize names to camel case.
func Depunct(ident string, initialCap bool) string {
	if ident == "" {
		return ""
	}
	if ident == "_id" {
		if initialCap {
			return "ID"
		}
		return "id"
	}
	if ident == "_rev" {
		if initialCap {
			return "Rev"
		}
		return "rev"
	}

	matches := camelcase.Split(ident, -1)
	for i, m := range matches {
		if initialCap || i > 0 {
			m = capFirst(m)
		}
		matches[i] = acronyms.ReplaceAllStringFunc(m, func(c string) string {
			if c == "CName" {
				return "CName"
			}
			if c == "Ids" {
				return "IDs"
			}
			if c == "Urls" {
				return "URLs"
			}
			return strings.ToUpper(c)
		})
		if i > 0 && matches[i] == "US" && matches[i-1] == "EN" {
			matches[i-1] = "En"
		}
	}
	ret := strings.Join(matches, "")
	return ret
}

func capFirst(ident string) string {
	if ident == "" {
		panic(fmt.Sprintf("invalid ident: %#v", ident))
	}
	r, n := utf8.DecodeRuneInString(ident)
	return string(unicode.ToUpper(r)) + ident[n:]
}

// commonInitialisms, taken from
// https://github.com/golang/lint/blob/32a87160691b3c96046c0c678fe57c5bef761456/lint.go#L702
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"C2C":   true,
	"CN":    true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EN":    true,
	"EOF":   true,
	"GMAT":  true,
	"GPA":   true,
	"GPS":   true,
	"GRE":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IDs":   true,
	"IELTS": true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"QQ":    true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TOEFL": true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"URI":   true,
	"URL":   true,
	"URLs":  true,
	"US":    true,
	"UTF8":  true,
	"UUID":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
	"GQL":   true,
}
