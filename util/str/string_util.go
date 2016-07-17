package str

import (
	"bytes"
	"fmt"
)

func Concat(args ...interface{}) string {
	var buffer bytes.Buffer

	for _, a := range args {
		buffer.WriteString(fmt.Sprintf("%v", a))
	}

	return buffer.String()
}
