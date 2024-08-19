package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

var usr, _ = user.Current()
var cacheDir = filepath.Join(usr.HomeDir, ".cache/incrementor")
var cacheFile = filepath.Join(cacheDir, "id")

func main() {
	var currentValue int
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, 0700); err != nil {
			fmt.Fprintln(os.Stderr, "Error creating cache directory: ", err)
			os.Exit(1)
		}
		currentValue = 0
	} else {
		file, err := os.Open(cacheFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening file: ", err)
			os.Exit(1)
		} else {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			if scanner.Scan() {
				line := scanner.Text()
				currentValue, err = strconv.Atoi(line)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error converting value to integer: ", err)
					os.Exit(2)
				}
			} else {
				fmt.Fprintln(os.Stderr, "Unexpected content of cache file")
				os.Exit(2)
			}
		}
	}

	currentValue++
	fmt.Println(currentValue)

	file, err := os.Create(cacheFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating cache file: ", err)
		os.Exit(3)
	}
	defer file.Close()

	err = file.Chmod(0600)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error setting file mode: ", err)
		os.Exit(3)
	}

	valueStr := strconv.Itoa(currentValue)
	_, err = file.WriteString(valueStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error writing to cache file: ", err)
		os.Exit(3)
	}
}
