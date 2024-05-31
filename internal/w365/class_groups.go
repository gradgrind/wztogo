package w365

import (
	"fmt"
	"slices"
	"strings"

	"gradgrind/wztogo/internal/wzbase"
)

// Manage the reading of classes and the associated students, groups
// and divisions.
func (w365data *W365Data) read_groups() {

	w365data.read_students()

}

func (w365data *W365Data) read_students() {

	type xstudent struct {
		wid     string
		student wzbase.Student
	}

	xnodes := []xstudent{}
	for _, node := range w365data.yeartables[w365_Student] {
		last := node[w365_Name]
		first := node[w365_First_Name]
		all_first := node[w365_Firstnames]
		if first == "" {
			first = all_first
		}
		snode := wzbase.Student{
			ID:         node[w365_StudentId],
			SORTNAME:   make_sortname(last, first),
			LASTNAME:   last,
			FIRSTNAMES: all_first,
			FIRSTNAME:  first,
			GENDER:     node[w365_Gender],
			DATE_BIRTH: convert_date(node[w365_DateOfBirth]),
			BIRTHPLACE: node[w365_PlaceOfBirth],
			DATE_ENTRY: node[w365_DateOfEntry],
			DATE_EXIT:  node[w365_DateOfExit],
			HOME:       node[w365_Home],
			POSTCODE:   node[w365_Postcode],
			STREET:     node[w365_Street],
			EMAIL:      node[w365_Email],
			PHONE:      node[w365_PhoneNumber],
		}
		xnodes = append(xnodes, xstudent{node[w365_Id], snode})
	}
	// Sort the students alphabetically
	slices.SortFunc(xnodes, func(a, b xstudent) int {
		if a.student.SORTNAME <= b.student.SORTNAME {
			return -1
		}
		return 1
	})
	for _, xs := range xnodes {
		w365data.add_node("STUDENTS", xs.student, xs.wid)
	}
}

// TODO: The result should perhaps be "ASCIIfied"
func make_sortname(last string, first string) string {
	s0 := fmt.Sprintf("%s,%s", last, first)
	return strings.ReplaceAll(s0, " ", "_")
}
