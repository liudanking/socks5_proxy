package main

import (
	"fmt"
)

func f(a ...interface{}) {
	for _, value := range a {
		if str, ok := value.(string); ok {
			fmt.Println(str)
		}
	}
}

func main() {
	f("123", "456")
}
