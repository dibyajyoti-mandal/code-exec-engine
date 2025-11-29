package models

type Job struct {
	ID       string `json:"id"`
	ClientID string `json:"client_id"`
	Language string `json:"language"`
	Code     string `json:"code"`
	Image    string `json:"image"`
}

type Result struct {
	JobID    string `json:"job_id"`    // ID to match response to request
	ClientID string `json:"client_id"` // Socket ID to route message
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Error    string `json:"error,omitempty"`
	Status   string `json:"status"`
}
