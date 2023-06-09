package helper

import (
	"fmt"
	"strconv"
)

func GenerateNextID(highestID int) string {
	// Increment the highest ID by 1
	nextID := highestID + 1

	// Convert the incremented ID back to a string with leading zeros
	id := strconv.Itoa(nextID)
	id = fmt.Sprintf("%03s", id)

	return id
}
