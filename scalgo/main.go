package main

import (
    "flag"
    "fmt"
	
)

func main() {
    enlistmentPath := flag.String("enlistment_path", "", "Path to the enlistment directory")
    flag.Parse()

    if *enlistmentPath == "" {
        fmt.Println("Error: enlistment_path is required")
        flag.Usage()
        return
    }

    fmt.Printf("Enlistment path: %s\n", *enlistmentPath)
}
