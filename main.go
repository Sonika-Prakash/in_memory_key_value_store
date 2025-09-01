package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	s := newStore()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		command, key, attributes, keyExpiry, err := parseInput(line)
		if err != nil {
			log.Fatal(err)
		}
		s.execCmd(command, key, attributes, keyExpiry)
	}
}
