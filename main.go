package main

import (
	"strings"
	"fmt"
	"bufio"
	"os"
)


func cleanInput(in string) []string {
	fields := strings.Fields(in)
	for i, s := range fields {
		fields[i] = strings.ToLower(s)
	}
	return fields
}

func main() {
	running := true
    scanner := bufio.NewScanner(os.Stdin)

	for running {
		fmt.Print("Pokedex > ")

		scanner.Scan()

		inputFields := cleanInput(scanner.Text())
		if inputFields[0] == "quit" {
			running = false
			break
		}
		fmt.Println("Your command was:", inputFields[0])
	}
}