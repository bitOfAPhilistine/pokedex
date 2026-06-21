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
		if len(inputFields) == 0 {
			continue
		} else {
			err := command(inputFields[0])
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}