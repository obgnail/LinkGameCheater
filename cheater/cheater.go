package cheater

import (
	"fmt"

	"github.com/obgnail/LinkGameCheater/config"
	"github.com/obgnail/LinkGameCheater/linker"
)

type Cheater struct {
	table      *linker.Table
	pointPairs map[*linker.PointPair]struct{}
}

func NewGame(table *linker.Table) *Cheater {
	g := &Cheater{table: table}
	g.collectPointPairs(table)
	return g
}

func (c *Cheater) collectPointPairs(table *linker.Table) {
	pointTypeMap := make(map[int][]*linker.Point)

	for rowIdx := 0; rowIdx < table.RowLen; rowIdx++ {
		for lineIdx := 0; lineIdx < table.LineLen; lineIdx++ {
			point := table.Table[rowIdx][lineIdx]
			typeCode := point.TypeCode
			if typeCode != config.PointTypeCodeEmpty {
				pointTypeMap[typeCode] = append(pointTypeMap[typeCode], point)
			}
		}
	}

	ret := make(map[*linker.PointPair]struct{})
	for _, points := range pointTypeMap {
		pps := linker.Compose(points)
		for _, pp := range pps {
			ret[pp] = struct{}{}
		}
	}
	c.pointPairs = ret
}

func (c *Cheater) removePointPairs(pointPair *linker.PointPair) {
	delete(c.pointPairs, pointPair)
}

func (c *Cheater) Play() error {
	step := 1
	for len(c.pointPairs) != 0 {
		hadLinked := false
		for pointPair := range c.pointPairs {
			if pointPair.Start.IsEmpty() || pointPair.End.IsEmpty() {
				c.removePointPairs(pointPair)
				continue
			}

			lt := linker.NewLinkTester(pointPair)
			canLink := lt.TestLink()
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
