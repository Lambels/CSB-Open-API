package engage

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Lambels/go-csb"
)

var (
	baseURL               = "https://cambridgeschoolportal.engagehosted.com/Services/ReportCommentServices.asmx/"
	academicYearsURL      = "GetMarksheetAcademicYears"
	reportingPeriodsURL   = "GetReportingPeriods"
	reportingSubjectsURL  = "GetPupilMarksheetSubjects"
	columnsForSubjectsURL = "GetColumnsForSubjects"
	marksheetRenderURL    = "RenderPupilMarksheet"
)

// Client is a client used to interface with the engage api.
type Client struct {
	cc *http.Client
}

// GetAcademicYears gets all the possible academic years for a PID.
func (c *Client) GetAcademicYears(ctx context.Context, pid string) ([]string, error) {
	resURL := baseURL + academicYearsURL

	res, err := c.post(ctx, resURL, engageContext{PupilIDs: pid})
	if err != nil {
		return nil, err
	}

	out := make([]string, len(res.D))
	for _, data := range res.D {
		out = append(out, data.Value)
	}

	return out, nil
}

// GetReportingPeriods gets the reporting periods for a PID in a specific range of academic years.
func (c *Client) GetReportingPeriods(ctx context.Context, pid string, academicYears []string) ([]string, error) {
	resURL := baseURL + reportingPeriodsURL

	res, err := c.post(ctx, resURL, engageContext{
		PupilIDs:      pid,
		AcademicYears: strings.Join(academicYears, ","),
	})
	if err != nil {
		return nil, err
	}

	out := make([]string, len(res.D))
	for _, data := range res.D {
		out = append(out, data.Value)
	}

	return out, nil
}

// GetReportingSubjects gets the reporting subjects for a PID in a specific range of academic years and reporting periods.
func (c *Client) GetReportingSubjects(ctx context.Context, pid string, academicYears, reportingPeriods []string) ([]csb.Subject, error) {
	resURL := baseURL + reportingSubjectsURL

	res, err := c.post(ctx, resURL, engageContext{
		PupilIDs:         pid,
		AcademicYears:    strings.Join(academicYears, ","),
		ReportingPeriods: strings.Join(reportingPeriods, ","),
	})
	if err != nil {
		return nil, err
	}

	out := make([]csb.Subject, len(res.D))
	for _, data := range res.D {
		out = append(out, csb.Subject(data.Value))
	}

	return out, nil
}

// GetColumnsForSubjects gets the "columns" for a pid in a specified academic years and periods range for the specified subjects.
// A column refers to the type of exam.
func (c *Client) GetColumnsForSubjects(ctx context.Context, pid string, academicYears, reportingPeriods []string, subjects []csb.Subject) ([]string, error) {
	resURL := baseURL + reportingPeriodsURL

	res, err := c.post(ctx, resURL, engageContext{
		PupilIDs:         pid,
		AcademicYears:    strings.Join(academicYears, ","),
		ReportingPeriods: strings.Join(reportingPeriods, ","),
		SubjectList:      csb.Concat(subjects),
	})
	if err != nil {
		return nil, err
	}

	out := make([]string, len(res.D))
	for _, data := range res.D {
		out = append(out, data.Value)
	}

	return out, nil
}

func (c *Client) GetMarksheetRender(ctx context.Context, pid string, academicYears, subjectColumns, reportingPeriods []string, reportingSubjects []csb.Subject) (string, error) {
	return "", nil
}

// post sends a post request to url with the specified engage context. It checks for any errors during
// the exchange process with engage. It returns an engage response which has
// at least one piece of data inside.
func (c *Client) post(ctx context.Context, url string, engCtx engageContext) (res *engageResponse, err error) {
	body, err := json.Marshal(engCtx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.cc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	// parse response here for status code not found since engage is very wierd with response
	// codes. (an invalid pid results in StatusCodeOK)
	if len(res.D) == 0 {
		return nil, csb.Errorf(csb.ENOTFOUND, "engage: invalid PID: %v", engCtx.PupilIDs)
	}

	return res, nil
}

// NewClient creates a new engage client with the provided token used for
// authentification.
func NewClient(c *http.Client, token string) *Client {
	c.Transport = &cookieHeaderTransport{
		cookie: token,
		d:      c.Transport,
	}

	return &Client{
		cc: c,
	}
}

type cookieHeaderTransport struct {
	cookie string
	d      http.RoundTripper
}

func (t *cookieHeaderTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("Cookie", t.cookie)
	return t.d.RoundTrip(r)
}
