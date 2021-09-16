package types

import (
	"encoding/json"
	"fmt"

	"github.com/obgnail/LinkGameCheater/config"
)

type Point struct {
	RowIdx   int
	LineIdx  int
	TypeCode int
}

// 仅用于创建GameTable,如果需要获取Point,请使用 Table.GetPoint()
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
	return 0 > p.RowIdx || p.RowIdx >= Table.rowLen || 0 > p.LineIdx || p.LineIdx >= Table.lineLen
}

func (p *Point) AtBorder() bool {
	return p.RowIdx == 0 || p.LineIdx == 0 || p.RowIdx == Table.rowLen-1 || p.LineIdx == Table.lineLen-1
}

func (p *Point) isEmpty() bool {
	return p.TypeCode == config.PointTypeCodeEmpty
}

func (p *Point) Equal(other *Point) bool {
	if p == other {
		return true
	}

	if (p == nil && other != nil) || (p != nil && other == nil) {
		return false
	}

	return *p == *other
}

func (p *Point) EqualTypeCode(other *Point) bool {
	return p.TypeCode == other.TypeCode
}

func (p *Point) RightThen(other *Point) bool {
	return p.LineIdx > other.LineIdx
}

func (p *Point) UnderThen(other *Point) bool {
	return p.RowIdx > other.RowIdx
}

// 获取临近点
func (p *Point) Direction(direction string) (*Point, error) {
	rowIdx := p.RowIdx
	lineIdx := p.LineIdx
	switch direction {
	case "right":
		lineIdx++
	case "left":
		lineIdx--
	case "up":
		rowIdx--
	case "down":
		rowIdx++
	}
	if p.isInValid() {
		return nil, fmt.Errorf("point(%d, %d) is out of boundary(%d, %d)", rowIdx, lineIdx, Table.rowLen, Table.lineLen)
	}
	newPoint, err := Table.GetPoint(rowIdx, lineIdx)
	if err != nil {
		return nil, err
	}
	return newPoint, nil
}

func (p *Point) Right() (*Point, error) {
	return p.Direction("right")
}

func (p *Point) Left() (*Point, error) {
	return p.Direction("left")
}

func (p *Point) Up() (*Point, error) {
	return p.Direction("up")
}
func (p *Point) Down() (*Point, error) {
	return p.Direction("down")
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
