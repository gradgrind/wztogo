package w365

import (
	"fmt"
	"testing"
	// "regexp"
	//"gradgrind/wztogo/internal/wzbase"
	//"github.com/RoaringBitmap/roaring"
)

func TestReadW365(t *testing.T) {
	fmt.Println("\n############## TestReadW365")
	db := ReadW365Raw("../_testdata/test.w365")
	db.ReadYear(db.ActiveYear)
	for _, yeardata := range db.Years {
		fmt.Printf("\n$$$ %#v\n", yeardata)
	}
	db.read_days()
	fmt.Printf("\n§§NodeList: %#v\n", db.NodeList)
	fmt.Printf("\n§§NodeMap: %#v\n", db.NodeMap)
	fmt.Printf("\n§§TableMap: %#v\n", db.TableMap)
	db.read_hours()
	fmt.Printf("\n§§Config: %#v\n", db.Config)
	db.read_subjects()
	db.read_rooms()
	db.read_absences()
	fmt.Printf("\n§§absences: %#v\n", db.absencemap)
	db.read_categories()
	fmt.Printf("\n§§categories: %#v\n", db.categorymap)
	db.read_teachers()
	db.read_groups()
	for i, n := range db.NodeList {
		fmt.Printf("\n§node %4d: %#v\n", i, n)
	}
}

func TestMisc(t *testing.T) {
	fmt.Println("\n############## TestMisc")
	fmt.Println(convert_date("24. 12. 2023"))
	fmt.Println(convert_colour("-16777216"))
	fmt.Println(convert_colour("-47834"))
	fmt.Println(convert_colour("-16000000"))

	ndays := 5
	absence_map := make([]([]int), ndays)
	for i := range ndays {
		absence_map[i] = []int{}
	}
	absence_map[1] = append(absence_map[1], 5)
	fmt.Printf("absence_map: %v\n", absence_map)

	var x int = 5
	xp := &x
	*xp += 1
	fmt.Printf("XXX %d\n", x)
	type sx struct {
		a int
		b []int
	}
	x1 := sx{a: 1}
	fmt.Printf("XXX1a: %+v\n", x1)
	x1.b = append(x1.b, 7)
	fmt.Printf("XXX1b: %+v\n", x1.b)
	x1p := &x1.b
	*x1p = append(*x1p, 8)
	fmt.Printf("XXX1c: %+v\n", x1.b)
	x2 := map[int]*sx{}
	x2[10] = &sx{a: 11}
	fmt.Printf("XXX2a: %+v\n", x2[10])
	x2x := *x2[10]
	fmt.Printf("XXX2b: %+v\n", x2x)
	x2x.a += 1
	fmt.Printf("XXX2c: %+v\n", x2x)
	fmt.Printf("XXX2d: %+v\n", x2[10])
	x2p := x2[10]
	(*x2p).a += 2
	fmt.Printf("XXX2e: %+v\n", x2[10])
	x3 := map[int]int{}
	x3[20] = 21
	x3[20]++
	fmt.Printf("XXX3: %+v\n", x3)
	x4 := map[int]*sx{}
	x4[30] = &sx{a: 41}
	x4[30].a++
	fmt.Printf("XXX4a: %+v\n", x4[30])
	x4p := x4[30]
	x4p.a++
	fmt.Printf("XXX4b: %+v\n", x4[30])
	x4pp := &(x4p.a)
	*x4pp++
	fmt.Printf("XXX4c: %+v\n", x4[30])

	fmt.Printf("??? %#v\n", *new(string))

	il := map[string][]int{}
	il["A1"] = append(il["A1"], 1)
	il["A1"] = append(il["A1"], 2)
	il["A2"] = append(il["A2"], 3)
	fmt.Printf("\n il: %+v\n", il)
}
