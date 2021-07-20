package globals

import "fmt"

type OldVersionError struct {
	Version string
}

func (e OldVersionError) Error() string {
	return fmt.Sprintf("Structure is of not currect version: %s", e.Version)
}
