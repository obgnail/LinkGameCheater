package main

import (
	"fmt"
	"time"

	"github.com/obgnail/LinkGameCheater/config"
	"github.com/obgnail/LinkGameCheater/types"
)

func run() error {
	t := time.Now()

	types.InitTable(config.GenTableMethod)
	fmt.Println("------ table ------")
	fmt.Println(types.Table)

	game := types.NewGame(types.Table)
	fmt.Println("------ step ------")
	if err := game.Play(); err != nil {
		return err
	}

	elapsed := time.Since(t)

	fmt.Println("\n------ result ------")
	fmt.Println(types.Table)

	fmt.Println("------ elapsed ------")
	fmt.Println("Program elapsed: ", elapsed)

	return nil
}
