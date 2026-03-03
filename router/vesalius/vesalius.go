package vesalius

import (
    "vesaliusdr/controller/vesalius"
	"vesaliusdr/middleware"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(router fiber.Router) {
    api := router.Group("/vesalius")
    api.Use(middleware.JWTProtected)
    api.Get("/patient-data/:branchId/:prn", vesalius.GetPatientData)
    api.Get("/patient-allergy/:branchId/:prn", vesalius.GetPatientAllergy)
    api.Get("/getOutpatientQueueSumarryList/:branchId/:mcr", vesalius.GetOutpatientQueueSummaryList)
    api.Get("/getOutpatientQueueDetailList/:branchId/:queueCriteria/:mcr", vesalius.GetOutpatientQueueDetailList)
    api.Get("/getInpatientQueueDetailList/:branchId/:mcr", vesalius.GetInpatientQueueDetailList)
    api.Get("/getInpatientQueueDetailByPrn/:branchId/:mcr/:prn", vesalius.GetInpatientQueueDetailByPrn)
    api.Get("/get-doctor-todo-notification/:branchId/:mcr", vesalius.GetDoctorToDoNotification)
    api.Get("/get-doctor-todo-notification-details/:branchId/:mcr/:prn", vesalius.GetDoctorToDoNotificationDetails)
    api.Post("/process-doctor-todo-ack/:branchId/:mcr", vesalius.ProcessDoctorToDoAcknowledge)
    api.Post("/vesalius/process-doctor-todo-ack/:branchId/:mcr", vesalius.ProcessDoctorReviewInvestigation)
    api.Get("/get-investigation-report/:branchId/:mcr/:orderDate", vesalius.GetInvestigationReport)
    api.Get("/get-pdf-investigation-report/:branchId/:mcr/:accessionNo", vesalius.GetPdfInvestigationReport)
}