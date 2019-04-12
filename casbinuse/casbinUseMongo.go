package main

import (
	"fmt"

	"github.com/casbin/casbin"
	"github.com/casbin/mongodb-adapter"
)

func main() {
	a := mongodbadapter.NewAdapter("127.0.0.1:27017/casbin") // Your MongoDB URL.

	e := casbin.NewEnforcer("./rbac_with_domains_model.conf", a)

	// Load the policy from DB.
	e.LoadPolicy()

	// Check the permission.
	//e.Enforce("alice", "data1", "read")

	// Modify the policy.
	e.AddPolicy("member", "domain1", "data1", "read")
	e.AddPolicy("admin", "domain1", "data1", "write")
	e.AddPolicy("member", "domain2", "data2", "read")
	e.AddPolicy("admin", "domain2", "data2", "write")
	// e.RemovePolicy(...)
	e.AddGroupingPolicy("alice", "admin", "domain1")
	e.AddGroupingPolicy("alice", "member", "domain1")

	roles := e.GetRolesForUserInDomain("alice", "domain1")
	fmt.Printf("roles = %+v\n", roles)

	// Save the policy back to DB.
	e.SavePolicy()
}
