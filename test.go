package main

import (
	"fmt"
	"strings"
	"github.com/bwesterb/go-ristretto"
)

func main() {
	device := "d1"

	deviceIds := [...]string{"d1", "d2"}

	result := StringInSlice(device, deviceIds[:])

	fmt.Printf("%s - %s: %v", device, deviceIds, result)
	// for _, b := range deviceIds {
	// 	if device == b {
	// 		fmt.Printf("true: %s", deviceIds)
	// 	}
	// }
}

func StringInSlice(a string, list []string) bool {

	for _, b := range list {
		if strings.Compare(a, b) == 0 {
			return true
		}
	}
	return false
}

func testCommitment(points []) {

}
