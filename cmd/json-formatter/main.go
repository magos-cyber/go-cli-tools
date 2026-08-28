package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "os"
)

func main() {
    indent := flag.String("indent", "  ", "Indentation string")
    compact := flag.Bool("compact", false, "Compact output")
    flag.Parse()

    data, err := io.ReadAll(os.Stdin)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
        os.Exit(1)
    }

    var parsed interface{}
    if err := json.Unmarshal(data, &parsed); err != nil {
        fmt.Fprintf(os.Stderr, "Invalid JSON: %v\n", err)
        os.Exit(1)
    }

    var output []byte
    if *compact {
        output, err = json.Marshal(parsed)
    } else {
        output, err = json.MarshalIndent(parsed, "", *indent)
    }
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error formatting: %v\n", err)
        os.Exit(1)
    }

    fmt.Println(string(output))
}
