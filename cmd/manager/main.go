package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spotinst/ocean-operator/internal/operator"
)

func main() {
	if err := operator.New().Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Exited with error: %v\n", err)
		os.Exit(1)
	}
}
