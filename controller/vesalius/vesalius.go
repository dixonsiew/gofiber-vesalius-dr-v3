package vesalius

import (
    "encoding/base64"
    "math"
    "strings"
    "vesaliusdr/config"
    "vesaliusdr/dto"
    "vesaliusdr/middleware"
    "vesaliusdr/model"
    doctorRequestService "vesaliusdr/service/doctor_request"
    doctorTodoNotificationService "vesaliusdr/service/doctor_todo_notification"
    inpatientQueueListService "vesaliusdr/service/inpatient_queue_list"
    investigationReportService "vesaliusdr/service/investigation_report"
    novaPatientAlertService "vesaliusdr/service/nova_patient_alert"
    outpatientQueueListService "vesaliusdr/service/outpatient_queue_list"
    patientInfoService "vesaliusdr/service/patient_info"
    "vesaliusdr/utils"
    "vesaliusdr/ws"

    "vesaliusdr/xmlmodel"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v3"
    "github.com/nleeper/goment"
)

// GetPatientData
//
// @Tags Vesalius
// @Produce json
// @Param        branchId   path      string  false "BranchId"
// @Param        prn        path      string  true  "PRN"
// @Security BearerAuth
// @Success 200 {object} xmlmodel.Patient
// @Router /vesalius/patient-data/{branchId}/{prn} [get]
func GetPatientData(c fiber.Ctx) error {
    prn := c.Params("prn")
    _, _, err := middleware.ValidateToken(c)
    if err != nil {
        return middleware.NoContent(c)
    }

    patientInfo, err := patientInfoService.FindByPrn(prn)
    if err != nil {
        return middleware.NoContent(c)
    }

    if patientInfo != nil {
        o := xmlmodel.Patient{
            Prn: prn,
            ContactNumber: xmlmodel.ContactNumber{
                Home:  patientInfo.Contact_number,
                Email: patientInfo.Email,
            },
            DOB:      patientInfo.DOB,
            Resident: patientInfo.Resident,
            HomeAddress: xmlmodel.HomeAddress{
                Address1:   patientInfo.Home_address1,
                Address2:   patientInfo.Home_address2,
                Address3:   patientInfo.Home_address3,
                CityState:  patientInfo.Home_address4,
                PostalCode: patientInfo.Home_address5,
            },
            Name: xmlmodel.Name{
                Title:      patientInfo.Title,
                FirstName:  patientInfo.First_name,
                MiddleName: patientInfo.Middle_name,
                LastName:   patientInfo.Last_name,
            },
            Nationality: xmlmodel.Nationality{
                Code:        patientInfo.Nationalityid,
                Description: patientInfo.Nationality_desc,
            },
            Sex: xmlmodel.Sex{
                Code:        patientInfo.Sex_code,
                Description: patientInfo.Sex_desc,
            },
            Document: []xmlmodel.Document{
                {
                    Code:        patientInfo.Document_code,
                    Description: patientInfo.Document_desc,
                    Value:       patientInfo.Document_value,
                },
            },
        }
        return c.JSON(o)
    } else {
        o, ex := ws.GetPatientData(prn)
        if ex != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "statusCode": fiber.StatusNotFound,
                "message":    "Information provided does not match with hospital patient profile. Please retry.",
            })
        }

        if o != nil {
            patientInfo := model.PatientInformation{
                Prn:              o.Prn,
                Contact_number:   o.ContactNumber.Home,
                Email:            o.ContactNumber.Email,
                DOB:              o.DOB,
                Resident:         o.Resident,
                Home_address1:    o.HomeAddress.Address1,
                Home_address2:    o.HomeAddress.Address2,
                Home_address3:    o.HomeAddress.Address3,
                Home_address4:    o.HomeAddress.CityState,
                Home_address5:    o.HomeAddress.PostalCode,
                Title:            o.Name.Title,
                First_name:       o.Name.FirstName,
                Middle_name:      o.Name.MiddleName,
                Last_name:        o.Name.LastName,
                Nationalityid:    o.Nationality.Code,
                Nationality_desc: o.Nationality.Description,
                Sex_code:         o.Sex.Code,
                Sex_desc:         o.Sex.Description,
                Last_update_date: localDateTime(),
            }
            for _, doc := range o.Document {
                if strings.ToUpper(doc.Description) == "NRIC / PASSPORT" {
                    patientInfo.Document_code = doc.Code
                    patientInfo.Document_desc = doc.Description
                    patientInfo.Document_value = doc.Value
                }
            }
            patientInfoService.AddPatientInfo(patientInfo)
            return c.JSON(o)
        }

        return middleware.NoContent(c)
    }
}

// GetPatientAllergy
//
// @Tags Vesalius
// @Produce json
// @Param        branchId   path      string  false "BranchId"
// @Param        prn        path      string  true  "PRN"
// @Security BearerAuth
// @Success 200 {array} model.NovaPatientAlert
// @Router /vesalius/patient-allergy/{branchId}/{prn} [get]
func GetPatientAllergy(c fiber.Ctx) error {
    prn := c.Params("prn")
    _, _, err := middleware.ValidateToken(c)
    if err != nil {
        return middleware.NoContent(c)
    }

    lx, err := novaPatientAlertService.FindPatientActiveAlertByPrn(prn)
    if err != nil {
        return middleware.NoContent(c)
    }

    if len(lx) > 0 {
        return c.JSON(lx)
    }

    return middleware.NoContent(c)
}

