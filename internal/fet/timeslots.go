package fet

import (
	"encoding/xml"
)

type fetDay struct {
	XMLName xml.Name `xml:"Day"`
	Name    string
}

type fetDaysList struct {
	XMLName        xml.Name `xml:"Days_List"`
	Number_of_Days int
	Day            []fetDay
}

type fetHour struct {
	XMLName xml.Name `xml:"Hour"`
	Name    string
}

type fetHoursList struct {
	XMLName         xml.Name `xml:"Hours_List"`
	Number_of_Hours int
	Day             []fetHour
}

/*
var daysin string = `
<Days_List>
    <Number_of_Days>5</Number_of_Days>
    <Day>
	  <Name>Mo</Name>
    </Day>
    <Day>
	  <Name>Di</Name>
    </Day>
    <Day>
	  <Name>Mi</Name>
    </Day>
    <Day>
	  <Name>Do</Name>
    </Day>
    <Day>
	  <Name>Fr</Name>
    </Day>
</Days_List>
`

func readDays() {
	reader := strings.NewReader(daysin)
	decoder := xml.NewDecoder(reader)
	var days fetDaysList
	err := decoder.Decode(&days)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
		return
	}
	fmt.Printf("\n*** Days: %+v\n", days)

	xmlData, err := xml.MarshalIndent(days, "", "  ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\n*** XML data:\n%v\n", string(xmlData))
}
*/

// func getDays(wzdb *wzbase.WZdata) string {
func getDays(ref2fet map[int]string, trefs []int) fetDaysList {
	days := []fetDay{}
	for _, ti := range trefs {
		days = append(days, fetDay{Name: ref2fet[ti]})
	}
	return fetDaysList{
		Number_of_Days: len(trefs),
		Day:            days,
	}
}

func getHours(ref2fet map[int]string, trefs []int) fetHoursList {
	hours := []fetHour{}
	for _, ti := range trefs {
		hours = append(hours, fetHour{Name: ref2fet[ti]})
	}
	return fetHoursList{
		Number_of_Hours: len(trefs),
		Day:             hours,
	}
}
