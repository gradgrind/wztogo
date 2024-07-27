package fet

import (
	"encoding/xml"
	"gradgrind/wztogo/internal/wzbase"
	"log"
)

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
		// Skip if all activities are "fixed".
		allfixed := true
		for _, ai := range sga.Activities {
			alist = append(alist, ai+1)
			if !fetinfo.fixed_activities[ai] {
				allfixed = false
			}
		}
		if allfixed {
			continue
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

type lunchBreak struct {
	XMLName             xml.Name `xml:"ConstraintStudentsSetMaxHoursDailyInInterval"`
	Weight_Percentage   int
	Students            string
	Interval_Start_Hour string
	Interval_End_Hour   string
	Maximum_Hours_Daily int
	Active              bool
}

func lunch_break(
	fetinfo *fetInfo,
	lbconstraints *([]lunchBreak),
	cname string,
	lunchperiods []int,
) bool {
	// Assume the lunch periods are sorted, but not necessarily contiguous,
	// which is necessary for this constraint.
	lb1 := lunchperiods[0]
	lb2 := lunchperiods[len(lunchperiods)-1] + 1
	if lb2-lb1 != len(lunchperiods) {
		log.Printf(
			"\n=========================================\n"+
				"  !!!  INCOMPATIBLE DATA: lunch periods not contiguous,\n"+
				"       can't generate lunch-break constraint for class %s.\n"+
				"=========================================\n",
			cname)
		return false
	}
	lb := lunchBreak{
		Weight_Percentage:   100,
		Students:            cname,
		Interval_Start_Hour: fetinfo.hours[lb1],
		Interval_End_Hour:   fetinfo.hours[lb2],
		Maximum_Hours_Daily: len(lunchperiods) - 1,
		Active:              true,
	}
	*lbconstraints = append(*lbconstraints, lb)
	return true
}

type maxGapsPerWeek struct {
	XMLName           xml.Name `xml:"ConstraintStudentsSetMaxGapsPerWeek"`
	Weight_Percentage int
	Max_Gaps          int
	Students          string
	Active            bool
}

type minLessonsPerDay struct {
	XMLName             xml.Name `xml:"ConstraintStudentsSetMinHoursDaily"`
	Weight_Percentage   int
	Minimum_Hours_Daily int
	Students            string
	Allow_Empty_Days    bool
	Active              bool
}
