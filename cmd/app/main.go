package main

import "github.com/fancar/tmp_xm/cmd/app/cmd"

/*
	Auhtor:		Mamaev Alexander
	Email:		fancatster@gmail.com
	year:		2023
*/

var version string // set by the compiler

func main() {
	cmd.Execute(version)
}
