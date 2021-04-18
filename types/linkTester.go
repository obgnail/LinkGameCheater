package types

import (
	"fmt"
	"log"
)

type LinkTester struct {
	*PointPair
}

func NewLinkTester(pp *PointPair) *LinkTester {
	return &LinkTester{PointPair: pp}
}

func (lt *LinkTester) getPathfinder(from string) *Point {
	var point *Point
	switch from {
	case "start":
		point = lt.Start
	case "end":
		point = lt.End
	}
	return point
}

// 一划：横
func (lt *LinkTester) CanLinkInSameLineAxis() bool {
	if isSamePoint := lt.Start.Equal(lt.End); isSamePoint {
		return true
	}
	currentPoint := lt.getPathfinder("start")
	for {
		nextPoint, err := currentPoint.Right()
		if err != nil {
			log.Fatal(err)
		}
		if nextPoint.RightThen(lt.End) {
			log.Fatal("------ move Right Over Then End Point", nextPoint)
		}
		// 移动之后，马上进行判断:
		arrived := nextPoint.Equal(lt.End)
		if arrived {
			pointIsEmpty := nextPoint.isEmpty() || currentPoint.isEmpty()
			currentPointEqualThenEndPoint := currentPoint.EqualTypeCode(lt.End)
			if pointIsEmpty || currentPointEqualThenEndPoint {
				return true
			} else {
				return false
			}
		} else {
			if nextPoint.isEmpty() {
				currentPoint = nextPoint
			} else {
				return false
			}
		}
	}
}

// 一划：竖
func (lt *LinkTester) CanLinkInSameRowAxis() bool {
	if isSamePoint := lt.Start.Equal(lt.End); isSamePoint {
		return true
	}
	currentPoint := lt.getPathfinder("start")
	for {
		nextPoint, err := currentPoint.Down()
		if err != nil {
			log.Fatal(err)
		}
		if nextPoint.UnderThen(lt.End) {
			log.Fatal("------  move Bottom Over Then End Point", nextPoint)
		}
		arrived := nextPoint.Equal(lt.End)
		if arrived {
			pointIsEmpty := nextPoint.isEmpty() || currentPoint.isEmpty()
			currentPointEqualThenEndPoint := currentPoint.EqualTypeCode(lt.End)
			if pointIsEmpty || currentPointEqualThenEndPoint {
				return true
			} else {
				return false
			}
		} else {
			if nextPoint.isEmpty() {
				currentPoint = nextPoint
			} else {
				return false
			}
		}
	}
}

// 一划
func (lt *LinkTester) CanLinkInOneStroke() bool {
	var canLink bool
	switch {
	case lt.Start.RowIdx == lt.End.RowIdx:
		canLink = lt.CanLinkInSameLineAxis()
	case lt.Start.LineIdx == lt.End.LineIdx:
		canLink = lt.CanLinkInSameRowAxis()
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
func (lt *LinkTester) CanLinkInTwoStrokes() bool {
	PointA := lt.Start
	PointB := lt.End

	PointC, err1 := Table.GetPoint(PointB.RowIdx, PointA.LineIdx)
	PointD, err2 := Table.GetPoint(PointA.RowIdx, PointB.LineIdx)
	if err1 != nil || err2 != nil {
		fmt.Println(err1, err2)
		return false
	}
	LinkAToC := NewLinkTester(NewPointPair(PointC, PointA))
	LinkCToB := NewLinkTester(NewPointPair(PointC, PointB))

	ALinkC := LinkAToC.CanLinkInSameRowAxis()
	CLinkB := LinkCToB.CanLinkInSameLineAxis()
	if ALinkC && CLinkB {
		return true
	}
	LinkAToD := NewLinkTester(NewPointPair(PointD, PointA))
	LinkDToB := NewLinkTester(NewPointPair(PointD, PointB))

	ALinkD := LinkAToD.CanLinkInSameLineAxis()
	DLinkB := LinkDToB.CanLinkInSameRowAxis()
	if ALinkD && DLinkB {
		return true
	}
	return false
}

// 获取 EndPoint为零点的坐标轴上 的所有可抵达点
func (lt *LinkTester) GetEndPointCanReachPointsOnAxis() []*Point {
	var ret []*Point
	end := lt.End

	collectCanReachPoints := func(direction string) {
		current := end
		for {
			newPoint, err := current.Direction(direction)
			if err != nil {
				log.Fatal(err)
			}
			current = newPoint
			if current.isEmpty() {
				ret = append(ret, current)
			} else {
				break
			}
			if current.AtBorder() {
				break
			}
		}
	}
	collectCanReachPoints("right")
	collectCanReachPoints("left")
	collectCanReachPoints("up")
	collectCanReachPoints("down")
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
func (lt *LinkTester) CanLinkInThreeStrokes() bool {
	Points := lt.GetEndPointCanReachPointsOnAxis()
	for _, PointPn := range Points {
		PointA := lt.Start
		AToPn := NewLinkTester(NewPointPair(PointA, PointPn))
		if AToPn.CanLinkInTwoStrokes() {
			return true
		}
	}
	return false
}

func (lt *LinkTester) CanLink() (canLink bool) {
	if !lt.TypeCodeEqual() || lt.Start.isEmpty() || lt.End.isEmpty() {
		return false
	}

	inSameAxis := lt.InSameAxis()
	if inSameAxis {
		canLink = lt.CanLinkInOneStroke()
	} else {
		canLink = lt.CanLinkInTwoStrokes()
	}

	if !canLink {
		canLink = lt.CanLinkInThreeStrokes()
	}
	return
}
