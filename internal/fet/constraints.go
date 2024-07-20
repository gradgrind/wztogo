package fet

import (
	"encoding/xml"
)

/*
<ConstraintActivityPreferredStartingTime>

	<Weight_Percentage>100</Weight_Percentage>
	<Activity_Id>126</Activity_Id>
	<Preferred_Day>Do</Preferred_Day>
	<Preferred_Hour>A</Preferred_Hour>
	<Permanently_Locked>true</Permanently_Locked>
	<Active>true</Active>
	<Comments></Comments>

</ConstraintActivityPreferredStartingTime>
*/
type startingTime struct {
	XMLName            xml.Name `xml:"ConstraintActivityPreferredStartingTime"`
	Weight_Percentage  int
	Activity_Id        int
	Preferred_Day      string
	Preferred_Hour     string
	Permanently_Locked bool
	Active             bool
}
