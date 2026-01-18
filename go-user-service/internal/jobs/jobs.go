package jobs

type Job struct {
	Type string
	PayLoad string
}

func NewemailJob(email string) Job {
	return Job{
		Type: "Welcome Email",
		PayLoad: email,
	}
}