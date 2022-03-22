package main

import (
	"fmt"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func readFile(fileName string) string {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("read fail", err)
	}
	return string(f)
}

func GetAllFile(pathname string) ([]string, error) {
	result := []string{}

	fis, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Printf("读取文件目录失败，pathname=%v, err=%v \n", pathname, err)
		return result, err
	}

	// 所有文件/文件夹
	for _, fi := range fis {
		fullname := pathname + "/" + fi.Name()
		// 是文件夹则递归进入获取;是文件，则压入数组
		if fi.IsDir() {
			temp, err := GetAllFile(fullname)
			if err != nil {
				fmt.Printf("读取文件目录失败,fullname=%v, err=%v", fullname, err)
				return result, err
			}
			result = append(result, temp...)
		} else {
			result = append(result, fullname)
		}
	}

	return result, nil
}

func checkWebProcess(processNameStr string) bool {
	var webProcessList []string
	webProcessList = []string{"python", "node", "anywhere"}
	for _, webProcessStr := range webProcessList {
		if strings.Contains(processNameStr, webProcessStr) {
			return true
		}
	}
	return false
}

func ProcessName() (pname []string) {
	pids, _ := process.Pids()
	for _, pid := range pids {
		pn, _ := process.NewProcess(pid)

		var processName string
		processName, _ = pn.Name()
		if strings.Contains(processName, "java") {
			fmt.Print(pn.Pid)
			fmt.Print("\n")
			fmt.Print(processName)
			fmt.Print("\n")
			//fmt.Print(pn.Cmdline())
			fmt.Print("\n")

			var envList []string
			envList, _ = pn.Environ()
			var pathFlag bool = false
			var tomcatPathStr string = ""
			for _, s := range envList {
				if strings.Contains(s, "CATALINA_BASE") {
					tomcatPathStr = strings.Split(s, "=")[1]
					var jspPathStr string = filepath.Join(
						tomcatPathStr,
						"webapps",
					)
					_, err := os.Stat(jspPathStr)
					if err == nil {
						fmt.Print(jspPathStr)
						pathFlag = true
					}
				}
			}

			if pathFlag == false && tomcatPathStr != "" {
				var tomcatConfPathStr string = filepath.Join(
					tomcatPathStr,
					"conf/Catalina",
				)
				fileList, _ := GetAllFile(tomcatConfPathStr)
				for _, s := range fileList {
					var fileContentStr string = readFile(s)
					reg := regexp.MustCompile(`docBase="(.*?)"`)
					result := reg.FindAllStringSubmatch(fileContentStr, -1)
					for _, mathList := range result {
						pathFlag = true
						fmt.Print(mathList[1])
					}
				}
			}

			fmt.Print("\n")

		} else if checkWebProcess(processName) {
			fmt.Print(pn.Pid)
			fmt.Print("\n")
			fmt.Print(processName)
			fmt.Print("\n")
			var connectionList []net.ConnectionStat
			connectionList, _ = pn.Connections()

			var portFlag = false
			for _, stat := range connectionList {
				if stat.Status == "LISTEN" {
					portFlag = true
					fmt.Print(stat.Laddr.IP, ":", stat.Laddr.Port)
					fmt.Print("\n")
				}
			}
			if portFlag == true {
				var cwdPath string
				cwdPath, _ = pn.Cwd()
				fmt.Print(cwdPath)

				fmt.Print("\n")

				var cmdLine string
				cmdLine, _ = pn.Cmdline()
				fmt.Print(cmdLine)
			}

			fmt.Print("\n")

		} else {
			//var connectionList []net.ConnectionStat
			//connectionList, _ = pn.Connections()
			//var portFlag = false
			//for _, stat := range connectionList {
			//	if stat.Status == "LISTEN" {
			//		portFlag = true
			//		fmt.Print(stat.Laddr.IP, ":", stat.Laddr.Port)
			//		fmt.Print("\n")
			//	}
			//}
			//if portFlag && strings.Index(processName, "[") != 0 {
			//	fmt.Print(pn.Pid)
			//	fmt.Print("\n")
			//	fmt.Print(processName)
			//	fmt.Print("\n")
			//
			//	var cwdPath string
			//	cwdPath, _ = pn.Cwd()
			//	fmt.Print(cwdPath)
			//
			//	fmt.Print("\n")
			//
			//	var cmdLine string
			//	cmdLine, _ = pn.Cmdline()
			//	fmt.Print(cmdLine)
			//
			//	fmt.Print("\n")
			//}
		}

		pName, _ := pn.Name()
		pname = append(pname, pName)
	}
	return pname
}

func main() {
	//fmt.Print(ProcessName())
	ProcessName()
}
