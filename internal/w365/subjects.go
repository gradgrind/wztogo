package w365

import (
	"log"
	"slices"
	"strconv"

	"gradgrind/wztogo/internal/wzbase"
)

func (w365data *W365Data) read_subjects() {
	type xsubject struct {
		sortnum float64
		wid     string
		subject wzbase.Subject
	}

	xnodes := []xsubject{}
	for _, node := range w365data.yeartables[w365_Subject] {
		snode := wzbase.Subject{
			ID:   node[w365_Shortcut],
			NAME: node[w365_Name],
		}
		sf, err := strconv.ParseFloat(node[w365_ListPosition], 64)
		if err != nil {
			log.Fatal(err)
		}
		xnodes = append(xnodes, xsubject{sf, node[w365_Id], snode})
	}
	// Sort the subjects according to the Waldorf 365 ListPosition
	slices.SortFunc(xnodes, func(a, b xsubject) int {
		if a.sortnum <= b.sortnum {
			return -1
		}
		return 1
	})
	for _, xs := range xnodes {
		w365data.add_node("SUBJECTS", xs.subject, xs.wid)
	}
}