// GetOutpatientQueueSummaryList
//
// @Tags Vesalius
// @Produce json
// @Param        branchId   path      string  false "BranchId"
// @Param        mcr        path      string  true  "MCR"
// @Security BearerAuth
// @Success 200 {array} model.OutpatientQueueSummary
// @Router /vesalius/getOutpatientQueueSumarryList/{branchId}/{mcr} [get]
func GetOutpatientQueueSummaryList(c fiber.Ctx) error {
    mcr := c.Params("mcr")
    _, _, err := middleware.ValidateToken(c)
    if err != nil {
        return middleware.NoContent(c)
    }

    queueCriteriaReg := "Registered"
    queueCriteriaAppt := "Appointment"
    queueCriteriaSeen := "Seen"
    queueCriteriaKiv := "KIV"

    regCount := 0
    apptCount := 0
    seenCount := 0
    kivCount := 0

    lastUpdateDateFromTable, err := outpatientQueueListService.GetLastUpdateDateByMCR(mcr)
    if err != nil {
        return middleware.NoContent(c)
    }

    if lastUpdateDateFromTable != "" {
        localLastUpdateDate := stringToGoment(lastUpdateDateFromTable)
        minutesDiff := getDateTimeDifferent(localLastUpdateDate)
        vesaliusWSInterval := config.VesaliusWSInterval()
        if minutesDiff > vesaliusWSInterval {
            latestUpdateDate := localDateTime()
            isLatestOutpatientQueueListUpdated := getLatestOutpatientQueueList(1, mcr)
            if isLatestOutpatientQueueListUpdated {
                outpatientQueueListSummary, err := outpatientQueueListService.GetOutpatientQueueSummaryByMCR(mcr)
                if err != nil {
                    return middleware.NoContent(c)
                }

                if len(outpatientQueueListSummary) > 0 {
                    infoMap := fiber.Map{}
                    regCount = 0
                    apptCount = 0
                    seenCount = 0
                    kivCount = 0
                    for _, info := range outpatientQueueListSummary {
                        queueCriteria := info["queueCriteria"].(string)
                        queueCount := info["queueCount"].(int)
                        if strings.EqualFold(queueCriteria, queueCriteriaReg) {
                            regCount = queueCount
                        } else if strings.EqualFold(queueCriteria, queueCriteriaAppt) {
                            apptCount = queueCount
                        } else if strings.EqualFold(queueCriteria, queueCriteriaSeen) {
                            seenCount = queueCount
                        } else if strings.EqualFold(queueCriteria, queueCriteriaKiv) {
                            kivCount = queueCount
                        }
                    }

                    outpatientQueueSummaryDtoReg := model.OutpatientQueueSummary{
                        QueueCriteria:      queueCriteriaReg,
                        QueueCount:         regCount,
                        ImageName:          "assets/images/Regd.png",
                        LinkName:           "regd",
                        LastUpdateDateTime: latestUpdateDate,
                    }
                    infoMap[queueCriteriaReg] = outpatientQueueSummaryDtoReg

                    outpatientQueueSummaryDtoAppt := model.OutpatientQueueSummary{
                        QueueCriteria:      queueCriteriaAppt,
                        QueueCount:         apptCount,
                        ImageName:          "assets/images/Appt.png",
                        LinkName:           "appt",
                        LastUpdateDateTime: latestUpdateDate,
                    }
                    infoMap[queueCriteriaAppt] = outpatientQueueSummaryDtoAppt

                    outpatientQueueSummaryDtoSeen := model.OutpatientQueueSummary{
                        QueueCriteria:      queueCriteriaSeen,
                        QueueCount:         seenCount,
                        ImageName:          "assets/images/Seen.png",
                        LinkName:           "seen",
                        LastUpdateDateTime: latestUpdateDate,
                    }
                    infoMap[queueCriteriaSeen] = outpatientQueueSummaryDtoSeen

                    outpatientQueueSummaryDtoKiv := model.OutpatientQueueSummary{
                        QueueCriteria:      queueCriteriaKiv,
                        QueueCount:         kivCount,
                        ImageName:          "assets/images/KIV.png",
                        LinkName:           "kiv",
                        LastUpdateDateTime: latestUpdateDate,
                    }
                    infoMap[queueCriteriaKiv] = outpatientQueueSummaryDtoKiv

                    outpatientQueueSummaryDtoList := []model.OutpatientQueueSummary{
                        outpatientQueueSummaryDtoReg,
                        outpatientQueueSummaryDtoAppt,
                        outpatientQueueSummaryDtoSeen,
                        outpatientQueueSummaryDtoKiv,
                    }

                    return c.JSON(outpatientQueueSummaryDtoList)
                }
            }
        } else {
            outpatientQueueListSummary, err := outpatientQueueListService.GetOutpatientQueueSummaryByMCR(mcr)
            if err != nil {
                return middleware.NoContent(c)
            }

            if len(outpatientQueueListSummary) > 0 {
                infoMap := fiber.Map{}
                regCount = 0
                apptCount = 0
                seenCount = 0
                kivCount = 0
                for _, info := range outpatientQueueListSummary {
                    queueCriteria := info["queueCriteria"].(string)
                    queueCount := info["queueCount"].(int)
                    if strings.EqualFold(queueCriteria, queueCriteriaReg) {
                        regCount = queueCount
                    } else if strings.EqualFold(queueCriteria, queueCriteriaAppt) {
                        apptCount = queueCount
                    } else if strings.EqualFold(queueCriteria, queueCriteriaSeen) {
                        seenCount = queueCount
                    } else if strings.EqualFold(queueCriteria, queueCriteriaKiv) {
                        kivCount = queueCount
                    }
                }

                outpatientQueueSummaryDtoReg := model.OutpatientQueueSummary{
                    QueueCriteria:      queueCriteriaReg,
                    QueueCount:         regCount,
                    ImageName:          "assets/images/Regd.png",
                    LinkName:           "regd",
                    LastUpdateDateTime: lastUpdateDateFromTable,
                }
                infoMap[queueCriteriaReg] = outpatientQueueSummaryDtoReg

                outpatientQueueSummaryDtoAppt := model.OutpatientQueueSummary{
                    QueueCriteria:      queueCriteriaAppt,
                    QueueCount:         apptCount,
                    ImageName:          "assets/images/Appt.png",
                    LinkName:           "appt",
                    LastUpdateDateTime: lastUpdateDateFromTable,
                }
                infoMap[queueCriteriaAppt] = outpatientQueueSummaryDtoAppt

                outpatientQueueSummaryDtoSeen := model.OutpatientQueueSummary{
                    QueueCriteria:      queueCriteriaSeen,
                    QueueCount:         seenCount,
                    ImageName:          "assets/images/Seen.png",
                    LinkName:           "seen",
                    LastUpdateDateTime: lastUpdateDateFromTable,
                }
                infoMap[queueCriteriaSeen] = outpatientQueueSummaryDtoSeen

                outpatientQueueSummaryDtoKiv := model.OutpatientQueueSummary{
                    QueueCriteria:      queueCriteriaKiv,
                    QueueCount:         kivCount,
                    ImageName:          "assets/images/KIV.png",
                    LinkName:           "kiv",
                    LastUpdateDateTime: lastUpdateDateFromTable,
                }
                infoMap[queueCriteriaKiv] = outpatientQueueSummaryDtoKiv

                outpatientQueueSummaryDtoList := []model.OutpatientQueueSummary{
                    outpatientQueueSummaryDtoReg,
                    outpatientQueueSummaryDtoAppt,
                    outpatientQueueSummaryDtoSeen,
                    outpatientQueueSummaryDtoKiv,
                }

                return c.JSON(outpatientQueueSummaryDtoList)
            }
        }
    } else {
        isLatestOutpatientQueueListUpdated := getLatestOutpatientQueueList(1, mcr)
        if isLatestOutpatientQueueListUpdated {
            latestUpdateDate := localDateTime()
            outpatientQueueListSummary, err := outpatientQueueListService.GetOutpatientQueueSummaryByMCR(mcr)
            if err != nil {
                return middleware.NoContent(c)
            }

            if len(outpatientQueueListSummary) > 0 {
                infoMap := fiber.Map{}
                regCount = 0
                apptCount = 0
                seenCount = 0
                kivCount = 0
                for _, info := range outpatientQueueListSummary {
                    queueCriteria := info["queueCriteria"].(string)
                    queueCount := info["queueCount"].(int)
                    if strings.EqualFold(queueCriteria, queueCriteriaReg) {
                        regCount = queueCount
                    } else if strings.EqualFold(queueCriteria, queueCriteriaAppt) {
                        apptCount = queueCount
                    } else if strings.EqualFold(queueCriteria, queueCriteriaSeen) {
                        seenCount = queueCount
                    } else if strings.EqualFold(queueCriteria, queueCriteriaKiv) {
                        kivCount = queueCount
                    }
                }

                outpatientQueueSummaryDtoReg := model.OutpatientQueueSummary{
                    QueueCriteria:      queueCriteriaReg,
                    QueueCount:         regCount,
                    ImageName:          "assets/images/Regd.png",
                    LinkName:           "regd",
                    LastUpdateDateTime: latestUpdateDate,
                }
                infoMap[queueCriteriaReg] = outpatientQueueSummaryDtoReg

                outpatientQueueSummaryDtoAppt := model.OutpatientQueueSummary{
                    QueueCriteria:      queueCriteriaAppt,
                    QueueCount:         apptCount,
                    ImageName:          "assets/images/Appt.png",
                    LinkName:           "appt",
                    LastUpdateDateTime: latestUpdateDate,
                }
                infoMap[queueCriteriaAppt] = outpatientQueueSummaryDtoAppt

                outpatientQueueSummaryDtoSeen := model.OutpatientQueueSummary{
                    QueueCriteria:      queueCriteriaSeen,
                    QueueCount:         seenCount,
                    ImageName:          "assets/images/Seen.png",
                    LinkName:           "seen",
                    LastUpdateDateTime: latestUpdateDate,
                }
                infoMap[queueCriteriaSeen] = outpatientQueueSummaryDtoSeen

                outpatientQueueSummaryDtoKiv := model.OutpatientQueueSummary{
                    QueueCriteria:      queueCriteriaKiv,
                    QueueCount:         kivCount,
                    ImageName:          "assets/images/KIV.png",
                    LinkName:           "kiv",
                    LastUpdateDateTime: latestUpdateDate,
                }
                infoMap[queueCriteriaKiv] = outpatientQueueSummaryDtoKiv

                outpatientQueueSummaryDtoList := []model.OutpatientQueueSummary{
                    outpatientQueueSummaryDtoReg,
                    outpatientQueueSummaryDtoAppt,
                    outpatientQueueSummaryDtoSeen,
                    outpatientQueueSummaryDtoKiv,
                }

                return c.JSON(outpatientQueueSummaryDtoList)
            }
        }
    }

    return middleware.NoContent(c)
}

