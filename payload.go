package hook

import (
	"fmt"
)

// User struct maps the `user` attributes coming from github API
type User struct {
	Name   string `json:"login"`
	ID     int64  `json:"id"`
	URL    string `json:"html_url"`
	Avatar string `json:"avatar_url"`
	Type   string `json:"type"`
}

// PullRequest struct maps the relavent `pull_request` attributes coming from github API
type PullRequest struct {
	URL       string `json:"url"`
	ID        int64  `json:"id"`
	Number    int    `json:"number"`
	State     string `json:"state"`
	Title     string `json:"title"`
	Locked    bool   `json:"locked"`
	User      User   `json:"user"`
	Assignees []User `json:"assignees"`
	Body      string `json:"body"`
	Merged    bool   `json:"merged"`
	Reviewers []User `json:"requested_reviewers"`
}

// Repository struct maps the relavent `repo` attributes coming from github API
type Repository struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Private     bool   `json:"private"`
	Description string `json:"description"`
	Fork        bool   `json:"fork"`
}

// Payload struct maps relavant attributes coming from github API
type Payload struct {
	Action     string      `json:"action"`
	PR         PullRequest `json:"pull_request"`
	Assignee   User        `json:"assignee"`
	Reviewer   User        `json:"requested_reviewer"`
	Sender     User        `json:"sender"`
	Repository Repository  `json:"repository"`
}

func (p *Payload) assignment() string {
	object := p.Assignee.Name
	if p.Sender.Name == p.Assignee.Name {
		object = "himself"
	}

	if p.Action == "assigned" {
		return fmt.Sprintf("%s has %s PR#%d to %s", p.Sender.Name, p.Action, p.PR.Number, object)
	}

	return fmt.Sprintf("%s has %s PR#%d %s", p.Sender.Name, p.Action, p.PR.Number, object)
}

func (p *Payload) review() string {
	var action string
	if p.Action == "review_requested" {
		action = "requested a review to"
	} else {
		action = "removed review request for"
	}

	return fmt.Sprintf("%s has %s %s on PR#%d", p.Sender.Name, action, p.Reviewer.Name, p.PR.Number)
}

func (p *Payload) operation() string {
	if p.Action == "closed" {
		if p.PR.Merged {
			return fmt.Sprintf("%s has merged the PR#%d", p.Sender.Name, p.PR.Number)
		}

		return fmt.Sprintf("%s has closed the PR#%d", p.Sender.Name, p.PR.Number)
	}

	return fmt.Sprintf("%s has %s the PR#%d", p.Sender.Name, p.Action, p.PR.Number)
}

// Process takes care of the title andmessage with proper grammar based on the action
func (p *Payload) Process() (string, string) {
	title := fmt.Sprintf("[%s] %s (#%d)", p.Repository.FullName, p.PR.Title, p.PR.Number)

	switch p.Action {
	case "assigned":
		fallthrough
	case "unassigned":
		return title, p.assignment()
	case "review_requested":
		fallthrough
	case "review_request_removed":
		return title, p.review()
	case "opened":
		fallthrough
	case "closed":
		fallthrough
	case "reopened":
		fallthrough
	case "edited":
		return title, p.operation()
	}
	return "", ""
}
