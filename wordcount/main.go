//go:build !solution

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	files := os.Args[1:]
	counter := make(map[string]int)

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()

		sc := bufio.NewScanner(f)
		for sc.Scan() {
			counter[sc.Text()]++
		}
	}

	for k, v := range counter {
		if v >= 2 {
			fmt.Printf("%d\t%s\n", v, k)
		}
	}
}