// GetOutpatientQueueDetailList
//
// @Tags Vesalius
// @Produce json
// @Param        branchId        path      string  false  "BranchId"
// @Param        queueCriteria   path      string  true   "QueueCriteria"
// @Param        mcr             path      string  true   "MCR"
// @Security BearerAuth
// @Success 200 {array} model.OutpatientQueueList
// @Router /vesalius/getOutpatientQueueDetailList/{branchId}/{queueCriteria}/{mcr} [get]
func GetOutpatientQueueDetailList(c fiber.Ctx) error {
    queueCriteria := c.Params("queueCriteria")
    mcr := c.Params("mcr")
    _, _, err := middleware.ValidateToken(c)
    if err != nil {
        return middleware.NoContent(c)
    }

    lastUpdateDate, err := outpatientQueueListService.GetLastUpdateDateByMCR(mcr)
    if err != nil {
        return middleware.NoContent(c)
    }

    if lastUpdateDate != "" {
        localLastUpdateDate := stringToGoment(lastUpdateDate)
        minutesDiff := getDateTimeDifferent(localLastUpdateDate)
        vesaliusWSInterval := config.VesaliusWSInterval()
        if minutesDiff > vesaliusWSInterval {
            isLatestOutpatientQueueListUpdated := getLatestOutpatientQueueList(1, mcr)
            if isLatestOutpatientQueueListUpdated {
                outpatientQueueListDetail, err := outpatientQueueListService.GetOutpatientQueueDetailListByMcr(mcr, queueCriteria)
                if err != nil {
                    return middleware.NoContent(c)
                }

                if len(outpatientQueueListDetail) > 0 {
                    return c.JSON(outpatientQueueListDetail)
                }
            }

        } else {
            outpatientQueueListDetail, err := outpatientQueueListService.GetOutpatientQueueDetailListByMcr(mcr, queueCriteria)
            if err != nil {
                return middleware.NoContent(c)
            }

            if len(outpatientQueueListDetail) > 0 {
                return c.JSON(outpatientQueueListDetail)
            }
        }
    } else {
        isLatestOutpatientQueueListUpdated := getLatestOutpatientQueueList(1, mcr)
        if isLatestOutpatientQueueListUpdated {
            outpatientQueueListDetail, err := outpatientQueueListService.GetOutpatientQueueDetailListByMcr(mcr, queueCriteria)
            if err != nil {
                return middleware.NoContent(c)
            }

            if len(outpatientQueueListDetail) > 0 {
                return c.JSON(outpatientQueueListDetail)
            }
        }
    }

    return middleware.NoContent(c)
}

