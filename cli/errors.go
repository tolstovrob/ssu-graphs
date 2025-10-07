package cli

import "fmt"

func ThrowInitRoot() error {
	return fmt.Errorf("Failed to init root of CLI")
}
