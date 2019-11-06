package helpers

import (
	"fmt"
	"runtime/debug"
)

// RecoverPanic -- will handle panic and logs instead of crashing
func RecoverPanic() {
	if r := recover(); r != nil {
		fmt.Println("recovered:", r, string(debug.Stack()))
		PushToSlack(fmt.Sprint(r), "panic", string(debug.Stack()))
	}
}
