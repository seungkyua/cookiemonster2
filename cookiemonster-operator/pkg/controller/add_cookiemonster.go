package controller

import (
	"github.com/rbxorkt12/coockiemonster2/pkg/controller/cookiemonster"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, cookiemonster.Add)
}
