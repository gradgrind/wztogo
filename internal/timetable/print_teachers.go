package timetable

import (
	"encoding/json"
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

func PrintTeachers(wzdb *wzbase.WZdata,
	plan_name string,
	activities []wzbase.Activity,
	datadir string,
	outpath string, // full path to output pdf
) {
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

	// Class-group infrastructure
	divmap := map[string][][]string{}
	for c, ad := range wzdb.ActiveDivisions {
		divlist := [][]string{}
		for _, div := range ad {
			gs := []string{}
			for _, g := range div {
				gs = append(gs, ref2id[g])
			}
			divlist = append(divlist, gs)
		}
		divmap[ref2id[c]] = divlist
		fmt.Printf(" $$$ AD %s: %+v\n", ref2id[c], divlist)
	}

	// Generate the tiles.
	teacherTiles := map[string][]Tile{}
	for _, a := range activities {
		if a.Day < 0 {
			// Unplaced activity, skip it.
			continue
		}
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

		//TODO: Is there any way of associating teachers with particular
		// (sub)groups? Probably not.

		// Gather student groups.
		var students string
		cgroups := []string{}
		for _, cg := range a.Groups {
			cgroups = append(cgroups, cg.Print(wzdb))
		}
		if len(cgroups) > 6 {
			students = strings.Join(cgroups[:5], ",") + "..."
		} else {
			students = strings.Join(cgroups, ",")
		}

		// Go through the teachers.
		for _, ti := range a.Teachers {
			teacher := ref2id[ti]
			tile := Tile{
				Day:      a.Day,
				Hour:     a.Hour,
				Duration: a.Duration,
				Fraction: 1,
				Offset:   0,
				Total:    1,
				Centre:   students,
				TL:       ref2id[a.Subject],
				BR:       room,
			}
			teacherTiles[teacher] = append(teacherTiles[teacher], tile)
		}

	}
	// Assume the teacher table is sorted!
	for _, ti := range wzdb.TableMap["TEACHERS"] {
		tnode := wzdb.GetNode(ti).(wzbase.Teacher)
		t := tnode.ID
		ctiles, ok := teacherTiles[t]
		if !ok {
			continue
		}
		pages = append(pages, []interface{}{
			fmt.Sprintf("%s %s (%s)", tnode.FIRSTNAMES, tnode.LASTNAME, t),
			ctiles,
		})
	}

	tt := Timetable{
		Title:  "Stundenpläne der Lehrer",
		School: wzdb.Schooldata["SchoolName"].(string),
		Plan:   plan_name,
		Pages:  pages,
	}
	b, err := json.MarshalIndent(tt, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	//os.Stdout.Write(b)
	jsonfile := filepath.Join("_out", "tmp.json")
	jsonpath := filepath.Join(datadir, jsonfile)
	err = os.WriteFile(jsonpath, b, 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Wrote json to: %s\n", jsonpath)
	cmd := exec.Command("typst", "compile",
		"--root", datadir,
		"--input", "ifile="+filepath.Join("..", jsonfile),
		filepath.Join(datadir, "resources", "print_timetable.typ"),
		outpath)
	fmt.Printf(" ::: %s\n", cmd.String())
	//TODO: I am not getting any output messages from typst here ...
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output))
}
