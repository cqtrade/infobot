/*
Go interfaces generally belong in the package that uses values of the interface type, not the package that implements those values. The implementing package should return concrete (usually pointer or struct) types: that way, new methods can be added to implementations without requiring extensive refactoring.
*/
package main

import (
	app "github.com/cqtrade/infobot/src"
)

func main() {
	app.Run()
}
