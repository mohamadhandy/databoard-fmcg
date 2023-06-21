package helper

import (
	"fmt"
	"strconv"
	"strings"
)

func GenerateNextID(highestID int) string {
	// Generate id for brand
	// Increment the highest ID by 1
	nextID := highestID + 1

	// Convert the incremented ID back to a string with leading zeros
	id := strconv.Itoa(nextID)
	id = fmt.Sprintf("%03s", id)

	return id
}

func GenerateNextIDCategory(highestID int) string {
	// Generate id for Category
	// Increment the highest ID by 1
	nextID := highestID + 1

	// Convert the incremented ID back to a string with leading zeros
	id := strconv.Itoa(nextID)
	id = fmt.Sprintf("%02s", id)

	return id
}

func GenerateProductID(highestID int) (string, string) {
	// Generate id for Product Id
	// increment the highest ID by 1
	nextId := highestID + 1
	// convert the incremented id back to a string with leading zeros
	id := strconv.Itoa(nextId)
	id = fmt.Sprintf("%04s", id)
	return "KDR" + id, id
}

func SplitProductID(id string) string {
	// split KDR & id
	parts := strings.Split(id, "KDR")
	onlyId := ""
	if len(parts) > 1 {
		onlyId = parts[1]
	}
	return onlyId
}
