package api

import "strconv"

// isValidID проверяет, что id имеет корректный формат
func IsValidID(id string) bool {
	_, err := strconv.Atoi(id)
	return err == nil
}
