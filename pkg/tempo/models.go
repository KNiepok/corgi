package tempo

import "time"

type userScheduleResponse struct {
	Schedules []userSchedule `json:"results"`
}

type userSchedule struct {
	Date            string        `json:"date"`
	RequiredSeconds time.Duration `json:"requiredSeconds"`
	Type            string        `json:"type"`
}

type worklogsResponse struct {
	Metadata metadata  `json:"metadata"`
	Worklogs []worklog `json:"results"`
}

type metadata struct {
	Count  int    `json:"count"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Next   string `json:"next"`
}

type issue struct {
	Self string `json:"self"`
	Key  string `json:"key"`
	ID   int    `json:"id"`
}
type worklog struct {
	Self             string        `json:"self"`
	TempoWorklogID   int           `json:"tempoWorklogId"`
	JiraWorklogID    int           `json:"jiraWorklogId"`
	Issue            issue         `json:"issue"`
	TimeSpentSeconds time.Duration `json:"timeSpentSeconds"`
	BillableSeconds  int           `json:"billableSeconds"`
	StartDate        string        `json:"startDate"`
	StartTime        string        `json:"startTime"`
	Description      string        `json:"description"`
	CreatedAt        time.Time     `json:"createdAt"`
	UpdatedAt        time.Time     `json:"updatedAt"`
}
