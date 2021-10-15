package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	out bytes.Buffer
)

func readTestCase(src string) (listFile []string, err error) {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if fileName := file.Name(); strings.Contains(fileName, ".in") {
			listFile = append(listFile, fileName)
		}
	}
	return listFile, nil
}

func unitTest(src string) (data string, err error) {
	f, err := os.ReadFile(src)
	if err != nil {
		return "", err
	}
	return string(f), nil

}
func buildFile(src *string) (string, error) {
	_, file := filepath.Split(*src)
	ext := filepath.Ext(file)
	if ext == ".cpp" {
		newFile := fmt.Sprintf("%s.exe", strings.Split(file, ".cpp")[0])
		cmd := exec.Command("g++", file, "-o", newFile)
		_, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return newFile, nil
	} else if ext == ".c" {
		newFile := fmt.Sprintf("%s.exe", strings.Split(file, ".c")[0])
		cmd := exec.Command("gcc", file, "-o", newFile)
		_, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return newFile, nil
	} else {
		return "", errors.New("[+] Invalid file extension!")
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	args := os.Args
	if len(args) != 5 {
		_, file := filepath.Split(args[0])
		fmt.Printf("Usage: %s -f file.cpp -d test\n-f: file needs to be run\n-d: path to folder testcase", file)
	}
	fmt.Println("[+] Building...")
	startBuilding := time.Now()
	newFile, err := buildFile(&args[2])
	if err != nil {
		fmt.Println("[+] Compile Error!")
		panic(err)
	}
	timeToBuild := time.Since(startBuilding)
	fmt.Println("[+] Compiled ", timeToBuild)
	listTest, err := readTestCase(args[4])
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range listTest {
		f, err := unitTest(fmt.Sprintf("%s\\%s\\%s", dir, args[4], file))
		timeStart := time.Now()
		cmd := exec.Command(newFile)
		cmd.Stdin = strings.NewReader(f + "\n")
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
		numberTest := strings.Split(file, ".in")[0]
		f, err = unitTest(fmt.Sprintf("%s\\%s\\%s.%s", dir, args[4], numberTest, "out"))
		if err != nil {
			panic(err)
		}
		if out.String() == f {
			fmt.Printf("[+] ✓ Case %s: passed %s\n", numberTest, time.Since(timeStart))
		} else {
			fmt.Printf("[+] ✕ Case %s: failed %s\n", numberTest, time.Since(timeStart))
		}
		out.Reset()
	}
	fmt.Println("[+] Done!")
}