// GetInpatientQueueDetailList
//
// @Tags Vesalius
// @Produce json
// @Param        branchId        path      string  false  "BranchId"
// @Param        mcr             path      string  true   "MCR"
// @Security BearerAuth
// @Success 200 {array} model.InpatientQueueList
// @Router /vesalius/getInpatientQueueDetailList/{branchId}/{mcr} [get]
func GetInpatientQueueDetailList(c fiber.Ctx) error {
    mcr := c.Params("mcr")
    _, _, err := middleware.ValidateToken(c)
    if err != nil {
        return middleware.NoContent(c)
    }

    lastUpdateDate, err := inpatientQueueListService.GetLastUpdateDateByMCR(mcr)
    if err != nil {
        return middleware.NoContent(c)
    }

    if lastUpdateDate != "" {
        localLastUpdateDate := stringToGoment(lastUpdateDate)
        minutesDiff := getDateTimeDifferent(localLastUpdateDate)
        vesaliusWSInterval := config.VesaliusWSInterval()
        if minutesDiff > vesaliusWSInterval {
            inpatientQueueListService.DeleteInpatientQueueListByMcr(mcr)
            isLatestInpatientQueueListUpdated := getLatestInpatientQueueList(1, mcr)
            if isLatestInpatientQueueListUpdated {
                inpatientQueueListSummary, err := inpatientQueueListService.GetInpatientQueueListByMcr(mcr)
                if err != nil {
                    return middleware.NoContent(c)
                }

                if len(inpatientQueueListSummary) > 0 {
                    return c.JSON(inpatientQueueListSummary)
                }
            }
        }
    } else {
        isLatestInpatientQueueListUpdated := getLatestInpatientQueueList(1, mcr)
        if isLatestInpatientQueueListUpdated {
            inpatientQueueListSummary, err := inpatientQueueListService.GetInpatientQueueListByMcr(mcr)
            if err != nil {
                return middleware.NoContent(c)
            }

            if len(inpatientQueueListSummary) > 0 {
                return c.JSON(inpatientQueueListSummary)
            }
        }
    }

    return middleware.NoContent(c)
}

// GetInpatientQueueDetailByPrn
//
// @Tags Vesalius
// @Produce json
// @Param        branchId        path      string  false  "BranchId"
// @Param        mcr             path      string  true   "MCR"
// @Param        prn             path      string  true   "PRN"
// @Security BearerAuth
// @Success 200 {object} model.InpatientQueueList
// @Router /vesalius/getInpatientQueueDetailByPrn/{branchId}/{mcr}/{prn} [get]
func GetInpatientQueueDetailByPrn(c fiber.Ctx) error {
    mcr := c.Params("mcr")
    prn := c.Params("prn")
    _, _, err := middleware.ValidateToken(c)
    if err != nil {
        return middleware.NoContent(c)
    }

    inpatientData, err := inpatientQueueListService.GetInpatientDetailByPrn(mcr, prn)
    if err != nil {
        return middleware.NoContent(c)
    }

    if inpatientData != nil {
        return c.JSON(inpatientData)
    }

    return middleware.NoContent(c)
}

