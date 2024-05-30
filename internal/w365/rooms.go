package w365

import (
	"log"
	"slices"
	"strconv"
	"strings"

	"gradgrind/wztogo/internal/wzbase"
)

func (w365data *W365Data) read_rooms() {
	type xroom struct {
		sortnum float64
		wid     string
		room    wzbase.Room
		rg      string // a "list" of w365ids for the component rooms
	}

	xnodes := []xroom{}
	for _, node := range w365data.yeartables[w365_Room] {
		sf, err := strconv.ParseFloat(node[w365_ListPosition], 64)
		if err != nil {
			log.Fatal(err)
		}
		a := w365data.absences(node)
		rnode := wzbase.Room{
			ID:            node[w365_Shortcut],
			NAME:          node[w365_Name],
			NOT_AVAILABLE: a,
		}
		xnodes = append(xnodes, xroom{
			sortnum: sf,
			wid:     node[w365_Id],
			room:    rnode,
			rg:      node[w365_RoomGroup], // component rooms ("list" of w365ids)
		})
	}
	// Sort the subjects according to the Waldorf 365 ListPosition
	slices.SortFunc(xnodes, func(a, b xroom) int {
		if a.sortnum <= b.sortnum {
			return -1
		}
		return 1
	})
	rglist := []xroom{} // collect room groups
	for i, xr := range xnodes {
		xr.room.X = i
		if xr.rg == "" {
			w365data.add_node("ROOMS", xr.room, xr.wid)
		} else {
			rglist = append(rglist, xr)
		}
	}
	for _, xr := range rglist {
		subrooms := []int{}
		for _, rgc := range strings.Split(xr.rg, LIST_SEP) {
			subrooms = append(subrooms, w365data.NodeMap[rgc])
		}
		xr.room.SUBROOMS = subrooms
		w365data.add_node("ROOMS", xr.room, xr.wid)
	}
}
