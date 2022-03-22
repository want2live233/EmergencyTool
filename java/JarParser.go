package java

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	filepath "path/filepath"
	"strings"
)

func getTargetJarFile(pathname string) ([]string, error) {
	var result []string
	fis, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Printf("读取文件目录失败，pathname=%v, err=%v \n", pathname, err)
		return result, err
	}
	// 所有文件/文件夹
	for _, fi := range fis {
		fullname := pathname + "/" + fi.Name()

		fileExt := path.Ext(fi.Name())

		// 是文件夹则递归进入获取;是文件，则压入数组
		if fi.IsDir() {
			temp, err := getTargetJarFile(fullname)
			if err != nil {
				fmt.Printf("读取文件目录失败,fullname=%v, err=%v", fullname, err)
				return result, err
			}
			result = append(result, temp...)
		} else {
			if fileExt == ".jar" {
				result = append(result, fullname)
			}

		}
	}
	return result, nil
}

func parseSingleJar(filename string) ([]string, []string) {
	var packageList []string
	var jarList []string
	// Open a zip archive for reading.
	r, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	// Iterate through the files in the archive,
	// printing some of their contents.

	for _, f := range r.File {
		filePath, fileName := filepath.Split(f.Name)
		fileExt := path.Ext(fileName)

		if fileExt == ".jar" {
			jarList = append(jarList, fileName)
		}
		if fileExt != ".class" || "" == filePath {
			continue
		}

		var inFlag = false
		for _, s := range packageList {
			if strings.Contains(filePath, s) {
				inFlag = true
				break
			}
		}

		if inFlag {
			continue
		}

		packageList = append(packageList, filePath)

	}

	return packageList, jarList
}

func Start(targetJarPath string) {

	var riskPackageList []string
	var riskJarLIst []string
	riskPackageList = []string{"com.thoughtworks.xstream", "com.alibaba.fastjson", "com.alibaba.dubbo", "org.apache.dubbo"}

	riskJarLIst = []string{"xstream-", "fastjson-", "dubbo-"}

	jarFilelist, _ := getTargetJarFile(targetJarPath)

	var allPackageList []string
	var allJarList []string
	for _, s := range jarFilelist {
		packageList, jarList := parseSingleJar(s)
		allPackageList = append(allPackageList, packageList...)
		allJarList = append(allJarList, jarList...)
	}

	for _, s := range allPackageList {
		s = strings.ReplaceAll(s, "/", ".")[:len(s)-1]
		for _, s2 := range riskPackageList {
			if strings.Contains(s, s2) {
				fmt.Println("有风险：", s)
			}
		}
	}

	fmt.Print("======================\n")

	for _, s := range allJarList {
		for _, s2 := range riskJarLIst {
			if strings.Contains(s, s2) {
				fmt.Println("有风险：", s)
			}
		}

	}
}
func main() {
	Start("")
}
