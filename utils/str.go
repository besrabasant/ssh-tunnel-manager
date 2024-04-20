package utils

import "strconv"

func IntToString(num int) string {
    if num == 0 {
        return "" // Return an empty string if the number is 0
    }
    return strconv.Itoa(num) // Convert and return the number as a string otherwise
}


func Bold(str string) string {
	return "\x1b[1m" + str + "\x1b[0m"
}