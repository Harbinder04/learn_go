package jobs

type Job struct {
	Type string	`json:"type"`
	Data any `json:"data"`
}

func NewemailJob(email string) Job {
	return Job{
		Type: "Welcome Email",
		Data: email,
	}
}