package linker

import (
	"fmt"
	"log"
)

type Linker struct {
	*PointPair
}

func NewLinker(pp *PointPair) *Linker {
	return &Linker{PointPair: pp}
}

func (l *Linker) getPoint(from string) *Point {
	var point *Point
	switch from {
	case "start":
		point = l.Start
	case "end":
		point = l.End
	}
	return point
}

// 一划：横
func (l *Linker) CanLinkInSameLineAxis() bool {
	if isSamePoint := l.EqualPoint(); isSamePoint {
		return true
	}
	currentPoint := l.getPoint("start")
	for {
		nextPoint, err := currentPoint.Right()
		if err != nil {
			log.Fatal(err)
		}
		if nextPoint.RightThen(l.End) {
			log.Fatal("------ move Right Over Then End Point", nextPoint)
		}
		// 移动之后，马上进行判断:
		arrived := EqualPoint(nextPoint, l.End)
		if arrived {
			// 两点重合,当其中一点为空 或 两点typeCode相等时:
			pointIsEmpty := nextPoint.IsEmpty() || currentPoint.IsEmpty()
			currentPointEqualThenEndPoint := EqualTypeCode(currentPoint, l.End)
			ok := pointIsEmpty || currentPointEqualThenEndPoint
			return ok
		} else {
			if nextPoint.IsEmpty() {
				currentPoint = nextPoint
			} else {
				return false
			}
		}
	}
}

// 一划：竖
func (l *Linker) CanLinkInSameRowAxis() bool {
	if isSamePoint := l.EqualPoint(); isSamePoint {
		return true
	}
	currentPoint := l.getPoint("start")
	for {
		nextPoint, err := currentPoint.Down()
		if err != nil {
			log.Fatal(err)
		}
		if nextPoint.UnderThen(l.End) {
			log.Fatal("------  move Bottom Over Then End Point", nextPoint)
		}
		arrived := EqualPoint(nextPoint, l.End)
		if arrived {
			pointIsEmpty := nextPoint.IsEmpty() || currentPoint.IsEmpty()
			currentPointEqualThenEndPoint := EqualTypeCode(currentPoint, l.End)
			ok := pointIsEmpty || currentPointEqualThenEndPoint
			return ok
		} else {
			if nextPoint.IsEmpty() {
				currentPoint = nextPoint
			} else {
				return false
			}
		}
	}
}

// 一划
func (l *Linker) CanLinkInOneStroke() bool {
	var canLink bool
	switch {
	case l.Start.RowIdx == l.End.RowIdx:
		canLink = l.CanLinkInSameLineAxis()
	case l.Start.LineIdx == l.End.LineIdx:
		canLink = l.CanLinkInSameRowAxis()
	}
	return canLink
}

/*
二划：有可能是A-C-B,有可能是A-D—B

       (X1,Y2)    (X2,Y2)
        C┌────────┐B
         │        │
         │        │
         │        │
        A└────────┘D
      (X1,Y1)     (X2,Y1)
*/
func (l *Linker) CanLinkInTwoStrokes() bool {
	table := GetTable()

	PointA := l.Start
	PointB := l.End

	PointC, err1 := table.GetPoint(PointB.RowIdx, PointA.LineIdx)
	PointD, err2 := table.GetPoint(PointA.RowIdx, PointB.LineIdx)
	if err1 != nil || err2 != nil {
		fmt.Println(err1, err2)
		return false
	}
	LinkAToC := NewLinker(NewPointPair(PointC, PointA))
	LinkCToB := NewLinker(NewPointPair(PointC, PointB))

	ALinkC := LinkAToC.CanLinkInSameRowAxis()
	CLinkB := LinkCToB.CanLinkInSameLineAxis()
	if ALinkC && CLinkB {
		return true
	}
	LinkAToD := NewLinker(NewPointPair(PointD, PointA))
	LinkDToB := NewLinker(NewPointPair(PointD, PointB))

	ALinkD := LinkAToD.CanLinkInSameLineAxis()
	DLinkB := LinkDToB.CanLinkInSameRowAxis()
	if ALinkD && DLinkB {
		return true
	}
	return false
}

// 获取 EndPoint为零点的坐标轴上 的所有可抵达点
func (l *Linker) GetEndPointCanReachPointsOnAxis() []*Point {
	var ret []*Point
	end := l.End

	collectCanReachPoints := func(direction int) {
		current := end
		for {
			newPoint, err := current.GetNextPoint(direction)
			if err != nil {
				log.Fatal(err)
			}
			current = newPoint
			if current.IsEmpty() {
				ret = append(ret, current)
			} else {
				break
			}
			if current.AtBorder() {
				break
			}
		}
	}
	collectCanReachPoints(DirectionRight)
	collectCanReachPoints(DirectionLeft)
	collectCanReachPoints(DirectionUp)
	collectCanReachPoints(DirectionDown)
	return ret
}

/*
三划：有可能是A-B1-P1-P,有可能是A-B2—P2-P

                    ^
                    │ P3
                    │
              P1    │ P  P2
        ──────+─────+────+───> (动态Pn点,P1-P3为P的可及范围)
              │     │    |
              │          |
      A───────B1---------B2'

	只要A能在二划之内抵达P1-P2的任意一个点(P为零点的坐标轴上的任意一点),A就一定能在三划之内抵达P
*/
func (l *Linker) CanLinkInThreeStrokes() bool {
	Points := l.GetEndPointCanReachPointsOnAxis()
	for _, PointPn := range Points {
		PointA := l.Start
		AToPn := NewLinker(NewPointPair(PointA, PointPn))
		if AToPn.CanLinkInTwoStrokes() {
			return true
		}
	}
	return false
}

func (l *Linker) TestLink() (canLink bool) {
	if !l.TypeCodeEqual() || l.Start.IsEmpty() || l.End.IsEmpty() {
		return false
	}

	inSameAxis := l.InSameAxis()
	if inSameAxis {
		canLink = l.CanLinkInOneStroke()
	} else {
		canLink = l.CanLinkInTwoStrokes()
	}

	if !canLink {
		canLink = l.CanLinkInThreeStrokes()
	}
	return
}
