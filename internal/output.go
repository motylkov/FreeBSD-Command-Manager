package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

func Output(data interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode terraform output: %v\n", err)
		return err
	}
	return nil
}
