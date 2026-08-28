package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"
)

func main() {
    logFile := flag.String("file", "", "Log file to analyze")
    pattern := flag.String("pattern", "", "Pattern to search")
    flag.Parse()

    if *logFile == "" {
        fmt.Println("Usage: log-analyzer --file <logfile> [--pattern <pattern>]")
        os.Exit(1)
    }

    file, err := os.Open(*logFile)
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    count := 0
    for scanner.Scan() {
        line := scanner.Text()
        if *pattern == "" || strings.Contains(line, *pattern) {
            count++
            if *pattern != "" {
                fmt.Println(line)
            }
        }
    }

    if *pattern == "" {
        fmt.Printf("Total lines: %d\n", count)
    } else {
        fmt.Printf("Matches: %d\n", count)
    }
}
