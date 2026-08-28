package main

import (
    "flag"
    "fmt"
    "os"
    "os/exec"
    "time"
)

func main() {
    path := flag.String("path", ".", "Path to watch")
    interval := flag.Int("interval", 5, "Check interval in seconds")
    command := flag.String("cmd", "", "Command to run on change")
    flag.Parse()

    fmt.Printf("Watching: %s\n", *path)
    fmt.Printf("Interval: %ds\n", *interval)

    var lastMod time.Time
    for {
        info, err := os.Stat(*path)
        if err != nil {
            time.Sleep(time.Duration(*interval) * time.Second)
            continue
        }

        if info.ModTime() != lastMod && !lastMod.IsZero() {
            fmt.Printf("Change detected in %s\n", *path)
            if *command != "" {
                cmd := exec.Command("sh", "-c", *command)
                cmd.Stdout = os.Stdout
                cmd.Stderr = os.Stderr
                cmd.Run()
            }
        }
        lastMod = info.ModTime()
        time.Sleep(time.Duration(*interval) * time.Second)
    }
}
