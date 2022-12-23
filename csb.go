package csb

type Subject string

func (s Subject) String() string {
	v, ok := subjectToString[s]
	if !ok {
		return "Unknown Subject"
	}
	return v
}

const (
	BIOLOGY            Subject = "CL1-106"
	CHEMISTRY          Subject = "CL1-108"
	COMPUTER_SCIENCE   Subject = "CL1-120"
	ECONOMICS          Subject = "CL1-112"
	ENGLISH            Subject = "CL1-102"
	ENGLISH_FIRST      Subject = "CL1-138"
	FRENCH             Subject = "CL1-114"
	FT_ASSEMBLY        Subject = "CL1-154"
	GEOGRAPHY          Subject = "CL1-116"
	HISTORY            Subject = "CL1-119"
	MATHEMATICS        Subject = "CL1-103"
	MUSIC              Subject = "CL1-122"
	PHYSICAL_EDUCATION Subject = "CL1-124"
	PHYSICS            Subject = "CL1-125"
	PSCHEE             Subject = "CL1-126"
	ROMANIAN           Subject = "CL1-128"
	SCIENCE            Subject = "CL1-129"
)

var subjectToString = map[Subject]string{
	BIOLOGY:            "Biology",
	CHEMISTRY:          "Chemistry",
	COMPUTER_SCIENCE:   "Computer Science",
	ECONOMICS:          "Economics",
	ENGLISH:            "English",
	ENGLISH_FIRST:      "English 1st Language",
	FRENCH:             "French",
	FT_ASSEMBLY:        "FT/Assembly",
	GEOGRAPHY:          "Geography",
	HISTORY:            "History",
	MATHEMATICS:        "Mathematics",
	MUSIC:              "Music",
	PHYSICAL_EDUCATION: "Physical Education (PE)",
	PHYSICS:            "Physics",
	PSCHEE:             "PSCHEE",
	ROMANIAN:           "Romanian",
	SCIENCE:            "Science",
}
