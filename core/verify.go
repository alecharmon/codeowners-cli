package core

import (
	"os"

	codeowners "github.com/alecharmon/codeowners/pkg"
	"github.com/fatih/color"
)

//MustGetCodeOwners ...
func MustGetCodeOwners(filePath string) *codeowners.CodeOwners {
	co, errs := codeowners.BuildFromFile(filePath)
	if errs != nil && len(errs) > 0 {
		color.New(color.FgRed).Printf("Could not load codeowner file %s:\n", filePath)
		for _, e := range errs {
			color.New(color.FgRed).Println(e.Error())
		}
		os.Exit(1)
		return nil
	}
	return co
}
