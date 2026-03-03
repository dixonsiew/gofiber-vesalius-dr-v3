package ws

import (
    "encoding/xml"
    "fmt"
    "time"
    "vesaliusdr/config"
    "vesaliusdr/utils"
    "vesaliusdr/xmlmodel"

    "github.com/gofiber/fiber/v2"
)

var vesaliusServerBaseUrl string = config.Config("ws.vesalius.server_baseurl")
var vesaliusServerCompanyCode string = config.Config("ws.vesalius.server_companycode")
var vesaliusServerSystemCode string = config.Config("ws.vesalius.server_systemcode")
var vesaliusServerPassword string = config.Config("ws.vesalius.server_password")

func AuthenticationLogin() (string, *xmlmodel.ResultToken, *xmlmodel.VesaliusWSException) {
    var token string = ""
    var result *xmlmodel.ResultToken = new(xmlmodel.ResultToken)
    var ex *xmlmodel.VesaliusWSException
    url := fmt.Sprintf("%sAUTHENTICATION/Login.cfc", vesaliusServerBaseUrl)
    r := utils.GetR("login")
    envelope :=
    `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:login>
                <axi:company_code>%s</axi:company_code>
                <axi:system_code>%s</axi:system_code>
                <axi:password>%s</axi:password>
            </axi:login>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, vesaliusServerCompanyCode, vesaliusServerSystemCode, vesaliusServerPassword)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return token, result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return token, result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := utils.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    if err != nil {
        return token, result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    token = result.Token.TokenNumber
    return token, result, ex
}

func Logout(token string) (*xmlmodel.ResultLogout, *xmlmodel.VesaliusWSException) {
    var result *xmlmodel.ResultLogout = new(xmlmodel.ResultLogout)
    var ex *xmlmodel.VesaliusWSException
    url := fmt.Sprintf("%sAUTHENTICATION/Logout.cfc", vesaliusServerBaseUrl)
    r := utils.GetR("Logout")
    envelope :=
    `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:Logout>
                <axi:token_number>%s</axi:token_number>
            </axi:Logout>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, token)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := utils.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    return result, ex
}

func PatientGetPatientData(prn string) (*xmlmodel.ResultListPatient, *xmlmodel.VesaliusWSException) {
    var result *xmlmodel.ResultListPatient = new(xmlmodel.ResultListPatient)
    var ex *xmlmodel.VesaliusWSException
    url := fmt.Sprintf("%sPATIENT/GetPatientData.cfc", vesaliusServerBaseUrl)
    r := utils.GetR("getPatientData")
    envelope :=
    `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:getPatientData>
                <axi:prn>%s</axi:prn>
            </axi:getPatientData>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, prn)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := utils.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    return result, ex
}

func GetOutpatientQueueList(mcr string, visitType string) (*xmlmodel.ResultListOutpatient, *xmlmodel.VesaliusWSException) {
    var result *xmlmodel.ResultListOutpatient = new(xmlmodel.ResultListOutpatient)
    var ex *xmlmodel.VesaliusWSException
    localToken, resx, ex := AuthenticationLogin()
    if ex != nil {
        return result, ex
    }

    if resx.Error.ErrorCode != "" {
        return result, ex
    }

    defer defLogout(localToken)

    sleep()
    url := fmt.Sprintf("%sPATIENT/GetPatientQueueList.cfc", vesaliusServerBaseUrl)
    r := utils.GetR("getPatientQueueList")
    envelope :=
    `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:getPatientQueueList>
                <axi:token_number>%s</axi:token_number>
                <axi:company_code>%s</axi:company_code>
                <axi:mcr>%s</axi:mcr>
                <axi:visit_type>%s</axi:visit_type>
            </axi:getPatientQueueList>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, localToken, vesaliusServerCompanyCode, mcr, visitType)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := utils.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    return result, ex
}

func GetInpatientQueueList(mcr string, visitType string) (*xmlmodel.ResultListInpatient, *xmlmodel.VesaliusWSException) {
    var result *xmlmodel.ResultListInpatient = new(xmlmodel.ResultListInpatient)
    var ex *xmlmodel.VesaliusWSException
    localToken, resx, ex := AuthenticationLogin()
    if ex != nil {
        return result, ex
    }

    if resx.Error.ErrorCode != "" {
        return result, ex
    }

    defer defLogout(localToken)

    sleep()
    url := fmt.Sprintf("%sPATIENT/GetPatientQueueList.cfc", vesaliusServerBaseUrl)
    r := utils.GetR("getPatientQueueList")
    envelope :=
    `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:getPatientQueueList>
                <axi:token_number>%s</axi:token_number>
                <axi:company_code>%s</axi:company_code>
                <axi:mcr>%s</axi:mcr>
                <axi:visit_type>%s</axi:visit_type>
            </axi:getPatientQueueList>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, localToken, vesaliusServerCompanyCode, mcr, visitType)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := utils.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    fmt.Println(result.Inpatients)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    return result, ex
}

func GetInvestigationReport(orderDate string, investigationType string) (*xmlmodel.ResultListInvestigation, *xmlmodel.VesaliusWSException) {
    var result *xmlmodel.ResultListInvestigation = new(xmlmodel.ResultListInvestigation)
    var ex *xmlmodel.VesaliusWSException
    localToken, resx, ex := AuthenticationLogin()
    if ex != nil {
        return result, ex
    }

    if resx.Error.ErrorCode != "" {
        return result, ex
    }

    defer defLogout(localToken)

    sleep()
    url := fmt.Sprintf("%sclinical/GetInvestigationReport.cfc", vesaliusServerBaseUrl)
    r := utils.GetR("GetInvestigationReport")
    envelope :=
    `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:GetInvestigationReport>
                <axi:token_no>%s</axi:token_no>
                <axi:company_code>%s</axi:company_code>
                <axi:investigation_type>%s</axi:investigation_type>
                <axi:account_no>%s</axi:account_no>
                <axi:prn>%s</axi:prn>
                <axi:date>%s</axi:date>
            </axi:GetInvestigationReport>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, localToken, vesaliusServerCompanyCode, investigationType, "", "", orderDate)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := utils.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    return result, ex
}

