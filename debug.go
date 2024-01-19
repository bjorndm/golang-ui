package ui

import "fmt"

var debugDisplay = false

func dprint(args ...any) {
	if debugDisplay {
		fmt.Print(args...)
	}
}

func dprintln(args ...any) {
	if debugDisplay {
		fmt.Println(args...)
	}
}

func dprintf(form string, args ...any) {
	if debugDisplay {
		fmt.Printf(form, args...)
	}
}
