package util

import (
	"rhinoman/couchdb-go"
	"strings"
	"time"
	t_conf "tvm-sdk/config"
	"unicode"
	"unicode/utf8"
)

const (
	special_tag = rune('$')
)

func CheckCouchDBExists(dbName string) bool {
	var timeout = time.Duration(500 * time.Millisecond)
	coon, err := couchdb.NewConnection("127.0.0.1", t_conf.TriasConfig.CouchdbInfo.Port, timeout)
	if (err != nil) {
		return false
	}
	auth := couchdb.BasicAuth{Username: t_conf.TriasConfig.CouchdbInfo.Username, Password: t_conf.TriasConfig.CouchdbInfo.Password}
	db := coon.SelectDB(dbName, &auth)
	if _, err := db.Compact(); err != nil {
		return false

	}
	return true
}

func GetCouchDBName(s string) string {
	isASCII, hasUpper := true, false
	var count int = 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= utf8.RuneSelf {
			isASCII = false
			break
		}
		if c >= 'A' && c <= 'Z' {
			count ++
		}
		hasUpper = hasUpper || (c >= 'A' && c <= 'Z')
	}

	if isASCII { // optimize for ASCII-only strings.
		if !hasUpper {
			return s
		}
		b := make([]byte, len(s)+count)
		count = 0
		for i := 0; i < len(s); i++ {
			c := s[i]
			if c >= 'A' && c <= 'Z' {
				c += 'a' - 'A'
				b[count] = byte(special_tag)
				count++
				b[count] = c
			} else {
				b[count] = c
			}
			count++
		}
		return string(b)
	}
	return strings.Map(unicode.ToLower, s)
}
