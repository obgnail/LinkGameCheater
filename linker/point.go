package linker

import (
	"encoding/json"
	"github.com/obgnail/LinkGameCheater/config"
)

const (
	DirectionRight = iota
	DirectionLeft
	DirectionUp
	DirectionDown
)

/*
   +────+───> Y(Line)
   │
   │
   X(Row)
*/
type Point struct {
	RowIdx   int
	LineIdx  int
	TypeCode int
}

// 仅用于创建GameTable,如果需要获取Point,请使用 table.GetPoint()
func NewPoint(rowIdx, lineIdx, typeCode int) *Point {
	p := new(Point)
	p.RowIdx = rowIdx
	p.LineIdx = lineIdx
	p.TypeCode = typeCode
	return p
}

func (p *Point) String() string {
	out, err := json.Marshal(p)
	if err != nil {
		return err.Error()
	}
	return string(out)
}

func (p *Point) isInValid() bool {
	table := GetTable()
	return 0 > p.RowIdx || p.RowIdx >= table.RowLen || 0 > p.LineIdx || p.LineIdx >= table.LineLen
}

func (p *Point) AtBorder() bool {
	table := GetTable()
	return p.RowIdx == 0 || p.LineIdx == 0 || p.RowIdx == table.RowLen-1 || p.LineIdx == table.LineLen-1
}

func (p *Point) IsEmpty() bool {
	return p.TypeCode == config.PointTypeCodeEmpty
}

func (p *Point) RightThen(other *Point) bool {
	return p.LineIdx > other.LineIdx
}

func (p *Point) UnderThen(other *Point) bool {
	return p.RowIdx > other.RowIdx
}

// 获取临近点
func (p *Point) GetNextPoint(direction int) (*Point, error) {
	rowIdx := p.RowIdx
	lineIdx := p.LineIdx
	switch direction {
	case DirectionRight:
		lineIdx++
	case DirectionLeft:
		lineIdx--
	case DirectionUp:
		rowIdx--
	case DirectionDown:
		rowIdx++
	}

	table := GetTable()
	newPoint, err := table.GetPoint(rowIdx, lineIdx)
	if err != nil {
		return nil, err
	}
	return newPoint, nil
}

func (p *Point) Right() (*Point, error) {
	return p.GetNextPoint(DirectionRight)
}

func (p *Point) Left() (*Point, error) {
	return p.GetNextPoint(DirectionLeft)
}

func (p *Point) Up() (*Point, error) {
	return p.GetNextPoint(DirectionUp)
}
func (p *Point) Down() (*Point, error) {
	return p.GetNextPoint(DirectionDown)
}

// 组合
func Compose(points []*Point) []*PointPair {
	var ret []*PointPair
	length := len(points)
	for i := 0; i < length; i++ {
		for j := 0; j < i; j++ {
			pp := NewPointPair(points[i], points[j])
			ret = append(ret, pp)
		}
	}
	return ret
}

func EqualPoint(a, b *Point) bool {
	if a == b {
		return true
	}
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}
	return *a == *b
}

func EqualTypeCode(a, b *Point) bool {
	return a.TypeCode == b.TypeCode
}
