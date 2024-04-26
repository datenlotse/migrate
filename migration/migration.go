package migration

import (
	"fmt"
	"os"
	"strings"
)

var FailedToSplitError = fmt.Errorf("Failed to split Migration.\nDoes it contain a --THDS:Up and --THDS:Down comment in that order?")

type Migration struct {
	Up   *string
	Down *string
}

// Creates a migration struct from the provided path
func FromFile(path string) (*Migration, error) {
	contents, err := os.ReadFile(path)
	contentsString := string(contents)
	if err != nil {
		return nil, err
	}

	// Split the migration into Up and Down
	splits := strings.Split(contentsString, "\n--THDS:Down\n")
	if len(splits) != 2 {
		return nil, FailedToSplitError
	}

	upR := strings.Replace(splits[0], "--THDS:Up\n", "", -1)
	m := &Migration{
		Up:   &upR,
		Down: &splits[1],
	}

	return m, nil
}
