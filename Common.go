package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func readConfigFile(fileName string) []string {
	var configList []string
	// open the file
	file, err := os.Open(fileName)
	//handle errors while opening
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	fileScanner := bufio.NewScanner(file)
	// read line by line
	for fileScanner.Scan() {
		fmt.Println(fileScanner.Text())
		configList = append(configList, strings.TrimSpace(fileScanner.Text()))
	}
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	file.Close()
	return configList
}

func main() {
	fmt.Println(readConfigFile("/etc/passwd"))
}
