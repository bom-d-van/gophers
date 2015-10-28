// Gophers runs scripts or commands concurrently and kills them when gophers was killed.
//
// Usage
//  	gophers gulp "find . -iname '*.go' | entr -rd grm"
//
// What is grm: https://github.com/bom-d-van/bin/blob/master/grm
package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mitchellh/go-ps"
)

func main() {
	if len(os.Args) == 1 || os.Args[0] == "-h" || os.Args[0] == "-help" {
		fmt.Print(`Usage of gophers:
	gophers command1 command2 ...
Gophers runs scripts or commands concurrently and kills them when gophers was killed.
`)
		os.Exit(0)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Kill, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		pses, err := ps.Processes()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		kill(pses, os.Getpid())
	}()

	cmds := os.Args[1:]
	var wg sync.WaitGroup
	for _, cmd := range cmds {
		wg.Add(1)
		go func(cmd string) {
			ecmd := exec.Command("sh", "-c", cmd)
			ecmd.Stdout = os.Stdout
			ecmd.Stdin = os.Stdin
			ecmd.Stderr = os.Stderr
			if err := ecmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "error (sh -c %s): %s\n", cmd, err)
			}
			wg.Done()
		}(cmd)
	}
	wg.Wait()
}

func kill(pses []ps.Process, pid int) {
	extinct(pses, pid)

	fmt.Printf("kill %d\n", pid)
	ps, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("failed to find process (%d): %s\n", pid, err)
		return
	}
	if err := ps.Kill(); err != nil {
		fmt.Printf("failed to kill process (%d): %s\n", pid, err)
	}
}

func extinct(pses []ps.Process, pid int) {
	for _, ps := range pses {
		if ps.PPid() != pid {
			continue
		}
		kill(pses, ps.Pid())
	}
}
