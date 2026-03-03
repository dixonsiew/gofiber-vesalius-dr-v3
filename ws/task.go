package ws

import (
	"fmt"
	"strings"
	srv "vesaliusdr/service/doctor_request"

	"github.com/nleeper/goment"
)

func localDateTime() string {
    g, _ := goment.New()
    return g.Local().Format("DD/MM/YYYY hh:mm A")
}

func ProcessDoctorToReviewData() {
    doctorRequestResult, err := srv.GetDoctorRequestByCreateDateTimeAsc()
    if err != nil {

    }

    if doctorRequestResult != nil {
        if strings.ToUpper(doctorRequestResult.Doctor_request_type) == "REV" {
            result, _ := ProcessDoctorReviewInvestigation(doctorRequestResult.Mcr, doctorRequestResult.Prn, doctorRequestResult.Accession_no, 
                doctorRequestResult.Review_doctor, doctorRequestResult.Review_date, doctorRequestResult.Review_time)
            if result.Success.Code != "" {
                srv.UpdateDoctorRequestStatus(doctorRequestResult.Accession_no, localDateTime(), fmt.Sprintf("%s - %s", result.Success.Code, result.Success.Message))
            } else if result.Error.ErrorCode != "" {
                srv.UpdateDoctorRequestStatus(doctorRequestResult.Accession_no, localDateTime(), fmt.Sprintf("%s - %s", result.Error.ErrorCode, result.Error.ErrorMessage))
            }
        } else if strings.ToUpper(doctorRequestResult.Doctor_request_type) == "ACK" {
            result, _ := ProcessDoctorToDoAcknowledge(doctorRequestResult.Mcr, doctorRequestResult.Notification_type, doctorRequestResult.Accession_no, doctorRequestResult.Remark)
            if result.Success.Code != "" {
                srv.UpdateDoctorRequestStatus(doctorRequestResult.Accession_no, localDateTime(), fmt.Sprintf("%s - %s", result.Success.Code, result.Success.Message))
            } else if result.Error.ErrorCode != "" {
                srv.UpdateDoctorRequestStatus(doctorRequestResult.Accession_no, localDateTime(), fmt.Sprintf("%s - %s", result.Error.ErrorCode, result.Error.ErrorMessage))
            }
        }
    }
}