// GetDoctorToDoNotification
//
// @Tags Vesalius
// @Produce json
// @Param        branchId        path      string  false  "BranchId"
// @Param        mcr             path      string  true   "MCR"
// @Security BearerAuth
// @Success 200 {array} model.DoctorToDoNotification
// @Router /vesalius/get-doctor-todo-notification/{branchId}/{mcr} [get]
func GetDoctorToDoNotification(c fiber.Ctx) error {
    mcr := c.Params("mcr")
    _, _, err := middleware.ValidateToken(c)
    if err != nil {
        return middleware.NoContent(c)
    }

    lastUpdateDate, err := doctorTodoNotificationService.GetLastUpdateDateByMcr(mcr)
    if err != nil {
        return middleware.NoContent(c)
    }

    if lastUpdateDate != "" {
        localLastUpdateDate := stringToGoment(lastUpdateDate)
        minutesDiff := getDateTimeDifferent(localLastUpdateDate)
        vesaliusWSInterval := config.VesaliusWSInterval()
        if minutesDiff > vesaliusWSInterval {
            doctorTodoNotificationService.DeleteDoctorToDoNotificationListByMcr(mcr)
            isLatestDoctorToDoNotificationListUpdated := getDoctorToDoNotificationList(1, mcr)
            if isLatestDoctorToDoNotificationListUpdated {
                doctorToDoNotificationList, err := doctorTodoNotificationService.GetDoctorToDoNotificationListByMcr(mcr)
                if err != nil {
                    return middleware.NoContent(c)
                }

                if len(doctorToDoNotificationList) > 0 {
                    return c.JSON(doctorToDoNotificationList)
                }
            }
        } else {
            doctorToDoNotificationList, err := doctorTodoNotificationService.GetDoctorToDoNotificationListByMcr(mcr)
            if err != nil {
                return middleware.NoContent(c)
            }

            if len(doctorToDoNotificationList) > 0 {
                return c.JSON(doctorToDoNotificationList)
            }
        }
    } else {
        isLatestDoctorToDoNotificationListUpdated := getDoctorToDoNotificationList(1, mcr)
        if isLatestDoctorToDoNotificationListUpdated {
            doctorToDoNotificationList, err := doctorTodoNotificationService.GetDoctorToDoNotificationListByMcr(mcr)
            if err != nil {
                return middleware.NoContent(c)
            }

            if len(doctorToDoNotificationList) > 0 {
                return c.JSON(doctorToDoNotificationList)
            }
        }
    }

    return middleware.NoContent(c)
}

// GetDoctorToDoNotificationDetails
//
// @Tags Vesalius
// @Produce json
// @Param        branchId        path      string  false  "BranchId"
// @Param        mcr             path      string  true   "MCR"
// @Param        prn             path      string  true   "PRN"
// @Security BearerAuth
// @Success 200 {array} model.DoctorToDoNotification
// @Router /vesalius/get-doctor-todo-notification-details/{branchId}/{mcr}/{prn} [get]
func GetDoctorToDoNotificationDetails(c fiber.Ctx) error {
    mcr := c.Params("mcr")
    prn := c.Params("prn")
    _, _, err := middleware.ValidateToken(c)
    if err != nil {
        return middleware.NoContent(c)
    }

    lastUpdateDate, err := doctorTodoNotificationService.GetLastUpdateDateByMcr(mcr)
    if err != nil {
        return middleware.NoContent(c)
    }

    if lastUpdateDate != "" {
        localLastUpdateDate := stringToGoment(lastUpdateDate)
        minutesDiff := getDateTimeDifferent(localLastUpdateDate)
        vesaliusWSInterval := config.VesaliusWSInterval()
        if minutesDiff > vesaliusWSInterval {
            doctorTodoNotificationService.DeleteDoctorToDoNotificationListByMcr(mcr)
            isLatestDoctorToDoNotificationListUpdated := getDoctorToDoNotificationList(1, mcr)
            if isLatestDoctorToDoNotificationListUpdated {
                doctorToDoNotificationList, err := doctorTodoNotificationService.GetDoctorToDoNotificationListByMcrPrnNotificationType(mcr, prn)
                if err != nil {
                    return middleware.NoContent(c)
                }

                if len(doctorToDoNotificationList) > 0 {
                    return c.JSON(doctorToDoNotificationList)
                }
            }
        } else {
            doctorToDoNotificationList, err := doctorTodoNotificationService.GetDoctorToDoNotificationListByMcrPrnNotificationType(mcr, prn)
            if err != nil {
                return middleware.NoContent(c)
            }

            if len(doctorToDoNotificationList) > 0 {
                return c.JSON(doctorToDoNotificationList)
            }
        }
    } else {
        isLatestDoctorToDoNotificationListUpdated := getDoctorToDoNotificationList(1, mcr)
        if isLatestDoctorToDoNotificationListUpdated {
            doctorToDoNotificationList, err := doctorTodoNotificationService.GetDoctorToDoNotificationListByMcr(mcr)
            if err != nil {
                return middleware.NoContent(c)
            }

            if len(doctorToDoNotificationList) > 0 {
                return c.JSON(doctorToDoNotificationList)
            }
        }
    }

    return middleware.NoContent(c)
}

