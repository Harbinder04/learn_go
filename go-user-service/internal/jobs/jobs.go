package jobs

type Job struct {
	Id  string `json:"id"`
	Type string `json:"type"`
	Data string `json:"data"`
}

func NewemailJob(id, data string) Job {
	return Job{
		Id: id,
		Type: "Welcome Email",
		Data: data,
	}
}