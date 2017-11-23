package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	mgo "gopkg.in/mgo.v2"
)

type urlInfo struct {
	addrs   []string
	user    string
	pass    string
	db      string
	options map[string]string
}

func isOptSep(c rune) bool {
	return c == ';' || c == '&'
}

func extractURL(s string) (*urlInfo, error) {
	if strings.HasPrefix(s, "mongodb://") {
		s = s[10:]
	}
	info := &urlInfo{options: make(map[string]string)}
	if c := strings.Index(s, "?"); c != -1 {
		for _, pair := range strings.FieldsFunc(s[c+1:], isOptSep) {
			l := strings.SplitN(pair, "=", 2)
			if len(l) != 2 || l[0] == "" || l[1] == "" {
				return nil, errors.New("connection option must be key=value: " + pair)
			}
			info.options[l[0]] = l[1]
		}
		s = s[:c]
	}
	if c := strings.Index(s, "@"); c != -1 {
		pair := strings.SplitN(s[:c], ":", 2)
		if len(pair) > 2 || pair[0] == "" {
			return nil, errors.New("credentials must be provided as user:pass@host")
		}
		var err error
		info.user, err = url.QueryUnescape(pair[0])
		if err != nil {
			return nil, fmt.Errorf("cannot unescape username in URL: %q", pair[0])
		}
		if len(pair) > 1 {
			info.pass, err = url.QueryUnescape(pair[1])
			if err != nil {
				return nil, fmt.Errorf("cannot unescape password in URL")
			}
		}
		s = s[c+1:]
	}
	if c := strings.Index(s, "/"); c != -1 {
		info.db = s[c+1:]
		s = s[:c]
	}
	info.addrs = strings.Split(s, ",")
	return info, nil
}

func main() {
	url := "mongodb://10.34.42.52:7088,10.34.42.51:7088,10.34.42.53:7088/hms?replicaSet=kirk_rs1"
	uinfo, err := extractURL(url)
	fmt.Printf("uinfo = %+v\n", uinfo)
	fmt.Printf("err = %+v\n", err)

	uinfo2, err := mgo.ParseURL(url)
	fmt.Printf("uinfo2 = %+v\n", uinfo2)
	fmt.Printf("err = %+v\n", err)
}