func GetDoctorToDoNotification(mcr string, notificationType string) (*xmlmodel.ResultListDoctorTodoNotification, *xmlmodel.VesaliusWSException) {
    var result *xmlmodel.ResultListDoctorTodoNotification = new(xmlmodel.ResultListDoctorTodoNotification)
    var ex *xmlmodel.VesaliusWSException
    localToken, resx, ex := AuthenticationLogin()
    if ex != nil {
        return result, ex
    }

    if resx.Error.ErrorCode != "" {
        return result, ex
    }

    defer defLogout(localToken)

    sleep()
    url := fmt.Sprintf("%sclinical/GetDoctorToDoNotification.cfc", vesaliusServerBaseUrl)
    r := utils.GetR("GetDoctorToDoNotification")
    envelope :=
    `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:GetDoctorToDoNotification>
                <axi:company_code>%s</axi:company_code>
                <axi:token_number>%s</axi:token_number>
                <axi:notification>%s</axi:notification>
                <axi:doctor_mcr>%s</axi:doctor_mcr>
            </axi:GetDoctorToDoNotification>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, vesaliusServerCompanyCode, localToken, notificationType, mcr)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := utils.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    return result, ex
}

func ProcessDoctorReviewInvestigation(mcr string, prn string, accessionNo string, reviewDoc string, reviewDate string, reviewTime string) (*xmlmodel.ResultList, *xmlmodel.VesaliusWSException) {
    var result *xmlmodel.ResultList = new(xmlmodel.ResultList)
    var ex *xmlmodel.VesaliusWSException
    localToken, resx, ex := AuthenticationLogin()
    if ex != nil {
        return result, ex
    }

    if resx.Error.ErrorCode != "" {
        return result, ex
    }

    defer defLogout(localToken)

    sleep()
    url := fmt.Sprintf("%sCLINICAL/ProcessReviewInvestigation.cfc", vesaliusServerBaseUrl)
    r := utils.GetR("ProcessReviewInvestigation")
    envelope :=
    `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:ProcessReviewInvestigation>
                <axi:token_no>%s</axi:token_no>
                <axi:company_code>%s</axi:company_code>
                <axi:prn>%s</axi:prn>
                <axi:accession_no>%s</axi:accession_no>
                <axi:review_doc>%s</axi:review_doc>
                <axi:review_date>%s</axi:review_date>
                <axi:review_time>%s</axi:review_time>
            </axi:ProcessReviewInvestigation>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, localToken, vesaliusServerCompanyCode, prn, accessionNo, reviewDoc, reviewDate, reviewTime)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := utils.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    return result, ex
}

func ProcessDoctorToDoAcknowledge(mcr string, notificationType string, accessionNo string, remark string) (*xmlmodel.ResultList, *xmlmodel.VesaliusWSException) {
    var result *xmlmodel.ResultList = new(xmlmodel.ResultList)
    var ex *xmlmodel.VesaliusWSException
    localToken, resx, ex := AuthenticationLogin()
    if ex != nil {
        return result, ex
    }

    if resx.Error.ErrorCode != "" {
        return result, ex
    }

    defer defLogout(localToken)

    sleep()
    url := fmt.Sprintf("%sCLINICAL/ProcessDoctorToDoAcknowledge.cfc", vesaliusServerBaseUrl)
    r := utils.GetR("ProcessDoctorToDoAcknowledge")
    envelope :=
    `
    <x:Envelope
        xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
        xmlns:axi="http://ws.apache.org/axis2">
        <x:Header/>
        <x:Body>
            <axi:ProcessDoctorToDoAcknowledge>
                <axi:company_code>%s</axi:company_code>
                <axi:token_number>%s</axi:token_number>
                <axi:type>%s</axi:type>
                <axi:doctor_mcr>%s</axi:doctor_mcr>
                <axi:accession_no>%s</axi:accession_no>
                <axi:remark>%s</axi:remark>
            </axi:ProcessDoctorToDoAcknowledge>
        </x:Body>
    </x:Envelope>
    `
    v := fmt.Sprintf(envelope, vesaliusServerCompanyCode, localToken, notificationType, mcr, accessionNo, remark)
    res, err := r.SetBody(v).Post(url)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: err.Error(),
        }
    }

    if res.StatusCode() != fiber.StatusOK {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: fmt.Errorf("SOAP Status %d", res.StatusCode()).Error(),
        }
    }

    bx := utils.GetXmlResult(res.String())
    err = xml.Unmarshal(bx, result)
    if err != nil {
        return result, &xmlmodel.VesaliusWSException{
            Code:    "ERROR",
            Message: "No response data!",
        }
    }

    return result, ex
}

func sleep() {
    time.Sleep(1200 * time.Millisecond)
}

func defLogout(token string) {
    _, _ = Logout(token)
}