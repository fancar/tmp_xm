package main

import "github.com/fancar/cobra_template/cmd/app/cmd"

/*
	Auhtor:		Mamaev Alexander
	Email:		fancatster@gmail.com
	year:		2022
*/

var version string // set by the compiler

func main() {
	cmd.Execute(version)
}
