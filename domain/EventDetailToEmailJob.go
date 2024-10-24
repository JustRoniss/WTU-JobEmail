package domain

type EventDetailToEmailJob struct {
	Title       string   `json:"title"`
	StartDate   string   `json:"start"`       // Corrigido para corresponder ao campo "start" da API
	UsersEmails []string `json:"usersEmails"` // Corrigido para corresponder ao campo "usersEmails" da API
}
