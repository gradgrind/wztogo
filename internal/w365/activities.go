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
		// should be replaced by (defined) blocks.
		slist := []int{}
		for _, s := range strings.Split(node[w365_Subjects], LIST_SEP) {
			if s != "" {
				slist = append(slist, w365data.NodeMap[s])
			}
		}
		subject := slist[0]
		//fmt.Printf("    --> Subject: %d\n", subject)
		// * Get groups
		glist := []int{}
		for _, s := range strings.Split(node[w365_Groups], LIST_SEP) {
			if s != "" {
				glist = append(glist, w365data.NodeMap[s])
			}
		}
		if len(slist) != 1 {
			snlist := []string{}
			for _, i := range slist {
				snlist = append(
					snlist, w365data.NodeList[i].Node.(wzbase.Subject).ID,
				)
			}
			stlist := strings.Join(snlist, ",")
			gnlist := []string{}
			for _, i := range glist {
				gnlist = append(
					gnlist, w365data.NodeList[i].Node.(wzbase.Class).ID,
				)
			}
			gtlist := strings.Join(gnlist, ",")
			log.Printf("\n=========================================\n"+
				"  !!!  INVALID SUBJECT (%s) in Class(es) %s\n"+
				"=========================================\n",
				stlist, gtlist,
			)
			continue
		}
		//fmt.Printf("    --> Groups: %+v\n", glist)
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
				gnlist := []string{}
				for _, i := range glist {
					gnlist = append(
						gnlist, w365data.NodeList[i].Node.(wzbase.Class).ID,
					)
				}
				gtlist := strings.Join(gnlist, ",")
				log.Printf("\n=========================================\n"+
					"  !!!  EpochWeeks without block tag (%s) in Class(es) %s\n"+
					"=========================================\n",
					sbj, gtlist,
				)
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
			GROUPS:          glist,
			SUBJECT:         subject,
			ROOM_WISH:       rlist,
			WORKLOAD:        workload,
			WORKLOAD_FACTOR: cat.WorkloadFactor,
			LESSONS:         lessons,
			BLOCK_UNITS:     block_units,
			FLAGS:           flags,
		}
		w365data.add_node("COURSES", cnode, wid)
		if len(lessons) > 0 || block_units > 0.0 {
			for _, g := range glist {
				active_groups[g] = true
			}
		}
	}
	w365data.ActiveGroups = active_groups
	// * Add the blocks to the database
	for b, xb := range blocks {
		xbi := w365data.NodeMap[xb.base]
		xcl := []int{}
		for _, xc := range xb.components {
			xcl = append(xcl, w365data.NodeMap[xc])
		}
		bnode := wzbase.Block{
			Tag:        b,
			Base:       xbi,
			Components: xcl,
		}
		w365data.add_node("BLOCKS", bnode, "")
	}
}
