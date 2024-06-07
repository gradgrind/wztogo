package w365

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"gradgrind/wztogo/internal/wzbase"
	// mapset "github.com/deckarep/golang-set/v2"
)

// Manage the reading of classes and the associated students, groups
// and divisions.
func (w365data *W365Data) read_groups() {
	w365data.read_students()
	w365data.read_subgroups()
	// Get all class divisions
	type xclass struct {
		node wzbase.Class
		wid  string
	}
	xclasses := []xclass{}
	wid2divgroups := map[string]wzbase.DivGroups{}
	for _, node := range w365data.yeartables[w365_YearDiv] {
		name := node[w365_Name]
		gklist := []int{}
		for _, n := range strings.Split(node[w365_Groups], LIST_SEP) {
			gklist = append(gklist, w365data.NodeMap[n])
		}
		wid2divgroups[node[w365_Id]] = wzbase.DivGroups{
			Tag: name, Groups: gklist,
		}
		//fmt.Printf("??? DivGroup %s: %+v\n", name, gklist)
	}
	// Get data associated with the classes
	for _, node := range w365data.yeartables[w365_Year] { // Waldorf365: "Grade"
		clevel := node[w365_Level]
		cletter := node[w365_Letter]
		cltag := clevel + cletter
		// Get the students associated with the class
		students := node[w365_Students]
		skeys := []int{}
		if students != "" {
			for _, s := range strings.Split(students, LIST_SEP) {
				skeys = append(skeys, w365data.NodeMap[s])
			}
		}
		// Get all groups associated with the class
		class_groups := map[int]int{}
		for _, n := range strings.Split(node[w365_Groups], LIST_SEP) {
			class_groups[w365data.NodeMap[n]]++
		}
		// Get the divisions associated with the class, and their groups
		divlist := []wzbase.DivGroups{}
		divs := node[w365_YearDivs]
		if divs != "" {
			for i, divid := range strings.Split(divs, LIST_SEP) {
				divgroups := wid2divgroups[divid]
				if divgroups.Tag == "" {
					divgroups.Tag = fmt.Sprintf("#%d", i)
				}
				for _, g := range divgroups.Groups {
					class_groups[g]--
				}
				divlist = append(divlist, divgroups)
			}
		}
		g0list := []int{}
		for g, n := range class_groups {
			if n > 0 {
				g0list = append(g0list, g)
			}
		}
		if len(g0list) > 1 {
			divlist = append(divlist, wzbase.DivGroups{
				Tag: "", Groups: g0list,
			})
		}
		af, err := strconv.ParseFloat(node[w365_EpochFactor], 64)
		if err != nil {
			log.Fatal(err)
		}
		xnode := wzbase.Class{
			ID:           cltag,
			SORTING:      fmt.Sprintf("%02s%s", clevel, cletter),
			BLOCK_FACTOR: af,
			STUDENTS:     skeys,
			DIVISIONS:    divlist,
		}
		xclasses = append(xclasses, xclass{xnode, node[w365_Id]})
	}
	// Sort the classes
	slices.SortFunc(xclasses, func(a, b xclass) int {
		if a.node.SORTING <= b.node.SORTING {
			return -1
		}
		return 1
	})
	// Add each class to the data structures and build atomic groups
	ag := wzbase.NewAtomicGroups()
	for _, xc := range xclasses {
		i := w365data.add_node("CLASSES", xc.node, xc.wid)
		ag.Add_class_groups(i, xc.node)
	}
	//TODO: Add the atomic group data to the general data structures.
	// Which form should I use? Probably cg2rbm.
	// Then I can comment out c2groups.
	// On the other hand I should separate out the class-handling bit
	// so that a class can be changed individually (when editing).
	// And it can be moved out of the w365 specific bit?
	// Because of the serial indexing it might be better to redo the
	// atomic groups completely when there is a change to class groups?
	// How are group changes handled as far as courses, etc., are
	// concerned? Surely if a group is removed, any associated courses
	// must also be removed (prompt!), mustn't they?

	//TODO: Compound groups ...

	//TODO:
	/*
	   constraints = {
	       _f: node[f]
	       for f, _f in (
	           (_ForceFirstHour, "ForceFirstHour"),
	           (_MaxLessonsPerDay, "MaxLessonsPerDay"),
	           (_MinLessonsPerDay, "MinLessonsPerDay"),
	           (_NumberOfAfterNoonDays, "NumberOfAfterNoonDays"),
	       )
	   }
	   xnode["CONSTRAINTS"] = constraints
	   a = absences(w365_db.idmap, node)
	   if a:
	       xnode["NOT_AVAILABLE"] = a
	   c = categories(w365_db.idmap, node)
	   if c:
	       xnode["EXTRA"] = c
	*/
}

func (w365data *W365Data) read_subgroups() {
	// I don't think sorting makes much sense here.
	for _, node := range w365data.yeartables[w365_Group] {
		students := node[w365_Students]
		skeys := []int{}
		if students != "" {
			for _, s := range strings.Split(students, LIST_SEP) {
				skeys = append(skeys, w365data.NodeMap[s])
			}
		}
		group := wzbase.Group{
			// Only the "Shortcut" is used for naming.
			ID:       node[w365_Shortcut],
			STUDENTS: skeys,
		}
		w365data.add_node("GROUPS", group, node[w365_Id])

	}
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
