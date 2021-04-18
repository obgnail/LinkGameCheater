package types

import (
	"fmt"
)

type Game struct {
	*GameTable
	pointPairs map[*PointPair]byte
}

func NewGame(table *GameTable) *Game {
	g := &Game{GameTable: table}
	g.collectPointPairs()
	return g
}

func (g *Game) collectPointPairs() {
	ret := make(map[*PointPair]byte)
	for _, points := range g.pointTypeMap {
		pps := Compose(points)
		for _, pp := range pps {
			ret[pp] = 1
		}
	}
	g.pointPairs = ret
}

func (g *Game) Play() error {
	step := 1
	for len(g.pointPairs) != 0 {
		hadLinked := false
		for pointPair := range g.pointPairs {
			if pointPair.Start.isEmpty() || pointPair.End.isEmpty() {
				delete(g.pointPairs, pointPair)
				continue
			}
			lt := NewLinkTester(pointPair)
			canLink := lt.CanLink()
			if canLink {
				fmt.Printf(" step %d %s\n", step, pointPair)
				step++
				hadLinked = true
				if err := g.SetEmpty(pointPair.Start.RowIdx, pointPair.Start.LineIdx); err != nil {
					return err
				}
				if err := g.SetEmpty(pointPair.End.RowIdx, pointPair.End.LineIdx); err != nil {
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
