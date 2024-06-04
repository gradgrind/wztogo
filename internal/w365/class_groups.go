package w365

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"gradgrind/wztogo/internal/wzbase"
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
	for _, xc := range xclasses {
		w365data.add_node("CLASSES", xc.node, xc.wid)
	}
	//TODO
	/*
	   gen_class_groups(w365_db.nodes, xnode)
	   #print("  *** $GROUP_ATOM_MAP:", xnode["$GROUP_ATOM_MAP"])
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

/*
#TODO: The following classes would also be relevant for other data
# sources. Perhaps they should be moved to a different folder?
class AG(frozenset):
    def __repr__(self):
        return f"{{*{','.join(sorted(self))}*}}"

    def __str__(self):
        return AG_SEP.join(sorted(self))


def gen_class_groups(key2node, node):
    """Produce "atomic" groups for the given class partitions.
    This should be rerun whenever any change is made to the partitions –
    including just name changes because the group names are used here.
    <parts> is a list of tuples:
        - name: the partition name (can be empty)
        - list of basic partition group keys
        - list of "compound" groups:
            [compound group key, basic group key, basic group key, ...]
    """
    parts = node["PARTITIONS"]
    if not parts:
        node["$GROUP_ATOM_MAP"] = {"": set()}
        return
    # Check the input
    gset = set()
    divs1 = []
    divs1x = []
    for n, d, dx in parts:
        gs = []
        xg = {}
        divs1.append(gs)
        divs1x.append(xg)
        for gk in d:
            g = key2node[gk]["ID"]
#TODO: use something more helpful than the assertion
            assert g not in gset
            gset.add(g)
            # "Compound" groups are combinations of "basic" groups,
            # as a convenience for input and display of multiple groups
            # within a division (not supported in Waldorf365).
            # Consider a division ["A", "BG", "R"]. There could be
            # courses, say, for combination "A" + "BG". The "compound"
            # group might then be "G", defined as "G=A+BG". Obviously, if
            # this format is used, the symbols "=" and "+" should not be
            # used in group names.
            gs.append(g)   # A "basic" group
        # Deal with compound groups
        for gx in dx:
#TODO: use something more helpful than the assertion
            assert len(gx) > 2
            gc = key2node[gx[0]]["ID"]
            xgl = []
            for gk in gx[1:]:
                g = key2node[gk]["ID"]
#TODO: use something more helpful than the assertion
                assert g in gs
                xgl.append(g)
            xg[gc] = xgl
#TODO: use something more helpful than the assertion
        assert len(gs) > 1
    # Generate "atomic" groups
    g2ag = {}
    aglist = []
    for p in product(*divs1):
        ag = AG(p)
        aglist.append(ag)
        for g in p:
            try:
                g2ag[g].add(ag)
            except KeyError:
                g2ag[g] = {ag}
    for xg in divs1x:
        for g, gl in xg.items():
            ags = set()
            for gg in gl:
                ags.update(g2ag[gg])
            g2ag[g] = ags
    # Add the atomic groups for the whole class
    g2ag[""] = set(aglist)
    node["$GROUP_ATOM_MAP"] = g2ag

*/

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
