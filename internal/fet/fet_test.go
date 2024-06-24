package fet

import (
	"fmt"
	"gradgrind/wztogo/internal/w365"
	"gradgrind/wztogo/internal/wzbase"
	"testing"
)

//func TestDays(t *testing.T) {
//	readDays()
//}

func TestFet(t *testing.T) {
	w365file := "../_testdata/fms.w365"
	// w365file := "../_testdata/test.w365"
	wzdb := w365.ReadW365(w365file)

	fmt.Println("\n *******************************************")
	//fmt.Printf("\n Class_Groups: %+v\n", wzdb.AtomicGroups.Class_Groups)
	fmt.Printf("\n Classes: %+v\n", wzdb.TableMap["CLASSES"])
	for _, c := range wzdb.TableMap["CLASSES"] {
		// ???
		ag_gs := map[int][]string{}

		cgs := wzdb.AtomicGroups.Class_Groups[c]
		cags := wzdb.AtomicGroups.Group_Atomics[wzbase.ClassGroup{
			CIX: c, GIX: 0,
		}]
		fmt.Printf("\n Class %s: %+v\n",
			wzdb.NodeList[c].Node.(wzbase.Class).ID,
			cags,
		)
		for _, cg := range cgs {
			g := wzdb.NodeList[cg.GIX].Node.(wzbase.Group).ID
			ags := wzdb.AtomicGroups.Group_Atomics[cg]
			fmt.Printf("  ++ %s: %+v\n", g, ags)
			for _, ag := range ags.ToArray() {
				ag_gs[int(ag)] = append(ag_gs[int(ag)], g)
			}
		}
		for ag, gs := range ag_gs {
			fmt.Printf("  ** %d: %+v\n", ag, gs)
		}
	}

	return

	//fmt.Printf("\nINPUT: %+v\n", wzdb)
	xmlitem := getDays(&wzdb)
	fmt.Printf("\n*** fet:\n%v\n", xmlitem)
	xmlitem = getHours(&wzdb)
	fmt.Printf("\n*** fet:\n%v\n", xmlitem)
	xmlitem = getSubjects(&wzdb)
	fmt.Printf("\n*** fet:\n%v\n", xmlitem)
	xmlitem = getTeachers(&wzdb)
	fmt.Printf("\n*** fet:\n%v\n", xmlitem)
	xmlitem = getClasses(&wzdb)
	fmt.Printf("\n*** fet:\n%v\n", xmlitem)
}
