package image

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sort"

	"github.com/obgnail/LinkGameCheater/config"
	"github.com/rivo/duplo"
)

type Image struct {
	*image.NRGBA
}

func NewImage(OriginImagePath string, x0, y0, x1, y1 int) (*Image, error) {
	originImg, err := OpenImage(OriginImagePath)
	if err != nil {
		return nil, err
	}
	img := &Image{originImg}

	// 修剪
	if x0 != -1 && y0 != -1 && x1 != -1 && y1 != -1 {
		newImg := img.GetSubImage(x0, y0, x1, y1)
		return &Image{newImg}, nil
	} else {
		return img, nil
	}
}

func (i *Image) GetSubImage(x0, y0, x1, y1 int) *image.NRGBA {
	img := i.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
	newImg := image.NewNRGBA(image.Rect(0, 0, x1-x0, y1-y0))
	for y := 0; y < y1-y0; y++ {
		for x := 0; x < x1-x0; x++ {
			r32, g32, b32, a32 := img.At(x0+x, y0+y).RGBA()
			r := uint8(r32 >> 8)
			g := uint8(g32 >> 8)
			b := uint8(b32 >> 8)
			a := uint8(a32 >> 8)
			newImg.SetNRGBA(x, y, color.NRGBA{R: r, G: g, B: b, A: a})
		}
	}
	return newImg
}

func (i *Image) GetSubImagesByCount(SubImageRowCount int, SubImageLineCount int) ([][]*image.NRGBA, error) {
	// NOTE: Max.X的X和Table的X坐标轴不一样
	originDW := i.Rect.Max.X
	originDH := i.Rect.Max.Y
	subDW := originDW / SubImageLineCount
	subDH := originDH / SubImageRowCount

	images := make([][]*image.NRGBA, SubImageRowCount)
	for row := 0; row < SubImageRowCount; row++ {
		images[row] = make([]*image.NRGBA, SubImageLineCount)
		for line := 0; line < SubImageLineCount; line++ {
			x0 := line * subDW
			y0 := row * subDH
			x1 := x0 + subDW
			y1 := y0 + subDH
			// 多裁剪12个像素, 减少边框的影响
			subImage := i.GetSubImage(x0+6, y0+6, x1-6, y1-6)
			images[row][line] = subImage
		}
	}
	return images, nil
}

func (i *Image) GetSubImagesByPixel(subImgDW int, subImgDH int) ([][]*image.NRGBA, error) {
	originDW := i.Rect.Max.X
	originDH := i.Rect.Max.Y

	var images [][]*image.NRGBA
	for dh := 0; dh < originDH; dh += subImgDH {
		var imgs []*image.NRGBA
		for dw := 0; dw < originDW; dw += subImgDW {
			x0 := dw
			y0 := dh
			x1 := dw + subImgDW
			y1 := dh + subImgDH
			subImage := i.GetSubImage(x0+6, y0+6, x1-6, y1-6)
			imgs = append(imgs, subImage)
		}
		images = append(images, imgs)
	}
	return images, nil
}

func (i *Image) Save(name string) error {
	create, _ := os.Create(name)
	err := png.Encode(create, i)
	if err != nil {
		return err
	}
	return nil
}

// 用户自定义,如果单纯使用索引还不够表示,可以手动实现该方法,返回用于参考的表示空格的子图片
// 在实现完此方法后,imageArr中,所有与emptyImages相似的图片都会被视为空格
func findEmptyByFunc(imageArr [][]*image.NRGBA) (emptyImages []*image.NRGBA) {
	return nil
}

type Idx struct {
	Row  int
	Line int
}

func NewIndies(indies [][2]int) []*Idx {
	ret := make([]*Idx, len(indies))
	for i, idx := range indies {
		ret[i] = &Idx{Row: idx[0], Line: idx[1]}
	}
	return ret
}

// 输入用于参考的表示空格的子图片
func findEmptyByIdx(imageArr [][]*image.NRGBA, emptyImagesIdx []*Idx) (emptyImages []*image.NRGBA) {
	for _, idx := range emptyImagesIdx {
		rowIdx := idx.Row
		lineIdx := idx.Line

		if rowIdx != -1 && lineIdx != -1 {
			emptyImages = append(emptyImages, imageArr[rowIdx][lineIdx])
		}
	}
	return
}

func GenTableArrByImages(imageArr [][]*image.NRGBA, emptyImagesIdx []*Idx) ([][]int, error) {
	imgStore := duplo.New()
	MapStoreIDToTypeCode := make(map[string]int)
	typeCode := 1

	if emptyImgsFromIdx := findEmptyByIdx(imageArr, emptyImagesIdx); emptyImgsFromIdx != nil {
		for idx, img := range emptyImgsFromIdx {
			imageHash, _ := duplo.CreateHash(img)
			ID := fmt.Sprintf("empty-fromIndex-%d", idx)
			imgStore.Add(ID, imageHash)
			MapStoreIDToTypeCode[ID] = config.PointTypeCodeEmpty
		}
	}
	if emptyImgsFromFunc := findEmptyByFunc(imageArr); emptyImgsFromFunc != nil {
		for idx, img := range emptyImgsFromFunc {
			imageHash, _ := duplo.CreateHash(img)
			ID := fmt.Sprintf("empty-fromFunction-%d", idx)
			imgStore.Add(ID, imageHash)
			MapStoreIDToTypeCode[ID] = config.PointTypeCodeEmpty
		}
	}

	rowLen := len(imageArr)
	lineLen := len(imageArr[0])
	table := make([][]int, len(imageArr))
	for rowIdx := 0; rowIdx < rowLen; rowIdx++ {
		table[rowIdx] = make([]int, lineLen)
		for lineIdx := 0; lineIdx < lineLen; lineIdx++ {
			if config.SaveSubImage {
				i := &Image{imageArr[rowIdx][lineIdx]}
				if err := i.Save(fmt.Sprintf("%s/%d-%d.png", config.SubImagePath, rowIdx, lineIdx)); err != nil {
					return nil, err
				}
			}
			var curImgTypeCode int
			toCompareImg := imageArr[rowIdx][lineIdx]
			imageHash, _ := duplo.CreateHash(toCompareImg)
			matches := imgStore.Query(imageHash)
			sort.Sort(matches)
			// 新的typeCode
			if len(matches) == 0 || matches[0].Score > config.SimilarScore {
				curImgTypeCode = typeCode
				ID := fmt.Sprintf("%d", typeCode)
				MapStoreIDToTypeCode[ID] = typeCode
				imgStore.Add(ID, imageHash)
				typeCode++
			} else {
				curImgTypeCode = MapStoreIDToTypeCode[matches[0].ID.(string)]
			}
			table[rowIdx][lineIdx] = curImgTypeCode
		}
	}
	return table, nil
}

func OpenImage(OriginImagePath string) (*image.NRGBA, error) {
	f, err := os.Open(OriginImagePath)
	if err != nil {
		return nil, err
	}
	m, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	nrgba := m.(*image.NRGBA)
	return nrgba, nil
}
