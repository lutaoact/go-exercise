package main

import (
	"fmt"

	"github.com/casbin/casbin"
)

func main() {
	e := casbin.NewEnforcer("./rbac_with_domains_model.conf", "./rbac_with_domains_policy.csv", true)
	sub := "alice"
	obj := "data1"
	dom := "domain1"
	act := "write"

	if e.Enforce(sub, dom, obj, act) {
		fmt.Println("ok")
	} else {
		fmt.Println("not ok")
	}

	roles := e.GetRolesForUserInDomain("alice", "domain1")
	fmt.Printf("roles = %+v\n", roles)

	permissions := e.GetPermissionsForUserInDomain("admin", "domain2")
	fmt.Printf("permissions = %+v\n", permissions)

	val := e.DeleteRoleForUserInDomain("alice", "member", "domain2")
	//val := e.DeleteRoleForUserInDomain("alice", "admin", "domain2")
	fmt.Println(val)
}
