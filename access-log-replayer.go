package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	accesslog "github.com/nekrassov01/access-log-parser"
)

func main() {
	inputFile := flag.String("input-file", "", "Path to CLF log file")
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

	handler := func(labels, values []string, isFirst bool) (string, error) {
		var method, path string
		for i, label := range labels {
			if label == "method" {
				method = values[i]
			}
			if label == "request_uri" {
				path = values[i]
			}
		}
		fmt.Printf("DEBUG: method='%s', path='%s'\n", method, path)
		if method != "GET" || path == "" {
			fmt.Println("DEBUG: Not a GET request or missing path, skipping")
			return "", nil
		}
		fullURL := "http://" + *httpHost + path
		fmt.Printf("DEBUG: Sending GET to %s\n", fullURL)
		resp, err := http.Get(fullURL)
		if err != nil {
			fmt.Printf("Request to %s failed: %v\n", fullURL, err)
			return "", nil
		}
		fmt.Printf("%s -> %d\n", fullURL, resp.StatusCode)
		resp.Body.Close()
		return "", nil
	}

	opt := accesslog.Option{LineHandler: handler}
	parser := accesslog.NewApacheCLFRegexParser(context.Background(), os.Stdout, opt)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue // skip empty lines
		}
		// ParseString will invoke the handler for each line
		_, _ = parser.ParseString(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}
