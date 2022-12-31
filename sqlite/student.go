package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	csb "github.com/Lambels/CSB-Open-API"
	"github.com/Lambels/CSB-Open-API/engage"
)

var _ csb.StudentService = (*StudentService)(nil)

// StudentService wraps around an engage client.
type StudentService struct {
	// db for persistance.
	db *DB
	// client for updates.
	c *engage.Client
	// saveNew indicates wether fetch to new students should be saved.
	saveNew bool
}

// NewStudentService creates a new student service with the provided database and engage client.
func NewStudentService(db *DB, client *engage.Client, saveNew bool) *StudentService {
	return &StudentService{
		db:      db,
		c:       client,
		saveNew: saveNew,
	}
}

// FindStudentByPID returns a student based on the passed pid.
//
// If the student isnt originally found in the database, the service will try to search
// engage using the engage client, if the user isnt found, ultimately ENOTFOUND is returned.
//
// If the user is found in engage and not in the db and saveNew is true the
// user is saved before returned.
func (s *StudentService) FindStudentByPID(ctx context.Context, pid int) (*csb.Student, error) {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	student, err := findStudentByPID(ctx, tx, pid)
	switch csb.ErrorCode(err) {
	case "":
		if err := attachStudentMarks(ctx, tx, student); err != nil {
			return nil, err
		}

		return student, nil

	case csb.ENOTFOUND:
		student, err := s.findStudentByPIDEngage(ctx, pid)
		if err != nil {
			return nil, err
		}
		if !s.saveNew {
			return student, nil
		}

		if err := createStudent(ctx, tx, student); err != nil {
			return student, nil
		}
		return student, tx.Commit()

	default:
		return nil, err
	}
}

// FindStudents returns a range of students based on the filter.
func (s *StudentService) FindStudents(ctx context.Context, filter csb.StudentFilter) ([]*csb.Student, error) {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	students, err := findStudents(ctx, tx, filter)
	if err != nil {
		return nil, err
	}

	for _, student := range students {
		if err := attachStudentMarks(ctx, tx, student); err != nil {
			return nil, err
		}
	}

	return students, nil
}

// DeleteStudent permanently deletes a student specified by pid.
// returns ENOTFOUND if student isnt found.
func (s *StudentService) DeleteStudent(ctx context.Context, pid int) error {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	if err := deleteStudent(ctx, tx, pid); err != nil {
		return err
	}

	return tx.Commit()
}

// RefreshStudents refreshes students incrementally starting from refresh.StartPID, refresh.N
// times.
//
// If a student is in engage but not in the local copy of students then the student is added
// to the local copy. If refresh.Purge is set to true and the student from engage is not
// attending school then the copy from engage to local storage wont be made.
//
// If refresh.Purge is set to true and a student in the local database is not attending the school
// any more in engage, then the user is deleted.
//
// If the student is both in engage and local storage, an update will be so that your local
// storage has the newest data.
func (s *StudentService) RefreshStudents(ctx context.Context, refresh csb.RefreshStudents) error {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	PIDCount := refresh.StartPID
	for i := 0; i < refresh.N; i++ {
		// engage copy.
		studentEngage, err := s.findStudentByPIDEngage(ctx, PIDCount)
		if err != nil && csb.ErrorCode(err) != csb.ENOTFOUND {
			return err
		}

		// local copy.
		studentLocal, err := findStudentByPID(ctx, tx, PIDCount)
		if err != nil && csb.ErrorCode(err) != csb.ENOTFOUND {
			return err
		}

		switch {
		case studentEngage == nil && studentLocal == nil:
			// no data from engage or local db.
		case !studentEngage.AttendsSchool && studentLocal != nil && refresh.Purge:
			// old student in db and willing to purge.
			if err := deleteStudent(ctx, tx, PIDCount); err != nil {
				return err
			}
		case studentEngage != nil && studentLocal == nil:
			// engage ahead of local db.
			// if engage is ahead of local db with students who dont attend the school
			// and this request is actively purgeing, skip the creation.
			if !studentEngage.AttendsSchool && refresh.Purge {
				break
			}

			if err := createStudent(ctx, tx, studentEngage); err != nil {
				return err
			}
		case studentEngage != nil && studentLocal != nil:
			// data from both engage and local db, update local db.
			if err := updateStudent(ctx, tx, studentLocal.PID, studentEngage); err != nil {
				return err
			}
		}

		PIDCount++

		// dont spam engage.
		select {
		case <-ctx.Done():
			return fmt.Errorf("refresh students: %w", ctx.Err())
		case <-time.After(engage.RequestTimeout):
		}
	}

	return nil
}

func findStudentByPID(ctx context.Context, tx *sql.Tx, id int) (*csb.Student, error) {

}

func (s *StudentService) findStudentByPIDEngage(ctx context.Context, id int) (*csb.Student, error) {

}

func findStudents(ctx context.Context, tx *sql.Tx, filter csb.StudentFilter) ([]*csb.Student, error) {
}

func createStudent(ctx context.Context, tx *sql.Tx, student *csb.Student) error {

}

func updateStudent(ctx context.Context, tx *sql.Tx, id int, student *csb.Student) error {

}

func deleteStudent(ctx context.Context, tx *sql.Tx, pid int) error {

}

func attachStudentMarks(ctx context.Context, tx *sql.Tx, student *csb.Student) (err error) {
	if student.Marks, err = findMarksByPID(ctx, tx, student.PID); err != nil {
		fmt.Errorf("attach student marks: %w", err)
	}
	return nil
}
