package types

import (
	"fmt"
)

/*
   +────+───> Y(Line)
   │
   │
   X(Row)
*/
type PointPair struct {
	Start *Point
	End   *Point
}

func (pp *PointPair) String() string {
	start := pp.Start
	end := pp.End
	return fmt.Sprintf("(%d,%d)->(%d,%d)", start.RowIdx, start.LineIdx, end.RowIdx, end.LineIdx)
}

func NewPointPair(p1, p2 *Point) *PointPair {
	var start, end *Point
	if p1.LineIdx < p2.LineIdx || (p1.LineIdx == p2.LineIdx && p1.RowIdx < p2.RowIdx) {
		start = p1
		end = p2
	} else {
		start = p2
		end = p1
	}
	pp := new(PointPair)
	pp.Start = start
	pp.End = end
	return pp
}

// 同轴
func (pp *PointPair) InSameAxis() bool {
	return pp.Start.RowIdx == pp.End.RowIdx || pp.Start.LineIdx == pp.End.LineIdx
}

func (pp *PointPair) TypeCodeEqual() bool {
	return pp.Start.TypeCode == pp.End.TypeCode
}
