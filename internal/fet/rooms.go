package fet

import (
	"encoding/xml"
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
	"log"
	"slices"
	"strconv"
	"strings"
)

type fetRoom struct {
	XMLName                      xml.Name `xml:"Room"`
	Name                         string   // e.g. k3 ...
	Long_Name                    string
	Capacity                     int           // 30000
	Virtual                      bool          // false`
	Number_of_Sets_of_Real_Rooms int           `xml:",omitempty"`
	Set_of_Real_Rooms            []realRoomSet `xml:",omitempty"`
	Comments                     string
}

type realRoomSet struct {
	Number_of_Real_Rooms int // normally 1, I suppose
	Real_Room            []string
}

type fetRoomsList struct {
	XMLName xml.Name `xml:"Rooms_List"`
	Room    []fetRoom
}

type fixedRoom struct {
	XMLName            xml.Name `xml:"ConstraintActivityPreferredRoom"`
	Weight_Percentage  int
	Activity_Id        int
	Room               string
	Permanently_Locked bool // true
	Active             bool // true
}

type roomChoice struct {
	XMLName                   xml.Name `xml:"ConstraintActivityPreferredRooms"`
	Weight_Percentage         int
	Activity_Id               int
	Number_of_Preferred_Rooms int
	Preferred_Room            []string
	Active                    bool // true
}

// Generate the fet entries for the basic ("real") rooms.
func getRooms(fetinfo *fetInfo) {
	rooms := []fetRoom{}
	for _, ti := range fetinfo.wzdb.TableMap["ROOMS"] {
		n := fetinfo.wzdb.GetNode(ti).(wzbase.Room)
		if len(n.SUBROOMS) == 0 {
			rooms = append(rooms, fetRoom{
				Name:      n.ID,
				Long_Name: n.NAME,
				Capacity:  30000,
				Virtual:   false,
				Comments:  fetinfo.wzdb.SourceReferences[ti],
			})
		} else {
			// Make a virtual room
			rrlist := []realRoomSet{}
			for _, ri := range n.SUBROOMS {
				rrlist = append(rrlist, realRoomSet{
					Number_of_Real_Rooms: 1,
					Real_Room:            []string{fetinfo.ref2fet[ri]},
				})
			}
			rooms = append(rooms, fetRoom{
				Name:                         n.ID,
				Long_Name:                    n.NAME,
				Capacity:                     30000,
				Virtual:                      true,
				Number_of_Sets_of_Real_Rooms: len(rrlist),
				Set_of_Real_Rooms:            rrlist,
				Comments:                     fetinfo.wzdb.SourceReferences[ti],
			})
		}
	}
	fetinfo.fetdata.Rooms_List = fetRoomsList{
		Room: rooms,
	}
}

