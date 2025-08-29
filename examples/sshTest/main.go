package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type sshEntry struct {
	IP   string
	Port int
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
		fmt.Println(scanner.Text())
		item := scanner.Text()
		startIP := strings.Index(item, "[")
		if startIP == -1 {
			fmt.Println("No IP found")
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
		sshEntrys = append(sshEntrys, sshEntry{IP: IP, Port: PORT})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Println(sshEntrys)
}
