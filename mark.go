package csb

import (
	"context"
	"time"
)

// Mark represents the structure of a mark from engage.
type Mark struct {
	// PK of the mark.
	ID int `json:"id"`

	// Links to student.
	StudentID int      `json:"student_id"`
	Student   *Student `json:"student"`

	// Subject at which the mark was recieved.
	SubjectID int     `json:"subject_id"`
	Subject   Subject `json:"subject"`
	// Teacher is the name of the teacher teaching the subject when the mark was recieved.
	Teacher string `json:"teacher"`
	// Percentage represents the grade recieved out of 100.
	Percentage int `json:"percentage"`
	// Exam period on which the mark was recieved.
	Period Period `json:"period"`
	// Timestamp.
	CreatedAt time.Time `json:"created_at"`
}

func (m *Mark) Validate() error {
	if m.StudentID == 0 {
		return Errorf(EINVALID, "validate: mark missing student id field")
	}
	if m.StudentID == 0 {
		return Errorf(EINVALID, "validate: mark missing subject id field")
	}

	ok, err := m.Period.Full()
	if err != nil {
		return err
	}
	if !ok {
		return Errorf(EINVALID, "validate: expecting full period")
	}

	if m.Teacher == "" {
		return Errorf(EINVALID, "validate: expecting teacher field")
	}
	return nil
}

// MarkService represents a mark service.
type MarkService interface {
	// FindMarkByID returns a mark with the id = id.
	//
	// return ENOTFOUND if the mark doesnt exist.
	FindMarkByID(ctx context.Context, id int) (*Mark, error)

	// FindMarksByPID returns the marks of the student with pid = pid.
	//
	// returns ENOTFOUND if the student doesnt exist.
	FindMarksByPID(ctx context.Context, pid int) ([]*Mark, error)

	// FindMarksByPeriod find the marks of the student with pid = pid
	// at a certain examination period.
	//
	// returns ENOTFOUND if the student doesnt exist.
	FindMarksByPeriod(ctx context.Context, pid int, period Period) ([]*Mark, error)

	// FindMarksByPeriodRange finds the marks between the two examination periods.
	// If the periods are narrowed to importance level the filter must provide a PID else
	// EINVALID is returned.
	// If the to period is before the from period, EINVALID is returned.
	FindMarksByPeriodRange(ctx context.Context, from, to Period, filter MarksFilter) ([]*Mark, error)

	// FindMarks finds the marks with the appropiate filter.
	FindMarks(ctx context.Context, filter MarksFilter) ([]*Mark, error)

	// DeleteMark permanently deletes the mark with id = id.
	//
	// returns ENOTFOUND if the mark doesnt exit.
	DeleteMark(ctx context.Context, id int) error

	// RefreshMarks refreshes the marks for a particular student over the exam period span
	// provided.
	//
	// returns any error in the exchange.
	RefreshMarks(ctx context.Context, pid int, from, to Period) error
}

// MarksFilter hardly replicates a RenderMarks request body for engage.
type MarksFilter struct {
	// ID filters on the mark id.
	ID *int `json:"id"`
	// PID filters on the student id.
	PID *int `json:"pid"`
	// Teacher filters on the name of the teacher.
	Teacher *string `json:"teacher"`
	// MinPercentage sets a minimum percentage for the results.
	MinPercentage *int `json:"min_percentage"`
	// MaxPercentage sets a maximum percentage for the results.
	MaxPercentage *int `json:"max_percentage"`
	// Periods filters on the POPULATED period fields on the period fields.
	//
	// If only the academic year is populated then the filter will only be applied on the academic year.
	Periods []Period `json:"periods"`
	// Subjects filters on the marks subjects and only lets through the marks with the specified
	// subjects.
	Subjects []Subject `json:"subjects"`
}
