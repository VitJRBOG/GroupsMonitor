package main

import "runtime/debug"

func main() {
	err := RunGui()
	if err != nil {
		ToLogFile(err.Error(), string(debug.Stack()))
		panic(err.Error())
	}
}
