package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
)

func main() {
	FindMatch()
	//	searchString := "[]#:ab c"
	//	re := regexp.MustCompile(`[\-\[\]{}()*+?.,\\^$|#\s]`)
	//	searchString = re.ReplaceAllString(searchString, "\\$0")
	//	fmt.Println(searchString)

	//	re := regexp.MustCompile(`^[a-z0-9]+(?:[._-]?[a-z0-9]+)+$`)
	//	match := re.MatchString("a.b-c")
	//	fmt.Printf("match = %+v\n", match)
}

func IsAKSK(username, password string) bool {
	re := regexp.MustCompile(`[a-zA-z0-9_-]{40}`)
	return re.MatchString(username) && re.MatchString(password)
}

func MakeKey() string {
	var b [30]byte
	io.ReadFull(rand.Reader, b[:])
	return base64.URLEncoding.EncodeToString(b[:])
}

func FindMatch() {
	re := regexp.MustCompile(`(?i)x-version:\s+(\S+)$`)
	matches := re.FindSubmatch([]byte(`*/* X-Version: 3.2`))
	fmt.Printf("%q\n", matches[1])
}
