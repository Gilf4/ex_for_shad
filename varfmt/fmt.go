//go:build !solution
package varfmt

import (
	"fmt"
	"strings"
	"strconv"
)

// Sprintf форматирует строку с использованием переменного числа аргументов.
func Sprintf(format string, args ...interface{}) string {
	var result strings.Builder
	bracketIndex := 0

	for i := 0; i < len(format); i++ {
		if format[i] == '{' {
			j := i + 1
			for format[j] != '}'{
				j++
			}
			index := 0
			if i+1 == j {
				index = bracketIndex
			} else {
				var err error
				index, err = strconv.Atoi(format[i+1 : j])
				if err != nil{
					fmt.Println(err)
				}
			}

			arg := fmt.Sprint(args[index])
			result.WriteString(arg)

			bracketIndex++
			i = j
		} else {
			result.WriteString(string(format[i]))
		}
	}
	return result.String()
}