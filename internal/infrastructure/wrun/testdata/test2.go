package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println(`{"name":"test_plugin","version":"2.0"}`)

	s := bufio.NewScanner(os.Stdin)
	s.Split(bufio.ScanLines)
	if !s.Scan() {
		return
	}
	line := s.Text()

	fmt.Println(line)
}
