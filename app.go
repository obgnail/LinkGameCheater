package main

import (
	"fmt"
	"time"

	"github.com/obgnail/LinkGameCheater/cheater"
	"github.com/obgnail/LinkGameCheater/config"
	"github.com/obgnail/LinkGameCheater/linker"
)

func run() error {
	t := time.Now()

	linker.InitTable(config.GenTableMethod)
	table := linker.GetTable()
	fmt.Println("------ table ------")
	fmt.Println(table)

	c := cheater.NewGame(table)
	fmt.Println("------ step ------")
	if err := c.Play(); err != nil {
		return err
	}

	elapsed := time.Since(t)

	fmt.Println("\n------ result ------")
	fmt.Println(table)

	fmt.Println("------ elapsed ------")
	fmt.Println("Program elapsed: ", elapsed)

	return nil
}
