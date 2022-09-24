// Copyright 2022 eightlay (github.com/eightlay). All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	"github.com/eightlay/rummikub-server/iternal/server"
)

const addr string = ":3000"

func main() {
	file, err := os.OpenFile("runtime.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(file)

	server.StartServer(addr)
}
