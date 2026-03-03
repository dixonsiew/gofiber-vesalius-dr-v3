package xmlmodel

import "encoding/xml"

type ResultListOutpatient struct {
    XMLName     xml.Name     `xml:"Result" json:"-"`
    Outpatients []Outpatient `xml:"Outpatient"`
    Success     Success      `xml:"Success"`
    Error       Error        `xml:"Error"`
}

type Outpatient struct {
    XMLName  xml.Name          `xml:"Outpatient" json:"-"`
    Patients []OutpatientQueue `xml:"Patient" json:"outpatientQueue"`
}

type OutpatientQueue struct {
    PRN                 string `xml:"PRN" json:"prn"`
    Name                Name   `xml:"Name" json:"name"`
    Sex                 Sex    `xml:"Sex" json:"sex"`
    Age                 string `xml:"Age" json:"age"`
    Nationality         string `xml:"Nationality" json:"nationality"`
    VIPFlag             string `xml:"VIPFlag" json:"vipFlag"`
    VisitType           string `xml:"VisitType" json:"visitType"`
    VisitNumber         string `xml:"VisitNumber" json:"visitNumber"`
    QueueNumber         string `xml:"QueueNumber" json:"queueNumber"`
    QueueCriteria       string `xml:"QueueCriteria" json:"queueCriteria"`
    PatientStatus       string `xml:"PatientStatus" json:"patientStatus"`
    RegistrationDate    string `xml:"RegistrationDate" json:"registrationDate"`
    RegistrationTime    string `xml:"RegistrationTime" json:"registrationTime"`
    VitalsAreAvailable  string `xml:"VitalsAreAvailable" json:"vitalsAreAvailable"`
    TriageScore         string `xml:"TriageScore" json:"triageScore"`
    TriageDiscriminator string `xml:"TriageDiscriminator" json:"triageDiscriminator"`
    HasOnArrivalOrders  string `xml:"HasOnArrivalOrders" json:"hasOnArrivalOrders"`
    RoutedBy            string `xml:"RoutedBy" json:"routedBy"`
    Ward                string `xml:"Ward" json:"ward"`
    Bed                 string `xml:"Bed" json:"bed"`
    AdmissionDate       string `xml:"AdmissionDate" json:"admissionDate"`
    AdmissionTime       string `xml:"AdmissionTime" json:"admissionTime"`
    AppointmentDate     string `xml:"AppointmentDate" json:"appointmentDate"`
    AppointmentTime     string `xml:"AppointmentTime" json:"appointmentTime"`
}