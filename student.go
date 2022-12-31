package csb

import (
	"context"
	"time"
)

// Student represents a student from engage.
type Student struct {
	// The pupil ID of the student, this field is required since its the only field which
	// makes sense for engage.
	PID int `json:"pid"`
	// The name of the student.
	Name string `json:"name"`
	// CurrentYear represents the current year of the student: 11 -> Y11.
	CurrentYear int `json:"current_year"`
	// Indicates if this student still attends the school, if false, current year will empty.
	AttendsSchool bool `json:"attends_school"`
	// Subjects are all the subjects the student ever took.
	Subjects []Subject `json:"subjects"`
	// Marks are all the marks the student ever took.
	Marks []*Mark `json:"marks"`
	// Timestamps.
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Student) Validate() error {
	if s.Name == "" {
		return Errorf(EINVALID, "validate: student missing name field")
	}
	if len(s.Subjects) == 0 {
		return Errorf(EINVALID, "validate: student has no subjects")
	}

	return nil
}

// StudentService represents a student service.
type StudentService interface {
	// FindStudentByPID returns the student with pid = pid.
	//
	// returns ENOTFOUND if the student doesent exist.
	FindStudentByPID(ctx context.Context, pid int) (*Student, error)

	// FindStudents finds the marks with the appropiate filter.
	FindStudents(ctx context.Context, filter StudentFilter) ([]*Student, error)

	// DeleteStudent permanently deletes the student with pid = pid.
	//
	// returns ENOTFOUND if the student doesnt exist.
	DeleteStudent(ctx context.Context, pid int) error

	// RefreshStudents refreshes the students in the system with the provided refresh students
	// filter.
	//
	// returns any error in the exchange.
	RefreshStudents(ctx context.Context, refresh RefreshStudents) error
}

// StudentFilter represents a filter to bulk get students.
type StudentFilter struct {
	// PID filters on the student pid.
	PID *int `json:"pid"`

	// Name filters on the students names.
	Name *string `json:"name"`

	// CurrentYear filters on the students current year.
	CurrentYear *int `json:"current_year"`

	// AttendsSchool filters wether the students currently attend school.
	AttendsSchool *bool `json:"attends_school"`

	// Subjects filters on the subjects each student takes.
	Subjects *[]Subject `json:"subjects"`
}

// RefreshStudents represents an request to the RefreshStudents serivce.
type RefreshStudents struct {
	// StartPID is the user you want to start refreshing from (including).
	StartPID int `json:"startPID"`
	// N is the amount of users you want to refresh.
	N int `json:"n"`
	// Purge indicates wether removed users should be deleted or not.
	Purge bool `json:"purge"`
}
