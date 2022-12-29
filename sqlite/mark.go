package sqlite

import (
	"context"
	"database/sql"

	"github.com/Lambels/go-csb"
	"github.com/Lambels/go-csb/engage"
)

var _ csb.MarkService = (*MarkService)(nil)

type MarkService struct {
	db            *DB
	c             *engage.Client
	periodService csb.PeriodService
	fallback      bool
}

func NewMarkService(db *DB, fallback bool, client *engage.Client, periodService csb.PeriodService) *MarkService {
	return &MarkService{
		db:            db,
		c:             client,
		periodService: periodService,
		fallback:      fallback,
	}
}

func (s *MarkService) FindMarkByID(ctx context.Context, id int) (*csb.Mark, error) {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	mark, err := findMarkByID(ctx, tx, id)
	if err != nil {
		return nil, err
	} else if err := attachMarkAssociations(ctx, tx, mark); err != nil {
		return nil, err
	}

	return mark, nil
}

func (s *MarkService) FindMarksByPID(ctx context.Context, pid int) ([]*csb.Mark, error) {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// find marks by PID only works for local data since we are potentially dealing with allot
	// of marks and we dont want to spam engage. For loading more data use refresh marks.
	marks, err := findMarksByPID(ctx, tx, pid)
	if err != nil {
		return nil, err
	} else if err := attachMarksAssociationsWithStudent(ctx, tx, pid, marks); err != nil {
		return nil, err
	}

	return marks, nil
}

func (s *MarkService) FindMarksByPeriod(ctx context.Context, pid int, period csb.Period) (marks []*csb.Mark, err error) {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	full, err := period.Full()
	if err != nil {
		return nil, err
	}

	if full {
		marks, err = s.findMarksByPeriodFallback(ctx, tx, pid, period)
	} else {
		periods, err := s.periodService.BuildPeriods(ctx, pid, period.AcademicYear, *period.Term)
		if err != nil {
			return nil, err
		}

		marks, err = findMarks(ctx, tx, csb.MarksFilter{PID: &pid, Periods: periods})
	}

	if err != nil {
		return nil, err
	} else if err := attachMarksAssociationsWithStudent(ctx, tx, pid, marks); err != nil {
		return nil, err
	}

	return marks, tx.Commit()
}

func (s *MarkService) FindMarksByPeriodRange(ctx context.Context, from, to csb.Period, filter csb.MarksFilter) (_ []*csb.Mark, err error) {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if filter.PID == nil {
		return nil, csb.Errorf(csb.EINVALID, "cannot generate marks over period range without a student id")
	}
	// if valid id provided save requests to engage for building period range.
	if filter.ID != nil {
		mark, err := findMarkByID(ctx, tx, *filter.ID)
		if err != nil {
			return nil, err
		} else if err := attachMarkAssociations(ctx, tx, mark); err != nil {
			return nil, err
		}

		return []*csb.Mark{mark}, nil
	}

	filter.Periods, err = s.periodService.PeriodRange(ctx, *filter.PID, from, to)
	if err != nil {
		return nil, err
	}

	marks, err := findMarks(ctx, tx, filter)
	if err != nil {
		return nil, err
	} else if err := attachMarksAssociationsWithStudent(ctx, tx, *filter.PID, marks); err != nil {
		return nil, err
	}
	return marks, nil
}

func (s *MarkService) FindMarks(ctx context.Context, filter csb.MarksFilter) ([]*csb.Mark, error) {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	marks, err := findMarks(ctx, tx, filter)
	if err != nil {
		return nil, err
	}

	for _, mark := range marks {
		if err := attachMarkAssociations(ctx, tx, mark); err != nil {
			return nil, err
		}
	}

	return marks, nil
}

func (s *MarkService) DeleteMark(ctx context.Context, id int) error {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteMark(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *MarkService) RefreshMarks(ctx context.Context, pid int, from, to csb.Period) error {
	// asign student manualy since we are again potentially dealing with allot of marks
	// and dont want to spam engage or the local database.
}

func findMarkByID(ctx context.Context, tx *sql.Tx, id int) (*csb.Mark, error) {

}

func findMarksByPID(ctx context.Context, tx *sql.Tx, pid int) ([]*csb.Mark, error) {

}

func (s *MarkService) findMarksByPeriodFallback(ctx context.Context, tx *sql.Tx, pid int, period csb.Period) ([]*csb.Mark, error) {

}

func findMarksByPeriod(ctx context.Context, tx *sql.Tx, pid int, period csb.Period) ([]*csb.Mark, error) {

}

func (s *MarkService) findMarksByPeriodEngage(ctx context.Context, pid int, period csb.Period) ([]*csb.Period, error) {

}

func findMarks(ctx context.Context, tx *sql.Tx, filter csb.MarksFilter) ([]*csb.Mark, error) {

}

func deleteMark(ctx context.Context, tx *sql.Tx, id int) error {

}

func attachMarkAssociations(ctx context.Context, tx *sql.Tx, mark *csb.Mark) (err error) {

}

func attachMarksAssociationsWithStudent(ctx context.Context, tx *sql.Tx, pid int, marks []*csb.Mark) (err error) {

}

func attachMarkSubject(ctx context.Context, tx *sql.Tx, mark *csb.Mark) (err error) {

}
