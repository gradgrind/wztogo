package fet

import (
	"encoding/xml"
	"fmt"
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
	XMLName   xml.Name `xml:"Hour"`
	Name      string
	Long_Name string
}

type fetHoursList struct {
	XMLName         xml.Name `xml:"Hours_List"`
	Number_of_Hours int
	Hour            []fetHour
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

func getDays(fetinfo *fetInfo) {
	days := []fetDay{}
	dlist := []string{}
	for _, ti := range fetinfo.wzdb.TableMap["DAYS"] {
		d := fetinfo.ref2fet[ti]
		days = append(days, fetDay{Name: d})
		dlist = append(dlist, d)
	}
	fetinfo.days = dlist
	fetinfo.fetdata.Days_List = fetDaysList{
		Number_of_Days: len(days),
		Day:            days,
	}
}

func getHours(fetinfo *fetInfo) {
	hours := []fetHour{}
	hlist := []string{}
	for _, ti := range fetinfo.wzdb.TableMap["HOURS"] {
		h := fetinfo.ref2fet[ti]
		hn := fetinfo.wzdb.GetNode(ti).(wzbase.Hour)
		hours = append(hours, fetHour{
			Name: hn.ID,
			Long_Name: fmt.Sprintf("%s@%s-%s",
				hn.NAME, hn.START_TIME, hn.END_TIME),
		})
		hlist = append(hlist, h)
	}
	fetinfo.hours = hlist
	fetinfo.fetdata.Hours_List = fetHoursList{
		Number_of_Hours: len(hours),
		Hour:            hours,
	}
}
