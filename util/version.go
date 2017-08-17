package util

import "fmt"

// Version will be populated at build time and houses the current version.
var Version = "undefined"

// ShowVersion will print the version information.
func ShowVersion() {
	fmt.Println("pr0nbot " + Version)
}
