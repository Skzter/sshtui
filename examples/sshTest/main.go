package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	args := os.Args

	// username, IP, Port
	programmArgs := args[1:]
	fmt.Println(programmArgs)
	/*
		name_and_ip := fmt.Sprintf("%s@%s", programmArgs[0], programmArgs[1])
		port := programmArgs[2]
	*/
	file, err := os.Create("/tmp/sshdata")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data := fmt.Sprintf("%s %s %s\n", programmArgs[0], programmArgs[1], programmArgs[2])
	n, err := file.WriteString(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("wrote %d bytes to the file\n", n)
	file.Sync()

	d, err := os.ReadFile("/tmp/sshdata")
	if err != nil {
		panic(err)
	}
	fmt.Println("data from file")
	fmt.Println(string(d))

	os.Remove("/tmp/sshdata")

	fmt.Println("after cleanup")
	cmd := exec.Command("ls", "/tmp/")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

}
