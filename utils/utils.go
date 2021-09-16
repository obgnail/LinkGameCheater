package utils

import (
	"errors"
	"math/rand"
	"time"

	"github.com/obgnail/LinkGameCheater/config"
)

func GenRandomTableList(typeCodeCount, totalCount int) ([]int, error) {
	if typeCodeCount <= 0 || totalCount <= 0 {
		return nil, errors.New("<= 0")
	}
	eachTypeCodeCount := totalCount / typeCodeCount

	var tableList []int
	for i := 0; i < eachTypeCodeCount; i++ {
		for j := 1; j < typeCodeCount; j++ {
			tableList = append(tableList, j)
		}
	}
	less := totalCount - len(tableList)

	for i := 0; i < less; i++ {
		tableList = append(tableList, typeCodeCount)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(tableList), func(i, j int) { tableList[i], tableList[j] = tableList[j], tableList[i] })

	return tableList, nil
}

// 在外围添加 PointTypeCodeEmpty
func AddOutEmptyPoint(table [][]int) [][]int {
	rowLen := len(table)
	lineLen := len(table[0])
	newLineLen := lineLen + 2
	newRowLen := rowLen + 2

	newTable := make([][]int, newRowLen)
	for r := 0; r < newRowLen; r++ {
		newTable[r] = make([]int, newLineLen)
		for l := 0; l < newLineLen; l++ {
			newTable[r][l] = config.PointTypeCodeEmpty
		}
	}

	for rowIdx := 0; rowIdx < rowLen; rowIdx++ {
		for LineIdx := 0; LineIdx < lineLen; LineIdx++ {
			newTable[rowIdx+1][LineIdx+1] = table[rowIdx][LineIdx]
		}
	}
	return newTable
}

func GenTableArr(tableList []int, rowLen, lineLen int) ([][]int, error) {
	tableLineLen := lineLen + 2
	tableRowLen := rowLen + 2
	table := make([][]int, tableRowLen)

	// 在外围添加 PointTypeCodeEmpty
	for i := 0; i < tableRowLen; i++ {
		table[i] = make([]int, tableLineLen)
		for k := 0; k < tableLineLen; k++ {
			table[i][k] = config.PointTypeCodeEmpty
		}
	}

	ListIdx := 0
	for rowIdx := 1; rowIdx < rowLen+1; rowIdx++ {
		for lineIdx := 1; lineIdx < lineLen+1; lineIdx++ {
			table[rowIdx][lineIdx] = tableList[ListIdx]
			ListIdx++
		}
	}
	return table, nil
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
