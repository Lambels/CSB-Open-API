package engage

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Lambels/go-csb"
)

// engageContext holds all relevant information when making a engage request.
// It is used by the default post method.
type engageContext struct {
	Text             string `json:"Text,omitempty"`
	NumberOfItems    int    `json:"NumberOfItems,omitempty"`
	AcademicYears    string `json:"academicYears,omitempty"`
	ReportingPeriods string `json:"reportingPeriods,omitempty"`
	YearGroupList    string `json:"yearGroupList,omitempty"`
	SubjectList      string `json:"subjectList,omitempty"`
	DivisionList     string `json:"divisionList,omitempty"`
	BatchList        string `json:"batchList,omitempty"`
	PupilIDs         string `json:"pupilIDs"`
}

// engageResponse encapsulates most engage responses.
// It is used by the default post method.
type engageResponse struct {
	D []engageData `json:"d"`
}

type engageData struct {
	Type       string `json:"__type"`
	Text       string `json:"Text"`
	Value      string `json:"Value"`
	Enabled    bool   `json:"Enabled"`
	Attributes struct {
		Checked     bool   `json:"Checked"`
		IsReporting bool   `json:"IsReporting"`
		ColumnType  string `json:"ColumnType"`
	} `json:"Attributes"`
}

// error represents an error message from engage, this message will be parsed to a csb error.
type engageError struct {
	Message       string `json:"Message"`
	StackTrace    string `json:"StackTrace"`
	ExceptionType string `json:"ExceptionType"`
}

// decodeError tries to decode the body of a non OK status request in an error
// understandable by the rest of the api.
func decodeError(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(body) == 0 {
		return csb.HttpStatusErrorf(resp.StatusCode, "engage: error body empty")
	}

	var engErr engageError
	if err := json.Unmarshal(body, &engErr); err != nil {
		return csb.HttpStatusErrorf(resp.StatusCode, "engage: couldnt decode error body")
	}

	return csb.HttpStatusErrorf(resp.StatusCode, "engage: %v , stack trace: %v , exception type: %v", engErr.Message, engErr.StackTrace, engErr.ExceptionType)
}
