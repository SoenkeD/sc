package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func UserConfirm(msg string) (bool, bool) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(msg + " ")

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return false, false
	}

	input = strings.TrimSpace(input)

	return input == "y", input == "s"
}
