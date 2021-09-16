package cheater

import (
	"fmt"

	"github.com/obgnail/LinkGameCheater/linker"
)

type Cheater struct {
	table      *linker.GameTable
	pointPairs map[*linker.PointPair]byte
}

func NewGame(table *linker.GameTable) *Cheater {
	g := &Cheater{table: table}
	g.collectPointPairs()
	return g
}

func (c *Cheater) collectPointPairs() {
	ret := make(map[*linker.PointPair]byte)
	for _, points := range c.table.PointTypeMap {
		pps := linker.Compose(points)
		for _, pp := range pps {
			ret[pp] = 1
		}
	}
	c.pointPairs = ret
}

func (c *Cheater) Play() error {
	step := 1
	for len(c.pointPairs) != 0 {
		hadLinked := false
		for pointPair := range c.pointPairs {
			if pointPair.Start.IsEmpty() || pointPair.End.IsEmpty() {
				delete(c.pointPairs, pointPair)
				continue
			}
			lt := linker.NewLinkTester(pointPair)
			canLink := lt.CanLink()
			if canLink {
				fmt.Printf(" step %d %s\n", step, pointPair)
				step++
				hadLinked = true
				if err := c.table.SetEmpty(pointPair.Start.RowIdx, pointPair.Start.LineIdx); err != nil {
					return err
				}
				if err := c.table.SetEmpty(pointPair.End.RowIdx, pointPair.End.LineIdx); err != nil {
					return err
				}
			}
		}
		if !hadLinked {
			break
		}
	}
	return nil
}
