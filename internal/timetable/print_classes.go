package timetable

import (
	"encoding/json"
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
	"log"
	"os"
	"slices"
	"strings"
)

type Tile struct {
	Day      int    `json:"day"`
	Hour     int    `json:"hour"`
	Duration int    `json:"duration"`
	Fraction int    `json:"fraction"`
	Offset   int    `json:"offset"`
	Total    int    `json:"total"`
	Centre   string `json:"centre"`
	TL       string `json:"tl"`
	TR       string `json:"tr"`
	BR       string `json:"br"`
	BL       string `json:"bl"`
}

type Timetable struct {
	Title  string
	School string
	Plan   string
	Pages  [][]interface{}
}

func PrintClasses(wzdb *wzbase.WZdata,
	plan_name string,
	activities []wzbase.Activity,
	filepath string,
) {
	var tiles []Tile
	pages := [][]interface{}{}
	//

	ref2id := wzdb.Ref2IdMap()
	// Get the rooms contained in room-groups
	room_groups := map[int][]string{}
	for _, ri := range wzdb.TableMap["ROOMS"] {
		rg := wzdb.GetNode(ri).(wzbase.Room).SUBROOMS
		if len(rg) != 0 {
			rglist := []string{}
			for _, r := range rg {
				rglist = append(rglist, ref2id[r])
			}
			slices.Sort(rglist)
			room_groups[ri] = rglist
		}
	}
	for _, a := range activities {
		// Gather the rooms.
		rooms := []string{}
		if len(a.Rooms) == 0 {
			// Check whether there are compulsory rooms (possible with
			// undeclared room-group).
			for _, r := range a.RoomNeeds.Compulsory {
				rooms = append(rooms, ref2id[r])
			}
			if len(rooms) > 1 {
				slices.Sort(rooms)
			}
		} else {
			for _, r := range a.Rooms {
				rg, ok := room_groups[r]
				if ok {
					rooms = append(rooms, rg...)
				} else {
					rooms = append(rooms, ref2id[r])
				}
			}
		}
		// Limit the length of the room list.
		var room string
		if len(rooms) > 6 {
			room = strings.Join(rooms[:5], ",") + "..."
		} else {
			room = strings.Join(rooms, ",")
		}
		// Gather the teachers.
		teachers := []string{}
		for _, t := range a.Teachers {
			teachers = append(teachers, ref2id[t])
		}
		var teacher string
		if len(teachers) > 6 {
			teacher = strings.Join(teachers[:5], ",") + "..."
		} else {
			teacher = strings.Join(teachers, ",")
		}

		// Is there any way of associating teachers with particular (sub)groups?
		// Probably not.

		// Gather student groups.
		classes := map[string][]string{} // mapping: class -> list of groups
		for _, cg := range a.Groups {
			c := ref2id[cg.CIX]
			g := ref2id[cg.GIX]
			gl, ok := classes[c]
			if ok {
				//TODO: Check for clashes? Or is that done previously?
				classes[c] = append(gl, g)
			} else if g == "" {
				log.Println("++ Full class")
				classes[c] = []string{}
			} else {
				classes[c] = []string{g}
			}
		}
		fmt.Printf("*** GROUPS (%s / %s): %#v\n", teacher, room, classes)
	}

	//
	tiles = []Tile{
		{Day: -1, Centre: "Tile 1"},
		{Day: 0, Centre: "Tile 2"},
	}
	pages = append(pages, []interface{}{"Class 1", tiles})
	tiles = []Tile{
		{Day: 1, Centre: "Tile 3"},
	}
	pages = append(pages, []interface{}{"Class 2", tiles})

	tt := Timetable{
		Title:  "Stundenpläne der Klassen",
		School: "Freie Michaelschule",
		Plan:   plan_name,
		Pages:  pages,
	}
	b, err := json.MarshalIndent(tt, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	os.Stdout.Write(b)
	fmt.Println()
}
