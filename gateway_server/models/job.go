package models

type Job struct {
	ID        string `json:"id"`
	ClientID  string `json:"client_id"`
	Language  string `json:"language"`
	Code      string `json:"code"`
	Timestamp int64  `json:"timestamp"`
}
