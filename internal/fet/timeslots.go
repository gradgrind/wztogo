package fet

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"
)

type Day struct {
	Name string
}

type Days_List struct {
	Number_of_Days int
	Day            []Day
}

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
	var days Days_List
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
