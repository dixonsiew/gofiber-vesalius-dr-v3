package ws

import (
	"strings"
	"vesaliusdr/config"
	"vesaliusdr/xmlmodel"
)

var patientDocumentCode string = config.Config("patient.document.code")
var myKadWSWithDash string = config.Config("mykad.ws.withDash")
// var nricWSWithDash string = config.Config("ws.vesalius.nric.withDash")

func GetPatientData(prn string) (*xmlmodel.Patient, *xmlmodel.VesaliusWSException) {
    localPrn := ""
    if strings.EqualFold(myKadWSWithDash, "Y") {
        localPrn = prn
    } else {
        localPrn = strings.ReplaceAll(prn, "-", "")
    }

    result, ex := PatientGetPatientData(localPrn)
    if ex != nil {
        return nil, ex
    }

    if result.Error.ErrorCode != "" {
        return nil, &xmlmodel.VesaliusWSException{
            Code:    result.Error.ErrorCode,
            Message: result.Error.ErrorMessage,
        }
    }

    if result.Patients == nil {
        return nil, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    if len(result.Patients) > 1 {
        return nil, &xmlmodel.VesaliusWSException{
            Code:    "99",
            Message: "More than 1 Patient Found",
        }
    }

    patientData := result.Patients[0]
    if patientData.Prn == localPrn {
        return &patientData, nil
    }

    if patientData.Document != nil {
        patientIDValue := ""
        for _, patientDoc := range patientData.Document {
            if strings.EqualFold(strings.TrimSpace(patientDoc.Code), patientDocumentCode) {
                patientIDValue = patientDoc.Value
            }
        }

        if patientIDValue != "" {
            if strings.EqualFold(strings.TrimSpace(patientIDValue), localPrn) {
                return &patientData, nil
            } else {
                return nil, &xmlmodel.VesaliusWSException{
                    Code:    "ERROR",
                    Message: "Patient Doc Not Found",
                }
            }
        } else {
            return nil, &xmlmodel.VesaliusWSException{
                Code:    "ERROR",
                Message: "Patient Doc Not Found",
            }
        }
    }

    return nil, &xmlmodel.VesaliusWSException{
        Code: "ERROR",
        Message: "PPatient PRN Not Found",
    }
}

func GetInpatientQueueListByMCR(mcr string) ([]xmlmodel.Inpatient, *xmlmodel.VesaliusWSException) {
    localVisitType := "I"
    result, ex := GetInpatientQueueList(mcr, localVisitType)
    if result.Inpatients != nil {
        return result.Inpatients, ex
    }

    return nil, ex
}

func GetOutpatientQueueListByMCR(mcr string) ([]xmlmodel.Outpatient, *xmlmodel.VesaliusWSException) {
    strVisitType := "O"
    result, ex := GetOutpatientQueueList(mcr, strVisitType)
    if result.Outpatients != nil {
        return result.Outpatients, ex
    }

    return nil, ex
}

func GetInvestigationReportList(orderDate string) ([]xmlmodel.Investigation, *xmlmodel.VesaliusWSException) {
    strInvestigationType := "ALL"
    result, ex := GetInvestigationReport(orderDate, strInvestigationType)
    if result.Investigations != nil {
        return result.Investigations, ex
    }

    return nil, ex
}

func GetDoctorToDoNotificationList(mcr string) ([]xmlmodel.DoctorTodoNotification, *xmlmodel.VesaliusWSException) {
    strNotificationType := "ALL"
    result, ex := GetDoctorToDoNotification(mcr, strNotificationType)
    if result.DoctorTodoNotifications != nil {
        return result.DoctorTodoNotifications, ex
    }

    return nil, ex
}