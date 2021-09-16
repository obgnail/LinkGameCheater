package types

import (
	"fmt"
	"image"
	"log"
	"strings"

	"github.com/obgnail/LinkGameCheater/config"
	"github.com/obgnail/LinkGameCheater/utils"
)

var table *GameTable

type GameTable struct {
	rowLen  int
	lineLen int
	table   [][]*Point

	pointTypeMap map[int][]*Point
}

func NewGameTable(linkGameTable [][]int) *GameTable {
	rowLen := len(linkGameTable)
	lineLen := len(linkGameTable[0])

	t := &GameTable{
		rowLen:       rowLen,
		lineLen:      lineLen,
		table:        make([][]*Point, rowLen),
		pointTypeMap: make(map[int][]*Point),
	}
	for rowIdx := 0; rowIdx < t.rowLen; rowIdx++ {
		t.table[rowIdx] = make([]*Point, lineLen)
		for lineIdx := 0; lineIdx < t.lineLen; lineIdx++ {
			typeCode := linkGameTable[rowIdx][lineIdx]
			point := newPoint(rowIdx, lineIdx, typeCode)
			t.table[rowIdx][lineIdx] = point
			if typeCode != config.PointTypeCodeEmpty {
				t.pointTypeMap[typeCode] = append(t.pointTypeMap[typeCode], point)
			}
		}
	}
	return t
}

func NewTableFromArr(tableArr [][]int) *GameTable {
	return NewGameTable(tableArr)
}

func NewTableFromRandom(typeCodeCount, rowLen, lineLen int) *GameTable {
	total := lineLen * rowLen
	TableList, err := utils.GenRandomTableList(typeCodeCount, total)
	if err != nil {
		log.Fatal("[ERROR] Gen TableList", err)
	}
	TableArr, err := utils.GenTableArr(TableList, lineLen, rowLen)
	if err != nil {
		log.Fatal("[ERROR] Gen TableArr", err)
	}
	table := NewGameTable(TableArr)
	return table
}

func NewTableFromImageByCount(
	imagePath string,
	fixRectangleMinPointX, fixRectangleMinPointY, fixRectangleMaxPointX, fixRectangleMaxPointY int,
	rowLen, lineLen int,
	emptyIndies []*Idx,
) *GameTable {
	img, err := NewImage(imagePath, fixRectangleMinPointX, fixRectangleMinPointY, fixRectangleMaxPointX, fixRectangleMaxPointY)
	if err != nil {
		log.Fatal(err)
	}
	subImages, err := img.GetSubImagesByCount(rowLen, lineLen)
	if err != nil {
		log.Fatal(err)
	}

	table, err := NewTableByImageArr(subImages, emptyIndies)
	if err != nil {
		log.Fatal(err)
	}
	return table
}

func NewTableFromImageByPixel(
	imagePath string,
	fixRectangleMinPointX, fixRectangleMinPointY, fixRectangleMaxPointX, fixRectangleMaxPointY int,
	subImgDW, subImgDH int,
	emptyIndies []*Idx,
) *GameTable {
	img, err := NewImage(imagePath, fixRectangleMinPointX, fixRectangleMinPointY, fixRectangleMaxPointX, fixRectangleMaxPointY)
	if err != nil {
		log.Fatal(err)
	}
	subImages, err := img.GetSubImagesByPixel(subImgDW, subImgDH)
	if err != nil {
		log.Fatal(err)
	}

	table, err := NewTableByImageArr(subImages, emptyIndies)
	if err != nil {
		log.Fatal(err)
	}
	return table
}

func NewTableByImageArr(imageArr [][]*image.NRGBA, emptyIndies []*Idx) (*GameTable, error) {
	linkGameTable, err := GenTableArrByImages(imageArr, emptyIndies)
	if err != nil {
		return nil, err
	}
	withEmpty := utils.AddOutEmptyPoint(linkGameTable)
	table := NewGameTable(withEmpty)
	return table, nil
}

func (t *GameTable) String() string {
	var rows []string
	for rowIdx := 0; rowIdx < t.rowLen; rowIdx++ {
		var line []string
		for lineIdx := 0; lineIdx < t.lineLen; lineIdx++ {
			point := t.table[rowIdx][lineIdx]
			s := fmt.Sprintf("%d", point.TypeCode)
			line = append(line, s)
		}
		rows = append(rows, strings.Join(line, "\t"))
	}
	return strings.Join(rows, "\n") + "\n"
}

func (t *GameTable) GetPoint(rowIdx, lineIdx int) (*Point, error) {
	if 0 > rowIdx || rowIdx >= t.rowLen || 0 > lineIdx || lineIdx >= t.lineLen {
		return nil, fmt.Errorf("point(%d, %d) is out of boundary(%d, %d)", rowIdx, lineIdx, t.rowLen, t.lineLen)
	}
	return t.table[rowIdx][lineIdx], nil
}

func (t *GameTable) SetEmpty(rowIdx, lineIdx int) error {
	p, err := t.GetPoint(rowIdx, lineIdx)
	if err != nil {
		return err
	}
	p.TypeCode = config.PointTypeCodeEmpty
	return nil
}

func InitTable(method string) {
	switch method {
	case "FromRandom":
		table = NewTableFromRandom(config.TypeCodeCount, config.LineLen, config.RowLen)
	case "FromArr":
		table = NewTableFromArr(config.Table)
	case "FromImageByCount":
		emptyIndies := NewIndies(config.EmptySubImageIndies)
		table = NewTableFromImageByCount(
			config.ImagePath,
			config.FixImageRectangleMinPointX,
			config.FixImageRectangleMinPointY,
			config.FixImageRectangleMaxPointX,
			config.FixImageRectangleMaxPointY,
			config.ImageRowCount,
			config.ImageLineCount,
			emptyIndies,
		)
	case "FromImageByPixel":
		emptyIndies := NewIndies(config.EmptySubImageIndies)
		table = NewTableFromImageByPixel(
			config.ImagePath,
			config.FixImageRectangleMinPointX,
			config.FixImageRectangleMinPointY,
			config.FixImageRectangleMaxPointX,
			config.FixImageRectangleMaxPointY,
			config.EachSubImageRowPixel,
			config.EachSubImageLinePixel,
			emptyIndies,
		)
	default:
		log.Fatal("ERROR Method")
	}
}

func GetTable() *GameTable {
	return table
}
