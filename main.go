package main

import (
	"mediasec/repository"
	"mediasec/service"
	"mediasec/ui"
)

func main() {
	sourceRepo := repository.NewMediaSourceRepository()
	pubRepo := repository.NewPublicationRepository()
	incidentRepo := repository.NewIncidentRepository()

	verifySvc := service.NewContentVerificationService(sourceRepo, pubRepo, incidentRepo)
	incidentSvc := service.NewIncidentService(incidentRepo)

	console := ui.NewConsoleUI(sourceRepo, pubRepo, verifySvc, incidentSvc)
	console.Run()
}
