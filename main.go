package main

import (
	"fmt"

	"github.com/khy20040121/orbit/cmd/orbit"
)

func main() {
	err := orbit.Execute()
	if err != nil {
		fmt.Println("execute error: ", err.Error())
	}
}
