package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	if len(os.Args) == 2 {
		file1 := os.Args[1]
		singleSetupMode(file1)
	} else if len(os.Args) == 3 {
		file1 := os.Args[1]
		file2 := os.Args[2]
		dualSetupMode(file1, file2)
	} else {
		processAllJsonFiles()
	}
}

func singleSetupMode(file1 string) {
	content1, err := ioutil.ReadFile(file1)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file1, err)
		return
	}
	urlEncodedSetup1 := url.QueryEscape(string(content1))
	fmt.Printf("Single setup mode for %s\n", file1)
	url := fmt.Sprintf("https://www.accsetupcomparator.com?source=mk&mode=single&filename1=%s&setup1=%s", url.QueryEscape(file1), urlEncodedSetup1)
	openBrowser(url)
}

func dualSetupMode(file1, file2 string) {
	content1, err := ioutil.ReadFile(file1)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file1, err)
		return
	}
	urlEncodedSetup1 := url.QueryEscape(string(content1))

	content2, err := ioutil.ReadFile(file2)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file2, err)
		return
	}
	urlEncodedSetup2 := url.QueryEscape(string(content2))

	fmt.Printf("Dual setup mode for %s and %s\n", file1, file2)
	url := fmt.Sprintf("https://www.accsetupcomparator.com?source=mk&mode=compare&filename1=%s&filename2=%s&setup1=%s&setup2=%s",
		url.QueryEscape(file1), url.QueryEscape(file2), urlEncodedSetup1, urlEncodedSetup2)
	openBrowser(url)
}

func processAllJsonFiles() {
	files, err := filepath.Glob("*.json")
	if err != nil {
		fmt.Printf("Error reading JSON files: %v\n", err)
		return
	}

	if len(files) < 2 {
		fmt.Println("Not enough JSON files to compare")
		return
	}

	var firstFilename string
	for _, filename := range files {
		if firstFilename != "" {
			fmt.Printf("%s -> %s\n", firstFilename, filename)
			dualSetupMode(firstFilename, filename)
		}
		firstFilename = filename
	}
}

// openBrowser opens the specified URL in the default web browser
func openBrowser(url string) {
	var err error
	switch os := runtime.GOOS; os {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Printf("Error opening URL: %v\n", err)
	}
}
