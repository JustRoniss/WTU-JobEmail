package domain

type ApiResponse struct {
	Message   string            `json:"message"`
	Status    int               `json:"status"`
	Links     map[string]string `json:"links"`
	Timestamp string            `json:"timestamp"`
	Data      EventData         `json:"data"`
}

type EventData struct {
	EventDetailToEmailJobMap map[int]EventDetailToEmailJob `json:"eventDetailToEmailJobMap"`
}
