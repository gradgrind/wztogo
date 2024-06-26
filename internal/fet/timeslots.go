package fet

import (
	"encoding/xml"
	"gradgrind/wztogo/internal/wzbase"
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

func getDays(wzdb *wzbase.WZdata) string {
	trefs := wzdb.TableMap["DAYS"]
	days := []fetDay{}
	for _, ti := range trefs {
		n := wzdb.NodeList[wzdb.IndexMap[ti]].Node
		days = append(days, fetDay{Name: n.(wzbase.Day).ID})
	}
	data := fetDaysList{
		Number_of_Days: len(trefs),
		Day:            days,
	}
	return string(makeXML(data, 0))
}

func getHours(wzdb *wzbase.WZdata) string {
	trefs := wzdb.TableMap["HOURS"]
	hours := []fetHour{}
	for _, ti := range trefs {
		n := wzdb.NodeList[wzdb.IndexMap[ti]].Node
		hours = append(hours, fetHour{Name: n.(wzbase.Hour).ID})
	}
	data := fetHoursList{
		Number_of_Hours: len(trefs),
		Day:             hours,
	}
	return string(makeXML(data, 0))
}
