package main

import (
	"fmt"
	"log"

	ldap "gopkg.in/ldap.v3"
)

func main() {
	username := "lutao"
	password := "lutao"

	bindusername := "cn=lutao,ou=People,dc=xxx,dc=com"
	bindpassword := "lutao"

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "ldap.dev.xxx.io", 389))
	//l, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", "ldaps.xxx.io", 636), &tls.Config{})
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer l.Close()

	// First bind with a read only user
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		log.Fatal("bind:", err)
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		"dc=xxx,dc=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%v=%v)", "cn", username),
		[]string{"cn", "sn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal("search:", err)
	}

	if len(sr.Entries) != 1 {
		log.Fatal("User does not exist or too many entries returned")
	}

	userdn := sr.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userdn, password)
	if err != nil {
		log.Fatal("bind2:", err)
	}
	fmt.Printf("userdn = %+v\n", userdn)

	// Rebind as the read only user for any further queries
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		log.Fatal(err)
	}
}
