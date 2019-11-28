package main

import (
	"fmt"
	"log"

	core "cipherassets.core"
	"cipherassets.core/rest"
)

func main() {
	fmt.Printf("cipherassets core\n")

	c, err := core.NewConfig()
	if err != nil {
		log.Fatalf("can't create config: %s", err.Error())
		return
	}

	r, err := rest.NewREST(c)
	if err != nil {
		log.Fatalf("can't create REST server: %s", err.Error())
		return
	}
	log.Fatalf("REST server shutdown: %s", r.Serve().Error())
}
