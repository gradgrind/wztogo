package fet

import (
	"fmt"
	"gradgrind/wztogo/internal/w365"
	"testing"
)

//func TestDays(t *testing.T) {
//	readDays()
//}

func TestFet(t *testing.T) {
	// w365file := "../_testdata/fms.w365"
	w365file := "../_testdata/test.w365"
	wzdb := w365.ReadW365(w365file)

	fmt.Println("\n *******************************************")
	/*
		//fmt.Printf("\n Class_Groups: %+v\n", wzdb.AtomicGroups.Class_Groups)
		fmt.Printf("\n Classes: %+v\n", wzdb.TableMap["CLASSES"])
		for _, c := range wzdb.TableMap["CLASSES"] {
			// ??
			ag_gs := map[int][]string{}

			ndivs := len(wzdb.ActiveDivisions[c])
			cgs := wzdb.AtomicGroups.Class_Groups[c]
			cags := wzdb.AtomicGroups.Group_Atomics[wzbase.ClassGroup{
				CIX: c, GIX: 0,
			}]
			cnode := wzdb.NodeList[c].Node.(wzbase.Class)
			fmt.Printf("\n Class %s: %+v // %d\n",
				cnode.ID,
				cags,
				ndivs,
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

				if ndivs == 0 {

				} else if ndivs == 1 {

				} else {

				}

		}
	*/
	// It looks like I don't need all that stuff for fet if I don't
	// use the fet-Categories. Where there is more than one division I
	// can use the atomic group numbers (e.g. "013") as subgroups.
	// The fet-Groups are all the groups used in W365 for actual lessons.
	// If I add compound-groups (e.g. B = BG + R) with only one division,
	// the subgroups would be BG and R, and so on.

	//return

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
	xmlitem = getCourses(&wzdb)
	fmt.Printf("\n*** fet:\n%v\n", xmlitem)

	/*
		cg0 := wzbase.CourseGroups{}
		cg := wzbase.CourseGroups{}
		cg = append(cg, wzbase.ClassDivGroups{Class: 386, Div: 0, Groups: []int{308}})
		cg = append(cg, wzbase.ClassDivGroups{Class: 2, Div: -1})
		cg = append(cg, wzbase.ClassDivGroups{Class: 386, Div: 0, Groups: []int{308, 328}})
		cg = append(cg, wzbase.ClassDivGroups{Class: 2, Div: 1, Groups: []int{10, 13}})
		//cg = append(cg, wzbase.ClassDivGroups{Class: 386, Div: 0, Groups: []int{311, 328}})
		cg = append(cg, wzbase.ClassDivGroups{Class: 2, Div: 2, Groups: []int{14, 15}})
		//cg = append(cg, wzbase.ClassDivGroups{Class: 386, Div: 1, Groups: []int{18}})
		if !cg0.AddCourseGroups(wzdb.NodeList, cg) {
			log.Fatalln("INCOMPATIBLE GROUP")
		}
		log.Printf("\n --> %+v\n", cg0)
	*/
}
