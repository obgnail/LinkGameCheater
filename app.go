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
	table := types.GetTable()
	fmt.Println("------ table ------")
	fmt.Println(table)

	game := types.NewGame(table)
	fmt.Println("------ step ------")
	if err := game.Play(); err != nil {
		return err
	}

	elapsed := time.Since(t)

	fmt.Println("\n------ result ------")
	fmt.Println(table)

	fmt.Println("------ elapsed ------")
	fmt.Println("Program elapsed: ", elapsed)

	return nil
}
