package commands

import (
	"context"
	"fmt"
)

const Version string = "v0.1.0-alpha"

type VersionRunner struct{}

func (r VersionRunner) Run(_ context.Context) error {
	fmt.Printf("CSB Open API %v\n", Version)
	return nil
}
