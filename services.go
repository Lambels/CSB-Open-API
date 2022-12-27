package csb

import "context"

// StudentFilter represents a filter to bulk get students.
type StudentFilter struct {
	// PIDs represents the pupil ids of the students.
	PIDs []int `json:"pIDs"`
	// Names represent the names which will be fetched and transformed to pupil IDs
	// by the db layer.
	Names []string `json:"names"`
}

// MarksFilter hardly replicates a RenderMarks request body for engage.
type MarksFilter struct {
	// AcademicYears on which to get marks.
	AcademicYears []string `json:"academic_years"`
	// ReportingPeriods over which the marks span, these have to be included in the academic
	// years.
	ReportingPeriods []string `json:"reporting_periods"`
	// Columns on when the marks have been recieved, these have to be included in the reporting
	// periods.
	Columns []string `json:"columns"`
	// Subjects for the marks, these have to be taken by the student.
	Subjects []Subject `json:"subjects"`
}

// RefreshStudents represents an request to the RefreshStudents serivce.
type RefreshStudents struct {
	// StartPID is the user you want to start refreshing from (including).
	StartPID *int `json:"startPID"`
	// N is the amount of users you want to refresh.
	N *int `json:"n"`
	// Overwrite indicates wether removed users should be deleted or not.
	Overwrite bool `json:"overwrite"`
}

// StudentService interfaces with engage and the db layer when needed to fetch students.
type StudentService interface {
	// GetStudents fetches students from the db layer (when the filter uses names) and then
	// fetches using ids from engage the students in their current form.
	GetStudents(context.Context, StudentFilter) ([]*Student, error)

	// RefreshStudents updates students with the provided refresh students request.
	RefreshStudents(context.Context, RefreshStudents) error

	// DeleteStudentPID deletes a student from the db layer via the pID.
	DeleteStudentPID(context.Context, int) error

	// GetStudentMarks fetches the marks for the student provided and populates his marks field.
	GetStudentMarks(context.Context, *Student, MarksFilter) error
}

// PersistentService represents a persistent backup of names paired with pids.
//
// The student service usually runs on top of this.
type PersistentService interface {
	// NamesToPIDs converts names to pupil ids.
	// If a name yields multiple pids all of them are returned.
	//
	// If the returned values are 1:1 the error will explain why.
	NamesToPIDs(context.Context, []string) ([][]int, error)

	// RefreshStudents refreshes a stream of sequentiall users.
	//
	// If any error is encountered an error is returned and refresh students stops reading
	// from the stream.
	RefreshStudents(context.Context, bool, int, <-chan *Student) error

	// DeleteStudentPID deletes students with the provided PIDs.
	DeleteStudentsPID(context.Context, []int) error
}
