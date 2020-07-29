package tempo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/jinzhu/now"
	"github.com/kniepok/corgi"
	"net/http"
	"time"
)

const yyyymmddFormat = "2006-01-02"

type NotificationGenerator struct {
	jiraClient *jira.Client

	tempoToken   string
	baseTempoURL string
	tempoClient  *http.Client
}

// todo refactor that to accept jira tempoClient as a constructor param
func NewNotificationGenerator(
	tempoToken,
	jiraToken,
	jiraUsername,
	jiraBaseURL,
	baseTempoURL string) (*NotificationGenerator, error) {
	ng := &NotificationGenerator{
		tempoToken:   tempoToken,
		baseTempoURL: baseTempoURL,
		tempoClient:  http.DefaultClient,
	}

	tp := jira.BasicAuthTransport{
		Username: jiraUsername,
		Password: jiraToken,
	}
	jc, err := jira.NewClient(tp.Client(), jiraBaseURL)
	if err != nil {
		return nil, err
	}

	ng.jiraClient = jc
	return ng, nil
}

func (n *NotificationGenerator) CreateNotification(ctx context.Context, sub corgi.Subscription) (string, error) {
	userAccountID, err := n.getUserAccountID(ctx, sub)
	if err != nil {
		return "", err
	}
	scheduledtime, err := n.getScheduledTime(ctx, sub, userAccountID)
	if err != nil {
		return "", err
	}

	workedTime, err := n.getLoggedTime(ctx, sub, userAccountID)
	if err != nil {
		return "", err
	}

	period := determinePeriod(sub)
	// todo probably rendering the message should belong to some other place.
	// Maybe this package should only return schedule and worklogs and some other place could print pretty messages
	if scheduledtime > workedTime {
		return fmt.Sprintf(`Looks like you logged %s on your Tempo timesheets 
			but you should have logged something around %s.
			Remeber to log the rest before end of the %s!`,
			workedTime,
			scheduledtime,
			period), nil
	}

	return fmt.Sprintf(`You logged %s on your Tempo timesheets on this %s, good job!`, workedTime, period), nil
}

func determinePeriod(sub corgi.Subscription) string {
	switch sub.Details.Mode {
	case corgi.SubModeDaily:
		return "day"
	case corgi.SubModeWeekly:
		return "week"
	default:
		return "day"
	}
}

func (n *NotificationGenerator) getScheduledTime(ctx context.Context, sub corgi.Subscription, accountID string) (time.Duration, error) {
	u := fmt.Sprintf("%s/user-schedule/%s?%s", n.baseTempoURL, accountID, n.getRangeParams(sub.Details))
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", n.tempoToken))
	resp, err := n.tempoClient.Do(req.WithContext(ctx))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("tempo-api: invalid status code for fetching schedule: %d", resp.StatusCode)
	}
	scheduleResponse := new(userScheduleResponse)
	err = json.NewDecoder(resp.Body).Decode(scheduleResponse)
	if err != nil {
		return 0, fmt.Errorf("tempo-api: could not decode response: %w", err)
	}
	var timeScheduled time.Duration
	for _, s := range scheduleResponse.Schedules {
		timeScheduled += time.Second * s.RequiredSeconds
	}

	return timeScheduled, nil
}

func (n *NotificationGenerator) getUserAccountID(ctx context.Context, sub corgi.Subscription) (string, error) {
	u := fmt.Sprintf("/rest/api/2/user/search?query=%s", sub.User.Email)
	req, err := n.jiraClient.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}

	var users []jira.User
	resp, err := n.jiraClient.Do(req.WithContext(ctx), &users)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("jira-api: invalid status code for fetching worklogs: %d", resp.StatusCode)
	}
	if len(users) == 0 {
		return "", fmt.Errorf("jira-api: user with email %s was not found", sub.User.Email)
	}

	return users[0].AccountID, nil
}

func (n *NotificationGenerator) getLoggedTime(ctx context.Context, sub corgi.Subscription, accountID string) (time.Duration, error) {
	u := fmt.Sprintf("%s/worklogs/user/%s?%s&limit=1000", n.baseTempoURL, accountID, n.getRangeParams(sub.Details))
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", n.tempoToken))
	resp, err := n.tempoClient.Do(req.WithContext(ctx))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("tempo-api: invalid status code for fetching worklogs: %d", resp.StatusCode)
	}
	worklogsResponse := new(worklogsResponse)
	err = json.NewDecoder(resp.Body).Decode(worklogsResponse)
	if err != nil {
		return 0, fmt.Errorf("tempo-api: could not decode response: %w", err)
	}
	var timeWorked time.Duration
	for _, s := range worklogsResponse.Worklogs {
		timeWorked += time.Second * s.TimeSpentSeconds
	}

	return timeWorked, nil
}

func (n *NotificationGenerator) getRangeParams(details corgi.SubscriptionDetails) string {
	switch details.Mode {
	case corgi.SubModeDaily:
		today := now.BeginningOfDay().Format(yyyymmddFormat)
		return fmt.Sprintf("from=%s&to=%s", today, today)
	case corgi.SubModeWeekly:
		return fmt.Sprintf("from=%s&to=%s",
			now.BeginningOfWeek().Format(yyyymmddFormat),
			now.EndOfWeek().Format(yyyymmddFormat))
	}
	today := now.BeginningOfDay().Format(yyyymmddFormat)
	return fmt.Sprintf("from=%s&to=%s", today, today)
}
