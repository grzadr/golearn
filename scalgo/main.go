package main

import (
	"flag"
	"fmt"

	"github.com/grzadr/golearn/scalgo/parsing"
	"github.com/grzadr/golearn/scalgo/path"
)

func main() {
	enlistmentPath := flag.String(
		"enlistment_path",
		"",
		"Path to the enlistment directory")
	maxUnits := flag.Int(
		"max_units",
		3,
		"Maximum number of units fot each output",
	)
	flag.Parse()

	if *enlistmentPath == "" {
		fmt.Println("Error: enlistment_path is required")
		flag.Usage()
		return
	}

	fmt.Printf("Enlistment path: %s\n", *enlistmentPath)

	enlistment_fs, enlistmant_path, err := path.ValidatePath(enlistmentPath)

	if err != nil {
		panic(err)
	}

	enlistment, err := parsing.NewEnlistmentFromFile(
		enlistment_fs,
		enlistmant_path,
		parsing.EmbeddedUnits,
	)

	if err != nil {
		panic(err)
	}

	for line := range enlistment.ScaleRecords(
		parsing.EmbeddedUnits,
		*maxUnits,
	) {
		_, err := fmt.Println(line)
		if err != nil {
			panic(err)
		}
	}
}
