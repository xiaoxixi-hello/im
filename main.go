package main

import (
	"github.com/ylinyang/im/router"
)

func main() {
	r := router.Router()
	if err := r.Run(); err != nil {
		panic(err)
	}
}
