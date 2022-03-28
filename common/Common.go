package common

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func ReadConfigFile(fileName string) map[string]string {
	var configMap map[string]string
	configMap = make(map[string]string)
	// open the file
	file, err := os.Open(fileName)
	//handle errors while opening
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	fileScanner := bufio.NewScanner(file)
	// read line by line
	for fileScanner.Scan() {

		splitList := strings.Split(fileScanner.Text(), ":")
		configMap[splitList[0]] = splitList[1]
	}
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	file.Close()
	return configMap
}

func main() {
	fmt.Println(ReadConfigFile("/etc/passwd"))
}
