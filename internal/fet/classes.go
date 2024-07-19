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

type notAvailableTime struct {
	XMLName xml.Name `xml:"Not_Available_Time"`
	Day     string
	Hour    string
}

type studentsNotAvailable struct {
	XMLName                       xml.Name `xml:"ConstraintStudentsSetNotAvailableTimes"`
	Weight_Percentage             int
	Students                      string
	Number_of_Not_Available_Times int
	Not_Available_Time            []notAvailableTime
	Active                        bool
}

// Note that any class divisions with no actual lessons should not appear
// in the atomic groups. This is handled before calling this function so
// that wzdb.AtomicGroups covers only these "active" divisions.
func getClasses(fetinfo *fetInfo) {
	//	trefs := wzdb.TableMap["CLASSES"]
	items := []fetClass{}
	natimes := []studentsNotAvailable{}
	for _, c := range fetinfo.wzdb.TableMap["CLASSES"] {
		//    for _, ti := range trefs {
		//		cl := wzdb.NodeList[wzdb.IndexMap[ti]].Node.(wzbase.Class)
		cl := fetinfo.wzdb.GetNode(c).(wzbase.Class)
		cgs := fetinfo.wzdb.AtomicGroups.Class_Groups[c]
		agmap := fetinfo.wzdb.AtomicGroups.Group_Atomics
		cags := agmap[wzbase.ClassGroup{
			CIX: c, GIX: 0,
		}]

		//		divs := cl.DIVISIONS
		//nc := 0
		//		if len(divs) > 0 {
		//if cags.GetCardinality() > 1 {
		//	nc = 1
		//}
		//calt := cl.SORTING //?
		cname := cl.ID
		groups := []fetGroup{}
		if cags.GetCardinality() > 1 {
			for _, cg := range cgs {
				g := fetinfo.ref2fet[cg.GIX]
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
			Name: cname,
			//Comments: calt,
			Group: groups,
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

		// "Not available" times
		nats := []notAvailableTime{}
		for d, dna := range cl.NOT_AVAILABLE {
			for _, h := range dna {
				nats = append(nats,
					notAvailableTime{
						Day: fetinfo.days[d], Hour: fetinfo.hours[h]})
			}
		}
		if len(nats) > 0 {
			natimes = append(natimes,
				studentsNotAvailable{
					Weight_Percentage:             100,
					Students:                      cname,
					Number_of_Not_Available_Times: len(nats),
					Not_Available_Time:            nats,
					Active:                        true,
				})
		}
		//fmt.Printf("==== %s: %+v\n", cname, nats)
	}
	fetinfo.fetdata.Students_List = fetStudentsList{
		Year: items,
	}
	fetinfo.fetdata.Time_Constraints_List.ConstraintStudentsSetNotAvailableTimes = natimes
}
