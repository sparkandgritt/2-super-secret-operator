package main

import (
	"crypto/sha1"
	"fmt"
)

func main() {
	text := "your_input_string"
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(text)))
	sixDigitHash := hash[:6]
	fmt.Println(sixDigitHash)
}
