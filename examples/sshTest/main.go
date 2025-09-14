package main

import (
	"bufio"
	"fmt"
	"os"
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

	sshEntrys := make(map[string]sshEntry)

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
		sshEntrys[IP] = sshEntry{IP: IP, Port: PORT}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	for {
		fmt.Println(sshEntrys)
	}
}
