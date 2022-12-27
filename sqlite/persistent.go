package sqlite

import (
	"context"
	"database/sql"

	"github.com/Lambels/go-csb"
)

var _ csb.PersistentService = (*PersistentService)(nil)

type PersistentService struct {
	db *DB
}

func NewPersistentService(db *DB) *PersistentService {
	return &PersistentService{
		db: db,
	}
}

func (s *PersistentService) NamesToPIDs(ctx context.Context, names []string) ([][]int, error) {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	pids := make([][]int, len(names))
	for _, name := range names {
		pid, err := nameToPID(ctx, tx, name)
		if err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}

	return pids, nil
}

// RefreshStudents refreshes a stream of students using the same connection for the whole stream.
func (s *PersistentService) RefreshStudents(ctx context.Context, overwrite bool, startPID int, stream <-chan *csb.Student) error {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	lastPID := startPID
	for student := range stream {
		if student.PID != lastPID && overwrite {
			if err := deleteStudents(ctx, tx, lastPID, student.PID); err != nil {
				return err
			}
		}

		_, err := PIDToName(ctx, tx, student.PID)
		switch csb.ErrorCode(err) {
		case csb.ENOTFOUND:
			if err := createStudent(ctx, tx, student); err != nil {
				return err
			}
		case "":

		default:
			return err
		}

		lastPID = student.PID
	}

	return tx.Commit()
}

func (s *PersistentService) DeleteStudentsPID(ctx context.Context, pids []int) error {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	for _, pid := range pids {
		if err := deleteStudent(ctx, tx, pid); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func nameToPIDs(ctx context.Context, tx *sql.Tx, name string) ([]int, error) {

}

func PIDToName(ctx context.Context, tx *sql.Tx, pid int) (string, error) {

}

func createStudent(ctx context.Context, tx *sql.Tx, student *csb.Student) error {

}

// [from, to)
func deleteStudents(ctx context.Context, tx *sql.Tx, from, to int) error {

}

func deleteStudent(ctx context.Context, tx *sql.Tx, pid int) error {

}
