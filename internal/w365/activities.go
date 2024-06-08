package w365

import (
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
	"log"
	"strings"
)

func (w365data *W365Data) read_activities() {
	/*type xteacher struct {
		sortnum float64
		wid     string
		teacher wzbase.Teacher
	}

	xnodes := []xteacher{}
	*/
	// Get teachers
	for _, node := range w365data.yeartables[w365_Course] {
		wid := node[w365_Id]
		tlist := []int{}
		for _, s := range strings.Split(node[w365_Teachers], LIST_SEP) {
			if s != "" {
				tlist = append(tlist, w365data.NodeMap[s])
			}
		}
		fmt.Printf("%s: %+v\n", wid, tlist)
		// There must be exactly one subject. Courses with multiple subjects
		// should be replaced by (defined) blocks.
		slist := []int{}
		for _, s := range strings.Split(node[w365_Subjects], LIST_SEP) {
			if s != "" {
				slist = append(slist, w365data.NodeMap[s])
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
			log.Printf("\n????????????????\n  INVALID SUBJECT in Course %s: "+
				"%+v\n????????????????\n", wid, stlist)
			continue
		}
		subject := slist[0]
		fmt.Printf("    --> Subject: %d\n", subject)
		// Get groups
		glist := []int{}
		for _, s := range strings.Split(node[w365_Groups], LIST_SEP) {
			if s != "" {
				glist = append(glist, w365data.NodeMap[s])
			}
		}
		fmt.Printf("    --> Groups: %+v\n", glist)
		// Get rooms
		rlist := []int{}
		for _, s := range strings.Split(node[w365_PreferredRooms], LIST_SEP) {
			if s != "" {
				rlist = append(rlist, w365data.NodeMap[s])
			}
		}
		fmt.Printf("    --> Rooms: %+v\n", rlist)
		workload := node[w365_HandWorkload]
		if workload == "555.555" { // that is the "empty" value (!)
			workload = ""
		}
		fmt.Printf("    --> HandWorkload: %s\n", workload)
	}
}
