package internal

import (
	"fmt"
	"os"
)

// OutputTerraform prints the given data as pretty JSON to stdout for Terraform integration.
func Output(data interface{}) error {
	// format ?
	fmt.Fprintf(os.Stdout, "%v\n", data)

	return nil
}
