package logs

import (
	"testing"
)

func TestAll(t *testing.T) {
	// test Info
	Info("%s", "Info")
	I("info")

	// test warn
	Warn("%s", "Warn")
	W("Warn")

	// test Debug
	Debug("%s", "Debug")
	D("Debug")

	// test Error
	Error("%s", "Error")
	E("Error")

	// test Critical
	Critical("%s", "Critical")
	C("Critical")

}
