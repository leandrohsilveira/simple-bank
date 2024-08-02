package components

import (
	"fmt"
)

const MainId = "main"

func RefId(id string) string {
	return fmt.Sprintf("#%s", id)
}