// fet can handle multiple compulsory rooms and choices by using virtual
// rooms. It is not, however, clear how additional ("user-input") rooms
// should be handled. So I will report them and then ignore them.
func addRoomConstraint(fetinfo *fetInfo,
	//TODO: Is there an advantage to using a fixed room over a roomchoice
	// when there is only one room? Would it be faster? Both variants seem
	// to work. An advantage of only using choices is that in the output
	// file (...data_and_timetable.fet) the generated rooms are then easy
	// to find, being the only fixed ones.
	fixed_rooms *([]fixedRoom),
	room_choices *([]roomChoice),
	virtual_rooms map[string]string,
	activity_indexes []int,
	roomspec wzbase.RoomSpec,
) {
	if roomspec.UserInput != 0 {
		log.Printf("WARNING: 'User-Input' rooms are not supported.")
	}
	nrooms := len(roomspec.Compulsory) + len(roomspec.Choices)
	if nrooms == 0 {
		return
	}
	if nrooms == 1 {
		if len(roomspec.Compulsory) == 1 {
			rm := fetinfo.ref2fet[roomspec.Compulsory[0]]
			for _, ai := range activity_indexes {
				/*
					*fixed_rooms = append(*fixed_rooms, fixedRoom{
						Weight_Percentage:  100,
						Activity_Id:        ai + 1,
						Room:               rm,
						Permanently_Locked: true,
						Active:             true,
					})
				*/
				*room_choices = append(*room_choices, roomChoice{
					Weight_Percentage:         100,
					Activity_Id:               ai + 1,
					Number_of_Preferred_Rooms: 1,
					Preferred_Room:            []string{rm},
					Active:                    true,
				})
			}
		} else {
			// A simple choice
			rlist := []string{}
			for _, ri := range roomspec.Choices[0] {
				rlist = append(rlist, fetinfo.ref2fet[ri])
			}
			for _, ai := range activity_indexes {
				*room_choices = append(*room_choices, roomChoice{
					Weight_Percentage:         100,
					Activity_Id:               ai + 1,
					Number_of_Preferred_Rooms: len(rlist),
					Preferred_Room:            rlist,
					Active:                    true,
				})
			}
		}
	} else {
		// Multiple rooms, use a new virtual room.
		// If a compulsory room is virtual, extract its component rooms.
		compulsory := []int{}
		for _, ri := range roomspec.Compulsory {
			srl := fetinfo.wzdb.GetNode(ri).(wzbase.Room).SUBROOMS
			if len(srl) == 0 {
				compulsory = append(compulsory, ri)
			} else {
				compulsory = append(compulsory, srl...)
			}
		}
		// Make a "key" for a map to preserve virtual rooms in case the
		// same one is needed more than once.
		allrooms := []string{}
		crooms := make([]int, len(compulsory))
		copy(crooms, compulsory)
		slices.Sort(crooms)
		for _, ri := range crooms {
			allrooms = append(allrooms, strconv.Itoa(ri))
		}
		xrooms := []string{}
		for _, ril := range roomspec.Choices {
			rl := []string{}
			slices.Sort(ril)
			for _, ri := range ril {
				rl = append(rl, strconv.Itoa(ri))
			}
			xrooms = append(xrooms, strings.Join(rl, "|"))
		}
		slices.Sort(xrooms)
		allrooms = append(allrooms, xrooms...)
		key := strings.Join(allrooms, "&")
		vr, ok := virtual_rooms[key]
		if !ok {
			// Make virtual room.
			rrlist := []realRoomSet{}
			for _, ri := range roomspec.Compulsory {
				rrlist = append(rrlist, realRoomSet{
					Number_of_Real_Rooms: 1,
					Real_Room:            []string{fetinfo.ref2fet[ri]},
				})
			}
			for _, ril := range roomspec.Choices {
				rl := []string{}
				for _, ri := range ril {
					rl = append(rl, fetinfo.ref2fet[ri])
				}
				rrlist = append(rrlist, realRoomSet{
					Number_of_Real_Rooms: len(rl),
					Real_Room:            rl,
				})
			}
			vr = fmt.Sprintf("v%03d", len(virtual_rooms)+1)
			vroom := fetRoom{
				Name:                         vr,
				Capacity:                     30000,
				Virtual:                      true,
				Number_of_Sets_of_Real_Rooms: len(rrlist),
				Set_of_Real_Rooms:            rrlist,
			}
			// Add the virtual room to the fet file
			fetinfo.fetdata.Rooms_List.Room = append(
				fetinfo.fetdata.Rooms_List.Room, vroom)
			// Remember key/value
			virtual_rooms[key] = vr
		}
		for _, ai := range activity_indexes {
			/*
				*fixed_rooms = append(*fixed_rooms, fixedRoom{
					Weight_Percentage:  100,
					Activity_Id:        ai + 1,
					Room:               vr,
					Permanently_Locked: true,
					Active:             true,
				})
			*/
			*room_choices = append(*room_choices, roomChoice{
				Weight_Percentage:         100,
				Activity_Id:               ai + 1,
				Number_of_Preferred_Rooms: 1,
				Preferred_Room:            []string{vr},
				Active:                    true,
			})
		}
	}
}
