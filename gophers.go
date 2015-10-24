// Gophers runs scripts or commands concurrently and kills them when gophers was killed.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func main() {
	if len(os.Args) == 1 || os.Args[0] == "-h" || os.Args[0] == "-help" {
		fmt.Print(`Usage of gophers:
	gophers command1 command2 ...
Gophers runs scripts or commands concurrently and kills them when gophers was killed.
`)
		os.Exit(0)
	}
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
				fmt.Fprintf(os.Stderr, "error (sh -c %s): %s", cmd, err)
			}
			wg.Done()
		}(cmd)
	}
	wg.Wait()
}
