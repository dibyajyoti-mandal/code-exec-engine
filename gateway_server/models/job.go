package models

type Job struct {
	ID        string `json:"id"`
	Language  string `json:"language"`
	Code      string `json:"code"`
	Timestamp int64  `json:"timestamp"`
}
