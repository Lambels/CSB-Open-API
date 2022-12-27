package csb

// Student represents a student in the most current year
type Student struct {
	// The pupil ID of the student, this field is required since its the only field which
	// makes sense for engage.
	PID int `json:"pid"`
	// The name of the student.
	Name string `json:"name"`
	// CurrentYear represents the current year of the student: "Y11".
	CurrentYear string `json:"current_year"`
	// Subjects are all the subjects the student ever took.
	Subjects []Subject `json:"subjects"`

	// Marks recieved by the student.
	Marks []*Mark `json:"marks"`
}

func (s *Student) Validate() error {
	if s.PID == 0 {
		return Errorf(EINVALID, "student missing PID field")
	}

	return nil
}

// Mark represents the structure of a mark from engage.
type Mark struct {
	// Year represents in what year the student was when the mark was recieved.
	Year string `json:"year"`
	// At represents the subject at which the mark was recieved.
	At Subject `json:"at"`
	// Teacher is the name of the teacher teaching the subject when the mark was recieved.
	Teacher string `json:"teacher"`
	// Percentage represents the grade recieved out of 100.
	Percentage int `json:"percentage"`
	// Column represents the "significance" of the exam: "Term 2 - Periodic Assessment 1"
	// or "Term 2 - Final Assessment (F)"
	Column string `json:"column"`
}

// Subject is the engage identifier of a subject.
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
