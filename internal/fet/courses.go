package fet

import (
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
	"log"
)

// TODO: Handle the groups ...
// Note that there may well be "superfluous" divisions – ones with no
// actual lessons associated. These should be stripped out for fet!
// That might be an argument for not generating the atomic groups until
// it is clear in which form they are needed.
func getCourses(wzdb *wzbase.WZdata, divisions []wzbase.DivGroups) string {
	trefs := wzdb.TableMap["COURSES"]
	//items := []fetCourse{}
	for _, ti := range trefs {
		cs := wzdb.NodeList[wzdb.IndexMap[ti]].Node.(wzbase.Course)
		// Determine what "type" of course it is. The groups participating
		// in any course with lessons must be in divisions. Multiple groups
		// within a class must be in the same division. The lessons can be
		// either normal (LESSONS) or in a block/Epoche (BLOCK_UNITS).

		//TODO: This is rather confused!!!
		with_activities := len(cs.LESSONS) > 0 || cs.BLOCK_UNITS > 0.0
		g_used := map[int]bool{}
		g_div := map[int]int{}
		realdiv := []bool{}
		for i, div := range divisions {
			nodiv := div.Tag == ""
			divr := false
			for _, g := range div.Groups {
				if with_activities {
					if nodiv {
						log.Fatalf(
							"Course has group (%d) not in division: %+v",
							g, cs,
						)
					}
					g_used[g] = true
					divr = true
				}
				g_div[g] = i
			}
			realdiv = append(realdiv, divr)
		}

		if len(cs.GROUPS) > 0 {

			fmt.Printf("§Groups: %+v\n", cs)
		} else {
			// Presumable a "special" activity not involving students – or
			// erroneous or incomplete data ...
			// TODO: check teachers? (If there aren't any, what is the point
			// of this item?)
			fmt.Printf("§No groups: %+v\n", cs)
		}

	}
	return ""
}
