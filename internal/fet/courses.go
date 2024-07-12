package fet

import (
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
)

// func getCourses(wzdb *wzbase.WZdata, divisions []wzbase.DivGroups) string {
func getCourses(wzdb *wzbase.WZdata) string {
	for _, bi := range wzdb.TableMap["BLOCKS"] {
		bnode := wzdb.NodeList[wzdb.IndexMap[bi]].Node.(wzbase.Block)
		bbnode := wzdb.NodeList[wzdb.IndexMap[bnode.Base]].Node.(wzbase.Course)
		fmt.Printf("* Block %s Base %+v\n", bnode.Tag, bbnode)
		// The GROUPS field is a list of the classes/groups declared in the
		// (W365) source for the "main" entry for the block, the one specifying
		// the lessons. They are ClassGroup items.

		//TODO: There should probably be a check that there is no conflict with
		// the groups declared in the components. But this is not the job of
		// the fet generator, it should be handled when reading the source.
		// Perhaps the best result would be to ensure that the entries in
		// GROUPS covers all needed ClassGroup items for the block. Then the
		// subcourses would not need consulting for the timetable.

		for _, bci := range bnode.Components {
			bcnode := wzdb.NodeList[wzdb.IndexMap[bci]].Node.(wzbase.Course)
			fmt.Printf("* Block %s Group(s) %s\n",
				bnode.Tag, bcnode.GROUPS.Print(wzdb.NodeList),
			)
		}
	}

	/*
		trefs := wzdb.TableMap["COURSES"]
		//items := []fetCourse{}
		for _, ti := range trefs {
			cs := wzdb.NodeList[wzdb.IndexMap[ti]].Node.(wzbase.Course)
			// Determine what "type" of course it is. The groups participating
			// in any course with lessons must be in divisions. Multiple groups
			// within a class must be in the same division (?? really?). The
			// lessons can be either normal (LESSONS) or in a block/Epoche
			// (BLOCK_UNITS).

			//TODO: This is rather confused!!!
			with_activities := len(cs.LESSONS) > 0 || cs.BLOCK_UNITS > 0.0

			// The blocks can be got from table BLOCKS. The components must be
			// scanned to determine teachers, groups and rooms. This scan could be
			// done before handling the actual lessons, accumulating the block
			// data in a separate structure which can be looked up when the main
			// node (with the lessons) is found.

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
	*/
	return ""
}
