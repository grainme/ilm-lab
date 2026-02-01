package cache

import (
	"fmt"

	"github.com/google/uuid"
)

const moviePrefix = "movie:"

func MovieKey(id uuid.UUID) string {
	return fmt.Sprintf("%s%s", moviePrefix, id.String())
}
