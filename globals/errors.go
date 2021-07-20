package globals

import "fmt"

type NotCurrentVersionError struct {
	Version string
}

func (e NotCurrentVersionError) Error() string {
	return fmt.Sprintf("Structure is of not currect version: %s", e.Version)
}
