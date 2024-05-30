package w365

import "gradgrind/wztogo/internal/wzbase"

const LIST_SEP = "#" // In the XML dumps it is ","

const w365_Id = "Id"
const w365_ContainerId = "ContainerId"
const w365_ListPosition = "ListPosition"
const w365_Teacher = "Teacher"
const w365_Shortcut = "Shortcut"
const w365_Name = "Name"
const w365_Firstname = "Firstname"
const w365_Gender = "Gender"
const w365_MaxDays = "MaxDays"
const w365_MaxLessonsPerDay = "MaxLessonsPerDay"
const w365_MaxGapsPerDay = "MaxWindowsPerDay"
const w365_MinLessonsPerDay = "MinLessonsPerDay"
const w365_NumberOfAfterNoonDays = "NumberOfAfterNoonDays"
const w365_Absence = "Absence"
const w365_Absences = "Absences"
const w365_day = "day"
const w365_Day = "Day"
const w365_hour = "hour"
const w365_Hour = "Hour"
const w365_Period = "TimedObject" // lesson slot
const w365_Start = "Start"
const w365_End = "End"
const w365_MiddayBreak = "MiddayBreak"
const w365_FirstAfternoonHour = "FirstAfternoonHour"
const w365_Category = "Category"
const w365_Categories = "Categories"
const w365_Subject = "Subject"

type ItemType map[string]string

type W365Data struct {
	Schooldata ItemType
	Years      map[string]YearData
	ActiveYear string
	NodeList   []wzbase.WZnode
	NodeMap    map[string]int
	TableMap   map[string][]int
	Config     map[string]interface{}

	tables0     map[string][]ItemType
	yeartables  map[string][]ItemType
	absencemap  map[string]wzbase.Timeslot
	categorymap map[string]Category
}

type YearData struct {
	Tag         string
	Name        string
	DATE_Start  string
	DATE_End    string
	EpochFactor float64
	W365Id      string
	LastChanged string // Datum und Uhrzeit - ISO-Format
}
