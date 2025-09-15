package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type sshEntry struct {
	IP       string
	Port     int
	Username string
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	sshDir := home + "/.ssh/known_hosts"
	data, err := os.Open(sshDir)
	if err != nil {
		panic(err)
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)

	var sshEntrys []sshEntry

	for scanner.Scan() {
		item := scanner.Text()
		startIP := strings.Index(item, "[")
		if startIP == -1 {
			continue
		}
		endIP := strings.Index(item, "]")
		IP := item[startIP+1 : endIP]
		startPort := strings.Index(item, ":")
		endPort := strings.Index(item, " ")
		PORT, err := strconv.Atoi(item[startPort+1 : endPort])
		if err != nil {
			panic(err)
		}
		entry := sshEntry{IP: IP, Port: PORT}
		if !slices.Contains(sshEntrys, entry) {
			sshEntrys = append(sshEntrys, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Println(sshEntrys)
}
