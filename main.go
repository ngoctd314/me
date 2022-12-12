package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open(".env")
	if err != nil {
		log.Println(err)
	}
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		tmp := strings.Split(fileScanner.Text(), "=")
		if len(tmp) != 2 {
			continue
		}
		os.Setenv(tmp[0], tmp[1])
	}
	fmt.Println(os.Getenv("PORT"))

}
