package main

import (
	"log"

	"github.com/nerney/dappy"
)

func main() {

	//create a new client
	client := dappy.New(dappy.Options{
		BaseDN:       "dc=xxx,dc=com",
		Filter:       "cn",
		BaseUser:     "cn=lutao,ou=People,dc=xxx,dc=com",
		BasePassword: "lutao",
		URL:          "ldap.dev.xxx.com:389",
		Attrs:        []string{"cn", "sn"},
	})

	//username and password to authenticate
	username := "lutao"
	password := "lutao"

	//attempt the authentication
	err := client.Authenticate(username, password)

	//see the results
	if err != nil {
		log.Print(err)
	} else {
		log.Print("user successfully authenticated!")
	}

	//get a user entry
	user, err := client.GetUserEntry(username)
	if err == nil {
		user.PrettyPrint(2)
	}
}
