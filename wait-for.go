package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"
)

const version = "0.0.1"

type addrs []string

func (a *addrs) String() string {
	return fmt.Sprint(*a)
}

func (a *addrs) Set(addr string) error {
	if ok, _ := regexp.MatchString(".+:\\d+", addr); !ok {
		return errors.New("Invalid format. host:port")
	}
	*a = append(*a, addr)
	return nil
}

var waits addrs
var wg sync.WaitGroup
var timeout int

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-w host:port...] <command>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Version: %s\n", version)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "  <command>\n\texecute command\n")
	}

	flag.Var(&waits, "w", "wait for host[s]. format: host:port")
	flag.IntVar(&timeout, "t", 300, "maximum time allowed for connection")
	log.SetFlags(0)
}

func dial(addr string) {
	for i := 0; i < timeout; i++ {
		if _, err := net.Dial("tcp", addr); err == nil {
			log.Printf("connected: %s\n", addr)
			wg.Done()
			break
		}
		if i%10 == 0 {
			log.Printf("waiting to be ready: %s\n", addr)
		}
		time.Sleep(1 * time.Second)
	}
}

func execCmd(command []string) {
	if len(command) == 0 {
		return
	}
	log.Printf("execute: %s\n", strings.Join(command, " "))
	binary, err := exec.LookPath(command[0])
	if err != nil {
		log.Fatalln(err)
	}
	syscall.Exec(binary, command, os.Environ())
}

func main() {
	flag.Parse()
	for _, addr := range waits {
		wg.Add(1)
		go dial(addr)
	}
	wg.Wait()

	execCmd(flag.Args())
}
