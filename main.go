package main

import (
	"WTU_GO_JobEmailSender/domain"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	totalEmailsSent int
)

func main() {
	log.Println("Iniciando a aplicação...")

	c := cron.New()
	c.AddFunc("@every 1m", reminderJob)
	c.Start()

	go countdownToNextCron()

	time.Sleep(5 * time.Minute)

	select {}
}

func countdownToNextCron() {
	for {
		for i := 3600; i > 0; i-- {
			log.Printf("Próxima execução do job em: %d segundos", i)
			time.Sleep(1 * time.Second)
		}
	}
}

func getUpComingEvents(start, end time.Time) (map[int]domain.EventDetailToEmailJob, error) {
	apiURL := fmt.Sprintf("http://localhost:8080/events/job/get-events?start=%s&end=%s", start.Format(time.RFC3339), end.Format(time.RFC3339))

	log.Printf("Chamando API com URL: %s", apiURL)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("Erro ao chamar a API: %v", err)
	}

	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	log.Printf("Resposta da API: Status Code: %d, Body: %s", resp.StatusCode, bodyString)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Erro ao chamar a API: status code %v", resp.StatusCode)
	}

	var response domain.ApiResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("Erro ao decodar JSON: %v", err)
	}

	return response.Data.EventDetailToEmailJobMap, nil
}

func sendReminderEmail(event domain.EventDetailToEmailJob) {
	log.Printf("Preparando o envio de e-mails para o evento: %s", event.Title)

	from := mail.NewEmail("Wise To Us", "fiapjrv@gmail.com")

	for _, email := range event.UsersEmails {
		log.Printf("Preparando e-mail para %s", email)
		to := mail.NewEmail("Destinatário", email)

		message := mail.NewV3Mail()
		message.SetFrom(from)
		message.SetTemplateID("d-d4459804123446609e7f764be695651e")

		personalization := mail.NewPersonalization()
		personalization.AddTos(to)

		personalization.SetDynamicTemplateData("title", event.Title)
		personalization.SetDynamicTemplateData("data_inicio", event.StartDate)

		message.AddPersonalizations(personalization)

		client := sendgrid.NewSendClient("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")
		response, err := client.Send(message)

		if err != nil {
			log.Printf("Erro ao enviar e-mail para %s: %v", email, err)
		} else if response.StatusCode >= 400 {
			log.Printf("Erro ao enviar e-mail para %s - Status Code: %d, Body: %s", email, response.StatusCode, response.Body)
		} else {
			log.Printf("E-mail enviado para %s - Status Code: %d", email, response.StatusCode)
			log.Printf("Resposta do SendGrid: Headers: %v, Body: %s", response.Headers, response.Body)
			totalEmailsSent++
		}
	}
}

func reminderJob() {
	log.Println("Iniciando o lote de envio de e-mails...")

	totalEmailsSent = 0
	now := time.Now()
	oneHourLater := now.Add(1 * time.Hour)

	events, err := getUpComingEvents(now, oneHourLater)
	if err != nil {
		log.Printf("Erro ao obter eventos: %v", err)
		return
	}

	for eventId, eventDetail := range events {
		log.Printf("Processando evento ID: %d", eventId)
		sendReminderEmail(eventDetail)
	}

	log.Printf("Lote de envio concluído. Total de e-mails enviados: %d", totalEmailsSent)
}
