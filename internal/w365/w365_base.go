package w365

import "gradgrind/wztogo/internal/wzbase"

const LIST_SEP = "#" // In the XML dumps it is ","

const w365_Id = "Id"
const w365_ContainerId = "ContainerId"
const w365_ListPosition = "ListPosition"
const w365_Teacher = "Teacher"
const w365_Teachers = "Teachers"
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
const w365_Subjects = "Subjects"
const w365_Room = "Room"
const w365_RoomGroup = "RoomGroup"
const w365_StudentId = "ExtraId"
const w365_Student = "Student"
const w365_Students = "Students"
const w365_Group = "Group"
const w365_Groups = "Groups"
const w365_Firstnames = "Firstname"
const w365_First_Name = "Zusatz 1" // !
const w365_DateOfBirth = "DateOfBirth"
const w365_PlaceOfBirth = "CityOfBirth"
const w365_DateOfEntry = "Zusatz 2" // !
const w365_DateOfExit = "Zusatz 3"  // !
const w365_Home = "City"
const w365_Postcode = "PLZ"
const w365_Street = "Street"
const w365_Email = "Email"
const w365_PhoneNumber = "PhoneNumber"
const w365_Year = "Grade"
const w365_YearDiv = "GradePartiton" // sic
const w365_YearDivs = "GradePartitions"
const w365_Level = "Level"
const w365_Letter = "Letter"
const w365_EpochFactor = "EpochFactor"
const w365_ForceFirstHour = "ForceFirstHour"
const w365_Course = "Course"
const w365_PreferredRooms = "PreferredRooms"
const w365_HandWorkload = "HandWorkload"
const w365_HoursPerWeek = "HoursPerWeek"
const w365_DoubleLessonMode = "DoubleLessonMode"
const w365_EpochWeeks = "EpochWeeks"
const w365_Schedule = "Schedule"
const w365_Lesson = "Lesson"
const w365_Lessons = "Lessons"
const w365_Fixed = "Fixed"
const w365_LocalRooms = "LocalRooms"

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
