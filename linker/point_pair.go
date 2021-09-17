package linker

import (
	"fmt"
)

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
		start, end = p1, p2
	} else {
		start, end = p2, p1
	}
	return &PointPair{Start: start, End: end}
}

func (pp *PointPair) GetStartPoint() *Point {
	return pp.Start
}

func (pp *PointPair) GetEndPoint() *Point {
	return pp.End
}

// 同轴
func (pp *PointPair) InSameAxis() bool {
	return pp.Start.RowIdx == pp.End.RowIdx || pp.Start.LineIdx == pp.End.LineIdx
}

func (pp *PointPair) TypeCodeEqual() bool {
	return EqualTypeCode(pp.Start, pp.End)
}

func (pp *PointPair) EqualPoint() bool {
	return EqualPoint(pp.Start, pp.End)
}
