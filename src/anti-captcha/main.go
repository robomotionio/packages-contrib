package main

import (
	"antiCaptcha/v1"

	"bitbucket.org/mosteknoloji/robomotion-go-lib/runtime"
)

func main() {

	runtime.RegisterNodes(

		&antiCaptcha.ImageCaptcha{},
	)

	runtime.Start()

}
