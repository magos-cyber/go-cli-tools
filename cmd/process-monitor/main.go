package main

import (
    "flag"
    "fmt"
    "os/exec"
    "strings"
    "time"
)

func main() {
    interval := flag.Int("interval", 5, "Update interval in seconds")
    filter := flag.String("filter", "", "Filter processes")
    flag.Parse()

    fmt.Println("=== Process Monitor ===")
    fmt.Printf("Interval: %ds\n\n", *interval)

    for {
        cmd := exec.Command("ps", "aux")
        output, err := cmd.Output()
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            return
        }

        lines := strings.Split(string(output), "\n")
        for _, line := range lines {
            if *filter == "" || strings.Contains(line, *filter) {
                fmt.Println(line)
            }
        }

        time.Sleep(time.Duration(*interval) * time.Second)
    }
}
