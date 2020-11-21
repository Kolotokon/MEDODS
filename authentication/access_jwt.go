package main

import (
	"crypto/sha512"
	"fmt"
	"io"
)

func main() {
	h512 := sha512.New()
	io.WriteString(h512, "Name")
	fmt.Printf("Name : %x\n", h512.Sum(nil))

	data := []byte("Autentivication")
	fmt.Printf("Password : %x\n", sha512.Sum512(data))
}