// ProcessDoctorToDoAcknowledge
//
// @Tags Vesalius
// @Produce json
// @Param request body dto.PostProcessDoctorToDoAcknowledgeDto true "ProcessDoctorToDoAcknowledge Request"
// @Param        branchId        path      string  false  "BranchId"
// @Param        mcr             path      string  true   "MCR"
// @Security BearerAuth
// @Success 200
// @Router /vesalius/process-doctor-todo-ack/{branchId}/{mcr} [post]
func ProcessDoctorToDoAcknowledge(c fiber.Ctx) error {
    mcr := c.Params("mcr")
    _, _, err := middleware.ValidateToken(c)
    if err != nil {
        return middleware.NoContent(c)
    }

    data := new(dto.PostProcessDoctorToDoAcknowledgeDto)
    if err := c.Bind().Body(data); err != nil {
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            errs := utils.GetValidationErrors(validationErrors)
            if errs != nil {
                return middleware.Unauthorized(c)
            }
        }

        return middleware.Unauthorized(c)
    }

    doctorTodoNotificationService.UpdateDoctorToDoNotificationAcknowledgementFlag(mcr, data.AccessionNo)
    doctorReq := model.DoctorRequest{
        Doctor_request_type: "ACK",
        Accession_no:        data.AccessionNo,
        Notification_type:   data.NotificationType,
        Remark:              data.Remark,
        Mcr:                 mcr,
        Posted_flag:         false,
        Create_date_time:    localDateTime(),
    }
    doctorRequestService.AddDoctorRequest(doctorReq)
    return c.JSON(fiber.Map{
        "success": 1,
    })
}

// ProcessDoctorReviewInvestigation
//
// @Tags Vesalius
// @Produce json
// @Param request body dto.PostProcessDoctorReviewInvestigationDto true "ProcessDoctorReviewInvestigation Request"
// @Param        branchId        path      string  false  "BranchId"
// @Param        mcr             path      string  true   "MCR"
// @Security BearerAuth
// @Success 200
// @Router /vesalius/process-doctor-review-investigation/{branchId}/{mcr} [post]
func ProcessDoctorReviewInvestigation(c fiber.Ctx) error {
    mcr := c.Params("mcr")
    _, _, err := middleware.ValidateToken(c)
    if err != nil {
        return middleware.NoContent(c)
    }

    data := new(dto.PostProcessDoctorReviewInvestigationDto)
    if err := c.Bind().Body(data); err != nil {
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            errs := utils.GetValidationErrors(validationErrors)
            if errs != nil {
                return middleware.Unauthorized(c)
            }
        }

        return middleware.Unauthorized(c)
    }

    investigationReportService.UpdateReviewInvestigationFlag(data.AccessionNo)
    doctorReq := model.DoctorRequest{
        Doctor_request_type: "REV",
        Accession_no:        data.AccessionNo,
        Prn:                 data.PRN,
        Review_date:         data.ReviewDate,
        Review_time:         data.ReviewTime,
        Review_doctor:       data.ReviewDoctor,
        Mcr:                 mcr,
        Posted_flag:         false,
        Create_date_time:    localDateTime(),
    }
    doctorRequestService.AddDoctorRequest(doctorReq)
    return c.JSON(fiber.Map{
        "success": 1,
    })
}

// GetInvestigationReport
//
// @Tags Vesalius
// @Produce json
// @Param        branchId        path      string  false  "BranchId"
// @Param        mcr             path      string  true   "MCR"
// @Param        orderDate       path      string  true   "OrderDate"
// @Security BearerAuth
// @Success 200 {array} model.InvestigationReport
// @Router /vesalius/get-investigation-report/{branchId}/{mcr}/{orderDate} [get]
func GetInvestigationReport(c fiber.Ctx) error {
    mcr := c.Params("mcr")
    orderDate := c.Params("orderDate")
    _, _, err := middleware.ValidateToken(c)
    if err != nil {
        return middleware.NoContent(c)
    }

    lastUpdateDate, err := investigationReportService.GetLastUpdateDateByOrderDate(orderDate)
    if err != nil {
        return middleware.NoContent(c)
    }

    if lastUpdateDate != "" {
        localLastUpdateDate := stringToGoment(lastUpdateDate)
        minutesDiff := getDateTimeDifferent(localLastUpdateDate)
        vesaliusWSInterval := config.VesaliusWSInterval()
        if minutesDiff > vesaliusWSInterval {
            investigationReportService.DeleteInvestigationListByOrderDate(orderDate)
            isLatestInvestigationListUpdated := getLatestInvestigationList(1, orderDate)
            if isLatestInvestigationListUpdated {
                investigationList, err := investigationReportService.GetInvestigationListByOrderDate(mcr, orderDate)
                if err != nil {
                    return middleware.NoContent(c)
                }

                if len(investigationList) > 0 {
                    for i, inv := range investigationList {
                        r := strings.ToLower(inv.Report_type)
                        if r == "pdf" {
                            investigationList[i].Result = nil
                        }
                    }

                    return c.JSON(investigationList)
                }
            }
        } else {
            investigationList, err := investigationReportService.GetInvestigationListByOrderDate(mcr, orderDate)
            if err != nil {
                return middleware.NoContent(c)
            }

            if len(investigationList) > 0 {
                for i, inv := range investigationList {
                    r := strings.ToLower(inv.Report_type)
                    if r == "pdf" {
                        investigationList[i].Result = nil
                    }
                }

                return c.JSON(investigationList)
            }
        }
    } else {
        isLatestInvestigationListUpdated := getLatestInvestigationList(1, orderDate)
        if isLatestInvestigationListUpdated {
            investigationList, err := investigationReportService.GetInvestigationListByOrderDate(mcr, orderDate)
            if err != nil {
                return middleware.NoContent(c)
            }

            if len(investigationList) > 0 {
                for i, inv := range investigationList {
                    r := strings.ToLower(inv.Report_type)
                    if r == "pdf" {
                        investigationList[i].Result = nil
                    }
                }

                return c.JSON(investigationList)
            }
        }
    }

    return middleware.NoContent(c)
}

