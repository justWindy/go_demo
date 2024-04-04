package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestBufioScanner(t *testing.T) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter some text:")
	for scanner.Scan() {
		text := scanner.Text()

		fmt.Println("You entered:", text)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
