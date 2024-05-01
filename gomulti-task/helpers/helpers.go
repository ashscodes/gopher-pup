/*
This package is intended for internal usage only in tests. Do not use any of the functions within.
*/
package helpers

import "os"

// Helper function to disable stdOut.
func DisableStdOut() {
	nullDevice := os.DevNull
	if os.PathSeparator == '\\' {
		// Windows uses 'NUL'
		nullDevice = "NUL"
	}

	os.Stdout, _ = os.Open(nullDevice)
}

// Helper function to re-enable stdOut
func EnableStdOut() {
	os.Stdout = os.NewFile(1, os.DevNull)
}
