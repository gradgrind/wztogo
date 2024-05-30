package w365

import (
	"log"
	"slices"
	"strconv"

	"gradgrind/wztogo/internal/wzbase"
)

func (w365data *W365Data) read_teachers() {
	type xteacher struct {
		sortnum float64
		wid     string
		teacher wzbase.Teacher
	}

	xnodes := []xteacher{}
	for _, node := range w365data.yeartables[w365_Teacher] {
		id := node[w365_Shortcut]
		name := node[w365_Name]
		constraints := map[string]string{
			"MaxDays":               node[w365_MaxDays],
			"MaxLessonsPerDay":      node[w365_MaxLessonsPerDay],
			"MaxGapsPerDay":         node[w365_MaxGapsPerDay],
			"MinLessonsPerDay":      node[w365_MinLessonsPerDay],
			"NumberOfAfterNoonDays": node[w365_NumberOfAfterNoonDays],
		}
		a := w365data.absences(node)
		// Get additional info from the "categories"
		c := w365data.categories(node)
		// log.Printf("***(%s) %#v", name, c)
		if c.Role["NoTeacher"] {
			constraints["NoTeacher"] = "true" // mark a "dummy" teacher
		}
		if c.MaxLunchDays >= 0 {
			constraints["MaxLunchDays"] = strconv.Itoa(c.MaxLunchDays)
		}
		tnode := wzbase.Teacher{
			ID:         id,
			LASTNAME:   name,
			FIRSTNAMES: node[w365_Firstname],
			// SEX: int(node[w365_Gender]),  // 0: male, 1: female
			// There is also other personal information ...
			CONSTRAINTS:   constraints,
			NOT_AVAILABLE: a,
		}
		sf, err := strconv.ParseFloat(node[w365_ListPosition], 64)
		if err != nil {
			log.Fatal(err)
		}
		xnodes = append(xnodes, xteacher{sf, node[w365_Id], tnode})
	}
	// Sort the teachers according to the Waldorf 365 ListPosition
	slices.SortFunc(xnodes, func(a, b xteacher) int {
		if a.sortnum <= b.sortnum {
			return -1
		}
		return 1
	})
	for _, xt := range xnodes {
		w365data.add_node("TEACHERS", xt.teacher, xt.wid)
	}
}
