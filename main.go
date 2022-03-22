package main

import (
	"EmergencyTool/java"
	"EmergencyTool/web"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("main web")
		fmt.Println("main java /path")
	}
	if os.Args[1] == "web" {
		web.GetWebProcessInfo()
	} else if os.Args[1] == "java" {
		targetJarPath := os.Args[2]
		java.Start(targetJarPath)
	}

}
