package w365

import (
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
	"log"
	"strconv"
	"strings"
)

func (w365data *W365Data) read_activities() {
	type xblock struct {
		// * The values below are w365 Ids of the courses
		base       string
		components []string
	}

	blocks := map[string]*xblock{}
	active_groups := map[int]bool{}
	for _, node := range w365data.yeartables[w365_Course] {
		wid := node[w365_Id]
		// * Get teachers
		tlist := []int{}
		for _, s := range strings.Split(node[w365_Teachers], LIST_SEP) {
			if s != "" {
				tlist = append(tlist, w365data.NodeMap[s])
			}
		}
		//fmt.Printf("%s: %+v\n", wid, tlist)
		// * There must be exactly one subject. Courses with multiple subjects
		// should be replaced by (defined) blocks at the source.
		slist := []int{}
		sstring := node[w365_Subjects]
		if sstring != "" {
			for _, s := range strings.Split(sstring, LIST_SEP) {
				if s != "" {
					slist = append(slist, w365data.NodeMap[s])
				}
			}
		}
		// * Get subject name(s)
		stlist := ""
		if len(slist) != 0 {
			snlist := []string{}
			for _, i := range slist {
				snlist = append(
					snlist, w365data.NodeList[i].Node.(wzbase.Subject).ID,
				)
			}
			stlist = strings.Join(snlist, ",")
		}
		// * Get groups for this course
		type cdxItem struct {
			div    int
			groups map[int]bool
		}
		cdmap := map[int]cdxItem{} // This collects class + division,
		// only one division per class is permitted.
		for _, s := range strings.Split(node[w365_Groups], LIST_SEP) {
			if s != "" {
				// The reference can be to either a W365-Group or to a
				// W365-Year (class). Internally these are transformed to
				// ClassGroup items, so that only a single type is used.
				gi := w365data.NodeMap[s]
				cg := w365data.group_classgroup[gi]
				cdx, ok := cdmap[cg.CIX]
				if ok {
					// If the whole class is already covered, anything else
					// is superfluous.
					if cdx.div == -1 {
						log.Printf("Class %s, subject %s: superfluous group",
							w365data.NodeList[cg.CIX].Node.(wzbase.Class).ID,
							stlist,
						)
					} else {
						// If the new group is "whole class", set the division
						// accordingly, previous entries are superfluous.
						if cg.GIX == 0 {
							cdmap[cg.CIX] = cdxItem{div: -1}
							log.Printf(
								"Class %s, subject %s: superfluous group",
								w365data.NodeList[cg.CIX].Node.(wzbase.Class).ID,
								stlist,
							)
						} else {
							// If the new division matches the existing one, add
							// the group.
							gd := w365data.class_group_div[cg.CIX][cg.GIX]
							if gd == cdx.div {
								cdx.groups[cg.GIX] = true
								// Check that the result is not equivalent
								// to "whole class"
								divs := w365data.NodeList[cg.CIX].Node.(wzbase.Class).DIVISIONS[gd]
								if len(cdx.groups) == len(divs.Groups) {
									// Substitute the whole class
									cdmap[cg.CIX] = cdxItem{div: -1}
								}
							} else {
								log.Printf("Class %s, subject %s:"+
									" groups from different divisions",
									w365data.NodeList[cg.CIX].Node.(wzbase.Class).ID,
									stlist,
								)
							}
						}
					}
				} else {
					// First group in this class
					if cg.GIX == 0 {
						// whole class
						cdmap[cg.CIX] = cdxItem{div: -1}
					} else {
						gd := w365data.class_group_div[cg.CIX][cg.GIX]
						cdmap[cg.CIX] = cdxItem{gd, map[int]bool{cg.GIX: true}}
					}
				}
			}
		}
		// Use cdmap to make a wzbase.CourseGroups list. Note that the groups
		// will not be properly sorted.
		cdglist := make(wzbase.CourseGroups, len(cdmap))
		i := 0
		for c, dig := range cdmap {
			keys := make([]int, len(dig.groups))
			if dig.div >= 0 {
				j := 0
				for k := range dig.groups {
					keys[j] = k
					j++
				}
			}
			cdglist[i] = wzbase.ClassDivGroups{
				Class: c, Div: dig.div, Groups: keys,
			}
			i++
		}

		//TODO: The groups (cdgmap) should later be added to the groups for
		// the block, which should then be checked against those of the base
		// course.

		// * Report invalid subject. It is placed here so that the group(s)
		// are available for the report.
		if len(slist) != 1 {
			log.Printf("\n=========================================\n"+
				"  !!!  INVALID SUBJECT (%s) in Class/Group(s) %s\n"+
				"=========================================\n",
				stlist, cdglist.Print(w365data.NodeList),
			)
			continue
		}
		subject := slist[0]
		// * Get rooms
		rlist := []int{}
		for _, s := range strings.Split(node[w365_PreferredRooms], LIST_SEP) {
			if s != "" {
				rlist = append(rlist, w365data.NodeMap[s])
			}
		}
		//fmt.Printf("    --> Rooms: %+v\n", rlist)
		// * Get the "workload override" (<0 => no override)
		var workload float64
		wl := node[w365_HandWorkload]
		if wl == "555.555" { // that is the "empty" value (!)
			workload = -1.0
		} else {
			f, err := strconv.ParseFloat(wl, 32)
			if err != nil {
				log.Fatal(err)
			}
			workload = f
		}
		//fmt.Printf("    --> HandWorkload: %f\n", workload)
		// * Divide lessons up according to duration – remove decimal places
		h, _, _ := strings.Cut(node[w365_HoursPerWeek], ".")
		total_duration, err := strconv.Atoi(h)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("    --> total_duration: %d\n", total_duration)
		//NOTE: Not all multiple lesson possibilities are supported here,
		// just take the first digit.
		dlm, _, _ := strings.Cut(node[w365_DoubleLessonMode], ",")
		ll, err := strconv.Atoi(dlm)
		if err != nil {
			log.Fatal(err)
		}
		lessons := []int{}
		nl := total_duration
		for nl > 0 {
			if nl < ll {
				// reduced length for last entry
				lessons = append(lessons, nl)
				break
			}
			lessons = append(lessons, ll)
			nl -= ll
		}
		//fmt.Printf("    --> Lessons: %+v\n", lessons)
		// * Get additional info from the "categories"
		cat := w365data.categories(node)
		block_units := 0.0
		epweeks := node[w365_EpochWeeks]
		if epweeks != "0.0" {
			if cat.Block == "" {
				sbj := w365data.NodeList[slist[0]].Node.(wzbase.Subject).ID
				log.Printf("\n=========================================\n"+
					"  !!!  EpochWeeks without block tag (%s) in Class/Group(s) %s\n"+
					"=========================================\n",
					sbj, cdglist.Print(w365data.NodeList),
				)
			} else {
				// Component of a named block.
				if len(lessons) > 0 {
					sbj := w365data.NodeList[slist[0]].Node.(wzbase.Subject).ID
					log.Fatalf("Class/Group %s: A course, subject %s,"+
						" in block %s has both lessons and weeks",
						cdglist.Print(w365data.NodeList), sbj, cat.Block)
				}
			}
			block_units, err = strconv.ParseFloat(epweeks, 32)
			if err != nil {
				log.Fatal(err)
			}
		}
		// * Is the node part of a block?
		if cat.Block != "" {
			bdata, ok := blocks[cat.Block]
			if ok {
				if block_units == 0.0 {
					// This is the "base" course of the block, there may
					// only be one of these. It specifies the actual lessons.
					if bdata.base != "" {
						log.Fatalf(
							"Block '%s' has more than one 'base'", cat.Block,
						)
					}
					if !cat.Role["NoReport"] {
						//TODO: Really fatal? Anyway, I would need a better
						// message, to identify the course.
						log.Fatalf(
							"'Epochenschiene' without 'NoReport': "+
								"course Id = %s\n", wid,
						)
					}
					bdata.base = wid
				} else {
					bdata.components = append(bdata.components, wid)
				}
			} else {
				if block_units == 0.0 {
					// This is the "base" course of the block.
					blocks[cat.Block] = &xblock{wid, []string{}}
				} else {
					blocks[cat.Block] = &xblock{"", []string{wid}}
				}
			}
			//fmt.Printf("    --> BLOCK DATA: %+v\n", blocks[cat.Block])
		}
		// ** Add the course to the database
		flags := map[string]bool{
			"NotColliding":  cat.NotColliding,
			"NoReport":      cat.Role["NoReport"],
			"NotInRegister": cat.Role["NotRegistered"],
			"WholeDayBlock": cat.Role["WholeDayBlock"],
		}
		cnode := wzbase.Course{
			TEACHERS:        tlist,
			GROUPS:          cdglist,
			SUBJECT:         subject,
			ROOM_WISH:       rlist,
			WORKLOAD:        workload,
			WORKLOAD_FACTOR: cat.WorkloadFactor,
			LESSONS:         lessons,
			BLOCK_UNITS:     block_units,
			FLAGS:           flags,
		}
		// * Mark the groups used in lessons – and blocks – as "active".
		// Only active groups are significant for the timetable.
		w365data.add_node("COURSES", cnode, wid)
		if len(lessons) > 0 || block_units > 0.0 {
			for _, cdgs := range cdglist {
				for _, g := range cdgs.Groups {
					active_groups[g] = true
				}
			}
		}
	}
	w365data.ActiveGroups = active_groups
	// * Add the blocks to the database, checking that the component groups
	// are compatible with the base groups, in a rather flexible way ...
	for b, xb := range blocks {
		xbi := w365data.NodeMap[xb.base]
		basenode := w365data.NodeList[xbi]
		basegroups := basenode.Node.(wzbase.Course).GROUPS
		//TODO--
		fmt.Printf("\n $$$ basegroups %s: %#v\n", b, basegroups)
		xcl := []int{}
		for _, xc := range xb.components {
			xci := w365data.NodeMap[xc]
			xcl = append(xcl, xci)
		}
		bnode := wzbase.Block{
			Tag:        b,
			Base:       xbi,
			Components: xcl,
		}
		w365data.add_node("BLOCKS", bnode, "")
	}
}