// GetPdfInvestigationReport
//
// @Tags Vesalius
// @Produce octet-stream
// @Param        branchId        path      string  false  "BranchId"
// @Param        mcr             path      string  true   "MCR"
// @Param        accessionNo     path      string  true   "AccessionNo"
// @Security BearerAuth
// @Success 200 {array} model.InvestigationReport
// @Router /vesalius/get-pdf-investigation-report/{branchId}/{mcr}/{accessionNo} [get]
func GetPdfInvestigationReport(c fiber.Ctx) error {
    mcr := c.Params("mcr")
    accessionNo := c.Params("accessionNo")
    _, _, err := middleware.ValidateToken(c)
    if err != nil {
        return middleware.NoContent(c)
    }

    pdfResult, err := investigationReportService.GetInvestigationPDFResultByAccessionNo(mcr, accessionNo)
    if err != nil {
        return middleware.NoContent(c)
    }

    if pdfResult != "" {
        data, err := base64.StdEncoding.DecodeString(pdfResult)
        if err != nil {
            return middleware.NoContent(c)
        }

        return c.Send(data)
    }

    return middleware.NoContent(c)
}

func getLatestInpatientQueueList(_ int, mcr string) bool {
    inpatientQueueList, _ := ws.GetInpatientQueueListByMCR(mcr)
    if inpatientQueueList != nil {
        lastUpdateDate := localDateTime()
        o := inpatientQueueList[0]
        for _, inpatientQueue := range o.Patients {
            patientQueueList := model.InpatientQueueList{
                Patient_type:         "I",
                Mcr:                  mcr,
                Prn:                  inpatientQueue.PRN,
                Title:                inpatientQueue.Name.Title,
                First_name:           inpatientQueue.Name.FirstName,
                Middle_name:          inpatientQueue.Name.MiddleName,
                Last_name:            inpatientQueue.Name.LastName,
                Sex_code:             inpatientQueue.Sex.Code,
                Sex_desc:             inpatientQueue.Sex.Description,
                Age:                  inpatientQueue.Age,
                Nationality:          inpatientQueue.Nationality,
                Vip_flag:             inpatientQueue.VIPFlag,
                Visit_type:           strings.ReplaceAll(inpatientQueue.VisitType, "&amp;", "&"),
                Visit_number:         inpatientQueue.VisitNumber,
                Ward:                 inpatientQueue.Ward,
                Bed:                  inpatientQueue.Bed,
                Admission_date:       inpatientQueue.AdmissionDate,
                Admission_time:       inpatientQueue.AdmissionTime,
                Triage_score:         inpatientQueue.TriageScore,
                Triage_discriminator: inpatientQueue.TriageDiscriminator,
                Last_update_date:     lastUpdateDate,
            }
            inpatientQueueListService.AddInpatientQueueList(patientQueueList)
        }

        return true
    }

    return false
}

func getLatestOutpatientQueueList(_ int, mcr string) bool {
    outpatientQueueList, _ := ws.GetOutpatientQueueListByMCR(mcr)
    if outpatientQueueList != nil {
        outpatientQueueListService.DeleteOutpatientQueueListByMcr(mcr)
        lastUpdateDate := localDateTime()
        o := outpatientQueueList[0]
        for _, outpatientQueue := range o.Patients {
            patientQueueList := model.OutpatientQueueList{
                Patient_type:          "O",
                Mcr:                   mcr,
                Prn:                   outpatientQueue.PRN,
                Title:                 outpatientQueue.Name.Title,
                First_name:            outpatientQueue.Name.FirstName,
                Middle_name:           outpatientQueue.Name.MiddleName,
                Last_name:             outpatientQueue.Name.LastName,
                Sex_code:              outpatientQueue.Sex.Code,
                Sex_desc:              outpatientQueue.Sex.Description,
                Age:                   outpatientQueue.Age,
                Nationality:           outpatientQueue.Nationality,
                Vip_flag:              outpatientQueue.VIPFlag,
                Visit_type:            strings.ReplaceAll(outpatientQueue.VisitType, "&amp;", "&"),
                Visit_number:          outpatientQueue.VisitNumber,
                Queue_number:          outpatientQueue.QueueNumber,
                Queue_criteria:        strings.ReplaceAll(outpatientQueue.QueueCriteria, "Arrived / Inprogress", "Registered"),
                Patient_status:        outpatientQueue.PatientStatus,
                Registration_date:     outpatientQueue.RegistrationDate,
                Registration_time:     outpatientQueue.RegistrationTime,
                Appointment_date:      outpatientQueue.AppointmentDate,
                Appointment_time:      outpatientQueue.AppointmentTime,
                Vital_are_available:   outpatientQueue.VitalsAreAvailable,
                Triage_score:          outpatientQueue.TriageScore,
                Triage_discriminator:  outpatientQueue.TriageDiscriminator,
                Has_on_arrival_orders: outpatientQueue.HasOnArrivalOrders,
                Routed_by:             outpatientQueue.RoutedBy,
                Last_update_date:      lastUpdateDate,
            }
            outpatientQueueListService.AddOutpatientQueueList(patientQueueList)
        }

        return true
    }

    return false
}

