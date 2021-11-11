package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/elliotchance/sshtunnel"
	sshconfig "github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type config struct {
	host   string
	method ssh.AuthMethod
	target string
	listen string
}

func main() {
	c := getConfig()
	tunnel := sshtunnel.NewSSHTunnel(c.host, c.method, c.target, "")
	tunnel.Local = sshtunnel.NewEndpoint(c.listen)
	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	tunnel.Start()
}

func usage() {
	fmt.Println("Usage: pivot-ssh <remote_host> <local_listener> <forward_to>")
	fmt.Println()
	fmt.Println("  Example:")
	fmt.Println()
	fmt.Println("  $ pivot-ssh 10.11.1.123 127.0.0.1:8080 10.1.1.55:80")
	fmt.Println()
	os.Exit(1)
}

func getConfig() config {
	args := os.Args
	if len(args) != 4 {
		usage()
	}
	host := args[1]
	listen := args[2]
	target := args[3]
	if host == "" || listen == "" || target == "" {
		usage()
	}
	if user := sshconfig.Get(host, "User"); user == "" {
		log.Fatalf("SSH config did not contain user for host %s\n", host)
		os.Exit(1)
	} else {
		host = user + "@" + host
	}
	port := "22"
	if portOverride := sshconfig.Get(host, "Port"); portOverride != "" {
		port = portOverride
	}
	host = host + ":" + port
	var method ssh.AuthMethod
	if key := sshconfig.Get(host, "IdentityFile"); key != "" && !strings.HasSuffix(key, "identity") {
		method = sshtunnel.PrivateKeyFile(key)
	} else {
		fmt.Print("Enter SSH password: ")
		b, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println()
		method = ssh.Password(string(b))
	}
	return config{
		host:   host,
		method: method,
		target: target,
		listen: listen,
	}
}
