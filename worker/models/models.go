package models

type Job struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Image    string `json:"image"`
}

type Result struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  string `json:"error"`
}