func getLatestInvestigationList(_ int, orderDate string) bool {
    investigationReportResult, _ := ws.GetInvestigationReportList(orderDate)
    if investigationReportResult != nil {
        if len(investigationReportResult) < 1{
            return false
        }

        lastUpdateDate := localDateTime()
        for _, visit := range investigationReportResult[0].Visits {
            for _, serviceComponent := range visit.Reports[0].ServiceComponents {
                investigationReport := model.InvestigationReport{
                    Last_update_date:          lastUpdateDate,
                    Account_no:                visit.AccountNo,
                    Prn:                       visit.PRN,
                    Patient_name:              visit.PatientName,
                    Accession_no:              serviceComponent.AccessionNo,
                    Investigation_type:        serviceComponent.Type,
                    Service_code:              serviceComponent.Code,
                    Service_desc:              serviceComponent.Description,
                    Order_doctor_mcr:          serviceComponent.OrderDoctorMcr,
                    Order_date:                serviceComponent.OrderDate,
                    Result_date:               serviceComponent.ResultDate,
                    Report_type:               serviceComponent.ReportType,
                    Result:                    &serviceComponent.Result,
                    Remark:                    serviceComponent.Remark,
                    Review_investigation_flag: false,
                }
                investigationReportService.AddInvestigationReport(investigationReport)
            }
        }

        return true
    }

    return false
}

func getDoctorToDoNotificationList(_ int, mcr string) bool {
    doctorToDoNotificationResult, _ := ws.GetDoctorToDoNotificationList(mcr)
    if doctorToDoNotificationResult != nil {
        lastUpdateDate := localDateTime()
        for _, todoNotification := range doctorToDoNotificationResult {
            if todoNotification.DrugDiscontinues != nil {
                for _, drugItem := range todoNotification.DrugDiscontinues[0].Items {
                    doctorToDoNotification := model.DoctorToDoNotification{
                        Mcr:                  mcr,
                        Last_update_date:     lastUpdateDate,
                        Notification_type:    "Drug Discontinue",
                        Prn:                  todoNotification.Prn,
                        Patient_name:         todoNotification.PatientName,
                        Account_no:           todoNotification.AccountNo,
                        Doctor_name:          todoNotification.DoctorName,
                        Accession_no:         drugItem.AccessionNo,
                        Item_code:            drugItem.Code,
                        Item_desc:            drugItem.Description,
                        Item_quantity:        drugItem.Quantity,
                        UOM:                  drugItem.UOM,
                        Instruction:          drugItem.Instruction,
                        Order_date:           drugItem.OrderDate,
                        Order_time:           drugItem.OrderTime,
                        Order_by:             drugItem.OrderedBy,
                        Discontinue_date:     drugItem.DiscontinueDate,
                        Discontinue_time:     drugItem.DiscontinueTime,
                        Discontinue_reason:   drugItem.DiscontinueReason,
                        Discontinue_by:       drugItem.DiscontinueBy,
                        Acknowledgement_flag: false,
                    }
                    doctorTodoNotificationService.AddDoctorToDoNotification(doctorToDoNotification)
                }
            }

            if todoNotification.VerbalOrders != nil {
                for _, verbalOrderItem := range todoNotification.VerbalOrders[0].Items {
                    doctorToDoNotification := model.DoctorToDoNotification{
                        Mcr:                  mcr,
                        Last_update_date:     lastUpdateDate,
                        Notification_type:    "Verbal Order",
                        Prn:                  todoNotification.Prn,
                        Patient_name:         todoNotification.PatientName,
                        Account_no:           todoNotification.AccountNo,
                        Doctor_name:          todoNotification.DoctorName,
                        Accession_no:         verbalOrderItem.AccessionNo,
                        Item_code:            verbalOrderItem.Code,
                        Item_desc:            verbalOrderItem.Description,
                        Item_quantity:        verbalOrderItem.Quantity,
                        UOM:                  verbalOrderItem.UOM,
                        Instruction:          verbalOrderItem.Instruction,
                        Order_date:           verbalOrderItem.OrderDate,
                        Order_time:           verbalOrderItem.OrderTime,
                        Order_by:             verbalOrderItem.OrderedBy,
                        Dispense_quantity:    verbalOrderItem.DispenseQty,
                        Last_dispense_time:   verbalOrderItem.LastDispenseTime,
                        Served_quantity:      verbalOrderItem.ServedQty,
                        Acknowledgement_flag: false,
                    }
                    doctorTodoNotificationService.AddDoctorToDoNotification(doctorToDoNotification)
                }
            }

            if todoNotification.DischargeSummary != "" {
                doctorToDoNotification := model.DoctorToDoNotification{
                    Mcr:                  mcr,
                    Last_update_date:     lastUpdateDate,
                    Notification_type:    "Discharge Summary",
                    Prn:                  todoNotification.Prn,
                    Patient_name:         todoNotification.PatientName,
                    Account_no:           todoNotification.AccountNo,
                    Doctor_name:          todoNotification.DoctorName,
                    Discharge_summary:    todoNotification.DischargeSummary,
                    Acknowledgement_flag: false,
                }
                doctorTodoNotificationService.AddDoctorToDoNotification(doctorToDoNotification)
            }
        }

        return true
    }

    return false
}

func getDateTimeDifferent(lastUpdateDateTime *goment.Goment) int {
    gdt := stringToGoment(localDateTime())
    x := gdt.Millisecond() - lastUpdateDateTime.Millisecond()
    diff := math.Abs(float64(x))
    a := math.Trunc(diff / (60 * 1000))
    diffmin := int(a)
    return diffmin
}

func stringToGoment(s string) *goment.Goment {
    g, _ := goment.New(s, "DD/MM/YYYY hh:mm A")
    return g
}

func localDateTime() string {
    g, _ := goment.New()
    return g.Local().Format("DD/MM/YYYY hh:mm A")
}