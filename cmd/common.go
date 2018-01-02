package cmd

import (
	"fmt"
	"os"
)

func must(err error) {
	if err == nil {
		return
	}

	fmt.Printf("ERROR: %+v\n", err)
	os.Exit(1)
}
