package fet

import (
	"encoding/xml"
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
	"slices"
	"strconv"
	"strings"
)

type fetRoom struct {
	XMLName                      xml.Name      `xml:"Room"`
	Name                         string        // e.g. k3 ...
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
		rooms = append(rooms, fetRoom{
			Name:     n.ID,
			Capacity: 30000,
			Virtual:  false,
			Comments: n.NAME,
		})
	}
	fetinfo.fetdata.Rooms_List = fetRoomsList{
		Room: rooms,
	}
}

// fet can handle multiple compulsory rooms and choices by using virtual
// rooms. It is not, however, clear how additional ("user-input") rooms
// should be handled. So I will report them and then ignore them.
func addRoomConstraint(fetinfo *fetInfo,
	fixed_rooms *([]fixedRoom),
	room_choices *([]roomChoice),
	activity_id int,
	roomspec wzbase.RoomSpec,
) {
	nrooms := len(roomspec.Compulsory) + len(roomspec.Choices)
	if nrooms == 0 {
		return
	}
	if nrooms == 1 {
		if len(roomspec.Compulsory) == 1 {
			*fixed_rooms = append(*fixed_rooms, fixedRoom{
				Weight_Percentage:  100,
				Activity_Id:        activity_id,
				Room:               fetinfo.ref2fet[roomspec.Compulsory[0]],
				Permanently_Locked: true,
				Active:             true,
			})
			//fetinfo.fetdata.Space_Constraints_List.ConstraintActivityPreferredRoom = cr
		} else {
			rlist := []string{}
			for _, ri := range roomspec.Choices[0] {
				rlist = append(rlist, fetinfo.ref2fet[ri])
			}
			*room_choices = append(*room_choices, roomChoice{
				Weight_Percentage:         100,
				Activity_Id:               activity_id,
				Number_of_Preferred_Rooms: len(rlist),
				Preferred_Room:            rlist,
				Active:                    true,
			})
		}
	} else {
		// Multiple rooms, use a virtual room.
		// Make a "key" for a map to preserve virtual rooms in case the
		// same one is needed more than once.
		allrooms := []string{}
		crooms := make([]int, len(roomspec.Compulsory))
		copy(crooms, roomspec.Compulsory)
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

		//TODO ...
		fmt.Printf("????? allrooms: %s\n", key)
	}

}

/*
<ConstraintActivityPreferredRooms>
	<Weight_Percentage>100</Weight_Percentage>
	<Activity_Id>190</Activity_Id>
	<Number_of_Preferred_Rooms>3</Number_of_Preferred_Rooms>
	<Preferred_Room>EuO</Preferred_Room>
	<Preferred_Room>EuK</Preferred_Room>
	<Preferred_Room>EuU</Preferred_Room>
	<Active>true</Active>
	<Comments></Comments>
</ConstraintActivityPreferredRooms>

<ConstraintActivityPreferredRoom>
	<Weight_Percentage>100</Weight_Percentage>
	<Activity_Id>595</Activity_Id>
	<Room>Sp</Room>
	<Permanently_Locked>true</Permanently_Locked>
	<Active>true</Active>
	<Comments></Comments>
</ConstraintActivityPreferredRoom>
*/

/*
	rooms = append(rooms, fetRoom{
		Name:     "Room1",
		Capacity: 30000,
		Virtual:  false,
		Comments: "First room",
	})
	rooms = append(rooms, fetRoom{
		Name:     "Room2",
		Capacity: 30000,
		Virtual:  false,
		Comments: "Second room",
	})

	rrlist := []realRoomSet{}
	rlist1 := []string{"Room1"}
	rlist2 := []string{"Room2"}

	rrlist = append(rrlist, realRoomSet{
		Number_of_Real_Rooms: 1, Real_Room: rlist1,
	})
	rrlist = append(rrlist, realRoomSet{
		Number_of_Real_Rooms: 1, Real_Room: rlist2,
	})
	rooms = append(rooms, fetRoom{
		Name:                         "V001",
		Capacity:                     30000,
		Virtual:                      true,
		Number_of_Sets_of_Real_Rooms: 2,
		Set_of_Real_Rooms:            rrlist,
		Comments:                     "Virtual room 1",
	})

	fetinfo.fetdata.Rooms_List = fetRoomsList{
		Room: rooms,
	}
}
*/
