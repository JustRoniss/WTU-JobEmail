package domain

type EventEmailJobResponse struct {
	EventDetailToEmailJobMap map[int]EventDetailToEmailJob `json:"eventDetailToEmailJobMap"`
}
