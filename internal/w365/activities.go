package w365

import (
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
		// * Get groups for this course. If a whole class is not selected,
		// only one division in the class is permitted.
		cgroups := wzbase.CourseGroups{}
		for _, s := range strings.Split(node[w365_Groups], LIST_SEP) {
			if s != "" {
				// The reference can be to either a W365-Group or to a
				// W365-Year (class). Internally these are transformed to
				// ClassGroup items, so that only a single type is used.
				gi := w365data.NodeMap[s]
				cg := w365data.group_classgroup[gi]
				if cg.GIX == 0 {
					cgroups = append(cgroups,
						wzbase.ClassDivGroups{Class: cg.CIX, Div: -1})
				} else {
					// Get the division
					d := w365data.class_group_div[cg.CIX][cg.GIX]
					cgroups = append(cgroups,
						wzbase.ClassDivGroups{
							Class:  cg.CIX,
							Div:    d,
							Groups: []int{cg.GIX},
						})
				}
			}
		}
		cdglist := wzbase.CourseGroups{}
		cdglist.AddCourseGroups(w365data.NodeList, cgroups)
		// The groups (cdglist) are later added to the groups for
		// the block, which are then be checked against those of the base
		// course.

		// * Report invalid subject. It is placed here so that the group(s)
		// are available for the report.
		if len(slist) != 1 {
			log.Printf("\n=========================================\n"+
				"  !!!  INVALID SUBJECT (%s) in Class/Group(s) %s\n"+
				"=========================================\n",
				stlist, cdglist.Print(w365data),
			)
			continue
		}
		subject := slist[0]
		// * Get rooms
		// I'd like to keep the room handling fairly flexible. That would
		// mean accepting compulsory rooms, rooms chosen from a list,
		// rooms requiring manual input. The wzbase.RoomSpec structure
		// should be able to cater for this. Waldorf 365 only supports a
		// subset.

		// Room groups should not be in the database. These get converted
		// to multiple compulsory rooms.

		var roomspec wzbase.RoomSpec
		rfield := node[w365_PreferredRooms]
		if rfield == "" {
			// Assume a user-input room
			roomspec = wzbase.RoomSpec{
				Compulsory: []int{},
				Choices:    [][]int{},
				UserInput:  1,
			}
		} else {
			rlist := [][]int{}
			for _, s := range strings.Split(rfield, LIST_SEP) {
				rg, ok := w365data.room_group[s]
				if ok {
					rlist = append(rlist, rg)
				} else {
					rlist = append(rlist, []int{w365data.NodeMap[s]})
				}
			}
			if len(rlist) == 1 {
				// A single compulsory room (list)
				roomspec = wzbase.RoomSpec{
					Compulsory: rlist[0],
					Choices:    [][]int{},
					UserInput:  0,
				}
			} else {
				// A single room-choice list
				choices := []int{}
				for _, rc := range rlist {
					if len(rc) != 1 {
						// Don't allow room groups in a choice list
						log.Printf("\n=========================================\n"+
							"  !!!  Room group in room choice list: %s (%s)\n"+
							"=========================================\n",
							cdglist.Print(w365data), stlist,
						)
						choices = choices[:0] // clear list
						break
					}
					choices = append(choices, rc[0])
				}
				roomspec = wzbase.RoomSpec{
					Compulsory: []int{},
					Choices:    [][]int{choices},
					UserInput:  0,
				}
			}
		}
		//fmt.Printf("    --> Rooms: %+v\n", roomspec)
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
					sbj, cdglist.Print(w365data),
				)
			} else {
				// Component of a named block.
				if len(lessons) > 0 {
					sbj := w365data.NodeList[slist[0]].Node.(wzbase.Subject).ID
					log.Fatalf("Class/Group %s: A course, subject %s,"+
						" in block %s has both lessons and weeks",
						cdglist.Print(w365data), sbj, cat.Block)
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
						//TODO: Really fatal?
						log.Fatalf("\n=========================================\n"+
							"  !!!  'Epochenschiene' without 'NoReport': %s (%s)\n"+
							"=========================================\n",
							cdglist.Print(w365data), stlist,
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
			ROOM_WISH:       roomspec,
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
	// Also the teachers and rooms should be gathered.
	for b, xb := range blocks {
		xbi := w365data.NodeMap[xb.base]
		basecourse := w365data.NodeList[xbi].Node.(wzbase.Course)
		blockgroups := wzbase.CourseGroups{}
		basegroups := basecourse.GROUPS
		bgmap := map[int]bool{}
		if len(basegroups) != 0 {
			for _, cdg := range basegroups {
				if cdg.Div == -1 {
					bgmap[-cdg.Class] = true
				} else {
					for _, g := range cdg.Groups {
						bgmap[g] = true
					}
				}
			}
		}
		//fmt.Printf("\n $$$ basegroups %s: %#v\n", b, basegroups)
		xcl := []int{}
		var tlist []int // collect teacher indexes
		tlist = append(tlist, basecourse.TEACHERS...)
		baserooms := basecourse.ROOM_WISH
		rspec := wzbase.RoomSpec{ // collect room indexes
			Compulsory: baserooms.Compulsory,
			Choices:    [][]int{},
			UserInput:  0,
		}
		for _, xc := range xb.components {
			xci := w365data.NodeMap[xc]
			xcl = append(xcl, xci)
			node := w365data.NodeList[xci]
			course := node.Node.(wzbase.Course)
			// + Deal with the groups
			groups := course.GROUPS
			if len(basegroups) == 0 {
				// Add this course's groups to blockgroups.
				if !blockgroups.AddCourseGroups(w365data.NodeList, groups) {
					log.Fatalf("Incompatible group in course %s\n",
						course.Print(w365data))
				}
			} else {
				// Check that this course's groups are a subset of basegroups
				for _, cdg := range groups {
					c := cdg.Class
					if !bgmap[-c] {
						// Full class not included, check groups
						if cdg.Div == -1 {
							log.Fatalf("Course class not in block groups: %s\n",
								course.Print(w365data))
						} else {
							for _, g := range cdg.Groups {
								if !bgmap[g] {
									log.Fatalf("Course group not subset of block groups: %s\n",
										course.Print(w365data))
								}
							}
						}
					}
				}
			}
			// + Deal with the teachers.
		tloop:
			for _, ti := range course.TEACHERS {
				for _, ti0 := range tlist {
					if ti0 == ti {
						continue tloop
					}
				}
				tlist = append(tlist, ti)
			}
			// + Deal with the rooms.
		rloop1:
			for _, ri := range course.ROOM_WISH.Compulsory {
				for _, ri0 := range rspec.Compulsory {
					if ri0 == ri {
						continue rloop1
					}
				}
				rspec.Compulsory = append(rspec.Compulsory, ri)
			}
			// Handling choices is somewhat tricky, to put it mildly ...
			// Indeed it is not at all obvious how they should be handled
			// in a useful way. So I'll cop out a bit and only pass through
			// choices in the base course.
			if len(course.ROOM_WISH.Choices) != 0 {
				log.Printf("Block course, room choices will be ignored: %s\n",
					course.Print(w365data))
			}
			// An empty room field is ignored here.
		}
		if len(basecourse.TEACHERS) > 0 {
			txlist := []string{}
		tloop2:
			for _, ti := range tlist {
				for _, ti0 := range basecourse.TEACHERS {
					if ti0 == ti {
						continue tloop2
					}
				}
				txlist = append(txlist,
					w365data.NodeList[ti].Node.(wzbase.Teacher).ID)
			}
			if len(txlist) > 0 {
				log.Printf("Block course, added teachers: %s\n  to %s\n",
					strings.Join(txlist, ","),
					basecourse.Print(w365data))
			}
		}
		if len(baserooms.Choices) == 1 {
			// Attempt at doing something with room choice list in the
			// base course ...
			//TODO: any better ideas?
			rclist := []int{}
		rloop2:
			for _, ri := range baserooms.Choices[0] {
				// Waldorf 365 can only provide a single list.
				for _, ri0 := range rspec.Compulsory {
					if ri0 == ri {
						continue rloop2
					}
				}
				rclist = append(rclist, ri)
			}
			if len(rclist) > 0 {
				if len(rclist) == 1 {
					rspec.Compulsory = append(rspec.Compulsory, rclist[0])
				} else {
					rspec.Choices = [][]int{rclist}
				}
			}
		}
		if len(basegroups) == 0 {
			basegroups = blockgroups
		}
		bnode := wzbase.Block{
			Tag:           b,
			Base:          xbi,
			Components:    xcl,
			BlockGroups:   basegroups,
			BlockTeachers: tlist,
			BlockRooms:    rspec,
		}
		w365data.add_node("BLOCKS", bnode, "")
	}
}
