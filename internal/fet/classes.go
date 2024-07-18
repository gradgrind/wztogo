package fet

import (
	"encoding/xml"
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
)

//type fetCategory struct {
//	//XMLName             xml.Name `xml:"Category"`
//	Number_of_Divisions int
//	Division            []string
//}

type fetSubgroup struct {
	Name string // 13.m.MaE
	//Number_of_Students int // 0
	//Comments string // ""
}

type fetGroup struct {
	Name string // 13.K
	//Number_of_Students int // 0
	//Comments string // ""
	Subgroup []fetSubgroup
}

type fetClass struct {
	//XMLName  xml.Name `xml:"Year"`
	Name     string
	Comments string
	//Number_of_Students int (=0)
	// The information regarding categories, divisions of each category,
	// and separator is only used in the dialog to divide the year
	// automatically by categories.
	//	Number_of_Categories int    // 0 or 1
	//	Separator            string // "."
	//	Category             []fetCategory
	Group []fetGroup
}

type fetStudentsList struct {
	XMLName xml.Name `xml:"Students_List"`
	Year    []fetClass
}

// Note that any class divisions with no actual lessons should not appear
// in the atomic groups. This is handled before calling this function so
// that wzdb.AtomicGroups covers only these "active" divisions.
func getClasses(wzdb *wzbase.WZdata) fetStudentsList {
	//	trefs := wzdb.TableMap["CLASSES"]
	items := []fetClass{}
	for _, c := range wzdb.TableMap["CLASSES"] {
		//    for _, ti := range trefs {
		//		cl := wzdb.NodeList[wzdb.IndexMap[ti]].Node.(wzbase.Class)
		cl := wzdb.GetNode(c).(wzbase.Class)
		cgs := wzdb.AtomicGroups.Class_Groups[c]
		agmap := wzdb.AtomicGroups.Group_Atomics
		cags := agmap[wzbase.ClassGroup{
			CIX: c, GIX: 0,
		}]

		//		divs := cl.DIVISIONS
		//nc := 0
		//		if len(divs) > 0 {
		//if cags.GetCardinality() > 1 {
		//	nc = 1
		//}
		cname := cl.SORTING //?
		calt := cl.ID       //?
		groups := []fetGroup{}
		if cags.GetCardinality() > 1 {
			for _, cg := range cgs {
				g := wzdb.GetNode(cg.GIX).(wzbase.Group).ID
				gags := agmap[cg]

				subgroups := []fetSubgroup{}
				for _, ag := range gags.ToArray() {
					subgroups = append(subgroups,
						fetSubgroup{Name: fmt.Sprintf("%s.%03d", cname, ag)},
					)
					//ag_gs[int(ag)] = append(ag_gs[int(ag)], g)
				}
				groups = append(groups, fetGroup{
					Name:     fmt.Sprintf("%s.%s", cname, g),
					Subgroup: subgroups,
				})
			}
		}
		items = append(items, fetClass{
			Name:     cname,
			Comments: calt,
			Group:    groups,
		})
		/*
			items = append(items, fetClass{
				Name:                 cl.SORTING, //?
				Comments:             cl.ID,      //?
				Number_of_Categories: nc,
				Separator:            ".",
			})
		*/
		//fmt.Printf("\nCLASS %s: %+v\n", cl.SORTING, cl.DIVISIONS)
	}
	return fetStudentsList{
		Year: items,
	}
}
