package internal

import (
	"fmt"
	"os"
)

func Output(data interface{}) error {
	// format ?
	fmt.Fprintf(os.Stdout, "%v\n", data)

	return nil
}
