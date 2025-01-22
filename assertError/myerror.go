package asserterror

import (
	"fmt"
	"os"
)

func Assert(cond bool, text string, err error) {
	if cond {
		fmt.Println("ERROR - ", text, " -- ", err)
		os.Exit(1)
	}
}
