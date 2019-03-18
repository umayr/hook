package hook

import (
	"fmt"
)

type User struct {
	Name   string `json:"login"`
	ID     int64  `json:"id"`
	URL    string `json:"html_url"`
	Avatar string `json:"avatar_url"`
	Type   string `json:"type"`
}

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

type Repository struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Private     bool   `json:"private"`
	Description string `json:"description"`
	Fork        bool   `json:"fork"`
}

type Payload struct {
	Action     string      `json:"action"`
	PR         PullRequest `json:"pull_request"`
	Assignee   User        `json:"assignee"`
	Reviewer   User        `json:"requested_reviewer"`
	Sender     User        `json:"sender"`
	Repository Repository  `json:"repository"`
}

func (p *Payload) Process() (string, string) {
	title := fmt.Sprintf("[%s] %s (#%d)", p.Repository.FullName, p.PR.Title, p.PR.Number)

	if p.Action == "assigned" || p.Action == "unassigned" {
		object := p.Assignee.Name
		if p.Sender.Name == p.Assignee.Name {
			object = "himself"
		}

		if p.Action == "assigned" {
			return title, fmt.Sprintf("%s has %s PR#%d to %s", p.Sender.Name, p.Action, p.PR.Number, object)
		}

		return title, fmt.Sprintf("%s has %s PR#%d %s", p.Sender.Name, p.Action, p.PR.Number, object)
	}

	if p.Action == "review_requested" || p.Action == "review_request_removed" {
		var action string
		if p.Action == "review_requested" {
			action = "requested a review to"
		} else {
			action = "removed review request for"
		}

		return title, fmt.Sprintf("%s has %s %s on PR#%d", p.Sender.Name, action, p.Reviewer.Name, p.PR.Number)
	}

	if p.Action == "opened" || p.Action == "closed" || p.Action == "reopened" || p.Action == "edited" {
		if p.Action == "closed" {
			if p.PR.Merged {
				return title, fmt.Sprintf("%s has merged the PR#%d", p.Sender.Name, p.PR.Number)
			}

			return title, fmt.Sprintf("%s has closed the PR#%d", p.Sender.Name, p.PR.Number)
		}

		return title, fmt.Sprintf("%s has %s the PR#%d", p.Sender.Name, p.Action, p.PR.Number)
	}

	return "", ""
}
