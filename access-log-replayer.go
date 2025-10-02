package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	inputFile := flag.String("input-file", "", "Path to ELF log file")
	httpHost := flag.String("http_host", "", "Target HTTP host, e.g. localhost:8983")
	flag.Parse()

	if *inputFile == "" || *httpHost == "" {
		fmt.Println("Usage: --input-file <elf_log_file.log> --http_host <host:port>")
		os.Exit(1)
	}

	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue // skip comments and empty lines
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue // skip malformed lines
		}

		method := fields[3]
		url := fields[4]
		if method != "GET" {
			continue // only replay GET requests
		}

		fullURL := "http://" + *httpHost + url
		resp, err := http.Get(fullURL)
		if err != nil {
			fmt.Printf("Request to %s failed: %v\n", fullURL, err)
			continue
		}
		fmt.Printf("%s -> %d\n", fullURL, resp.StatusCode)
		resp.Body.Close()
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}
