//Util package of Leolife project.
//Contains utility methods.
package util

import (
	"fmt"
	"os"
)

func GenerateUuid() (string, error) {
	var random *os.File
	var err error

	if random, err = os.Open("/dev/urandom"); err == nil {
		b := make([]byte, 16)
		random.Read(b)
		return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
	} else {
		return "", err
	}
}
