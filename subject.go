package csb

import "strings"

// Subject is the engage identifier of a subject.
type Subject string

func (s Subject) String() string {
	v, ok := subjectToString[s]
	if !ok {
		return "Unknown Subject"
	}
	return v
}

// Concat is used to concatonate a list of subjects in a format usable as a parameter
// for an engage request.
func Concat(elems []Subject) string {
	// Code from strings.Join .
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return string(elems[0])
	}
	n := (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(string(elems[0]))
	for _, s := range elems[1:] {
		b.WriteString(",")
		b.WriteString(string(s))
	}
	return b.String()
}

const (
	BIOLOGY                  Subject = "CL1-106"
	CHEMISTRY                Subject = "CL1-108"
	COMPUTER_SCIENCE         Subject = "CL1-120"
	ECONOMICS                Subject = "CL1-112"
	ENGLISH                  Subject = "CL1-102"
	ENGLISH_FIRST            Subject = "CL1-138"
	FRENCH                   Subject = "CL1-114"
	FT_ASSEMBLY              Subject = "CL1-154"
	GEOGRAPHY                Subject = "CL1-116"
	HISTORY                  Subject = "CL1-119"
	MATHEMATICS              Subject = "CL1-103"
	MUSIC                    Subject = "CL1-122"
	PHYSICAL_EDUCATION       Subject = "CL1-124"
	PHYSICS                  Subject = "CL1-125"
	PSCHEE                   Subject = "CL1-126"
	ROMANIAN                 Subject = "CL1-128"
	SCIENCE                  Subject = "CL1-129"
	PHYSICAL_EDUCATION_IGCSE Subject = "CL1-155"
	BUSINESS                 Subject = "CL1-107"
	COMBINED_SCIENCE         Subject = "CL1-153"
	ENGLISH_SECOND           Subject = "CL1-139"
	GLOBAL_PERSPECTIVES      Subject = "CL1-117"
)

var subjectToString = map[Subject]string{
	BIOLOGY:                  "Biology",
	CHEMISTRY:                "Chemistry",
	COMPUTER_SCIENCE:         "Computer Science",
	ECONOMICS:                "Economics",
	ENGLISH:                  "English",
	ENGLISH_FIRST:            "English 1st Language",
	FRENCH:                   "French",
	FT_ASSEMBLY:              "FT/Assembly",
	GEOGRAPHY:                "Geography",
	HISTORY:                  "History",
	MATHEMATICS:              "Mathematics",
	MUSIC:                    "Music",
	PHYSICAL_EDUCATION:       "Physical Education (PE)",
	PHYSICS:                  "Physics",
	PSCHEE:                   "PSCHEE",
	ROMANIAN:                 "Romanian",
	SCIENCE:                  "Science",
	PHYSICAL_EDUCATION_IGCSE: "Physical Eduaction (PE) IGCSE",
	BUSINESS:                 "Business",
	COMBINED_SCIENCE:         "Combined Sciences",
	ENGLISH_SECOND:           "English 2nd Language",
	GLOBAL_PERSPECTIVES:      "Global Perspectives (GP)",
}
