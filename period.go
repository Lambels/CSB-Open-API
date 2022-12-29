package csb

import "context"

// Period represents a period of examination, the academic year, term and importance of the exam
// at which the mark was recieved.
type Period struct {
	AcademicYear int  `json:"academic_year"`
	Term         *int `json:"term"`
	// Importance differs from student to student since some students take exams in other periods.
	Importance *string `json:"importance"`
}

func (p Period) Full() (bool, error) {
	if err := p.Validate(); err != nil {
		return false, err
	}

	if p.Term == nil || p.Importance == nil {
		return false, nil
	}
	return true, nil
}

func (p Period) Validate() error {
	if p.AcademicYear < 2020 {
		return Errorf(EINVALID, "validate: period has invalid academic year: %v", p.AcademicYear)
	}
	if p.Importance != nil && p.Term == nil {
		return Errorf(EINVALID, "validate: period cannot have importance field without term field")
	}
	if p.Term == nil {
		return nil
	}

	if *p.Term > 4 || *p.Term < 1 {
		return Errorf(EINVALID, "validate: term must be between 1 and 4 inclusive, but got: %v", p.Term)
	}
	return nil
}

// PeriodService represents a period service.
//
// PeriodService should usually be implemented over engage since periods are volatile and
// unpredictable.
type PeriodService interface {
	// BuildPeriods builds underlying periods from the base academicYear and optional term.
	//
	// If term is provided and pid isnt EINVALID is returned.
	BuildPeriods(ctx context.Context, pid int, academicYear int, term int) ([]Period, error)

	// Exists checks wether a period exists.
	Exists(ctx context.Context, pid int, period Period) (bool, error)

	// PeriodRange generates a range of periods [from, to].
	//
	// If pid is provided, the returned periods will have accuracy to importance level.
	//
	// If pid isnt provided but one of the periods has the importance level populated the
	// field will be ignored and the term field will be used as the most accurate selection.
	PeriodRange(ctx context.Context, pid int, from, to Period) ([]Period, error)
}
