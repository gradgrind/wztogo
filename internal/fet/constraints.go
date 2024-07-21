package fet

import (
	"encoding/xml"
	"gradgrind/wztogo/internal/wzbase"
)

/*
<ConstraintMinDaysBetweenActivities>

	<Weight_Percentage>100</Weight_Percentage>
	<Consecutive_If_Same_Day>true</Consecutive_If_Same_Day>
	<Number_of_Activities>2</Number_of_Activities>
	<Activity_Id>108</Activity_Id>
	<Activity_Id>109</Activity_Id>
	<MinDays>1</MinDays>
	<Active>true</Active>
	<Comments></Comments>

</ConstraintMinDaysBetweenActivities>
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

type minDaysBetweenActivities struct {
	XMLName                 xml.Name `xml:"ConstraintMinDaysBetweenActivities"`
	Weight_Percentage       int
	Consecutive_If_Same_Day bool
	Number_of_Activities    int
	Activity_Id             []int
	MinDays                 int
	Active                  bool
}

func gap_subject_activities(fetinfo *fetInfo,
	subject_activities []wzbase.SubjectGroupActivities,
) {
	gsalist := []minDaysBetweenActivities{}
	for _, sga := range subject_activities {
		l := len(sga.Activities)
		// Adjust indexes for fet
		alist := []int{}
		for _, ai := range sga.Activities {
			alist = append(alist, ai+1)
		}
		gsalist = append(gsalist, minDaysBetweenActivities{
			Weight_Percentage:       100,
			Consecutive_If_Same_Day: true,
			Number_of_Activities:    l,
			Activity_Id:             alist,
			MinDays:                 1,
			Active:                  true,
		})
	}
	fetinfo.fetdata.Time_Constraints_List.ConstraintMinDaysBetweenActivities = gsalist
}
