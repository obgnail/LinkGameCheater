package config

const (
	// 代表空格的数字
	PointTypeCodeEmpty = 0
	// 判断图片是否相似的分数(越低越相似)
	SimilarScore = -40
)

// 共有四种生成Table的方法
const (
	// FromRandom or FromArr or FromImageByCount or FromImageByPixel
	GenTableMethod = "FromImageByCount"
)

// 随机生成Table
var (
	// 牌种数
	TypeCodeCount = 8
	// 长
	LineLen = 16
	// 宽
	RowLen = 8
)

// 使用Arr生成Table
var (
	Table = [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 7, 5, 7, 5, 8, 4, 4, 7, 0},
		{0, 3, 5, 5, 8, 6, 1, 7, 8, 0},
		{0, 1, 2, 1, 3, 2, 6, 5, 3, 0},
		{0, 4, 2, 3, 3, 4, 2, 5, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
)

// 使用图片生成Table
var (
	// 原图片位置
	ImagePath = "images/test.png"
	// 修剪原图片(都为-1表示不修剪)
	FixImageRectangleMinPointX = -1
	FixImageRectangleMinPointY = -1
	FixImageRectangleMaxPointX = -1
	FixImageRectangleMaxPointY = -1

	// 是否保存子图片
	SaveSubImage = false
	// 保存子图片的位置
	SubImagePath = "images/sub"

	// 用于参考的 表示空格的子图片 的索引值(即:与这些图片相似的都会被认为是空格)
	// (rowIdx,LineIdx)
	// (-1,-1)表示图片中没有空格
	// 如果单纯使用索引还不够表示,可以手动实现types/images的findEmptyByFunc()方法
	EmptySubImageIndies = [][2]int{
		{-1, -1},
	}

	// 下面配置二选一
	// 通过像素切割图片
	EachSubImageRowPixel  = 131
	EachSubImageLinePixel = 166
	// 通过数量切割图片
	ImageRowCount  = 9
	ImageLineCount = 16
)
