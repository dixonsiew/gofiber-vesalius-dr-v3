package xmlmodel

import "encoding/xml"

type ResultListDoctorTodoNotification struct {
    XMLName                 xml.Name                 `xml:"Result" json:"-"`
    DoctorTodoNotifications []DoctorTodoNotification `xml:"Notification"`
    Success                 Success                  `xml:"Success"`
    Error                   Error                    `xml:"Error"`
}

type DoctorTodoNotification struct {
    XMLName          xml.Name          `xml:"Notification" json:"-"`
    Prn              string            `xml:"PRN" json:"prn"`
    PatientName      string            `xml:"PatientName" json:"patientName"`
    AccountNo        string            `xml:"AccountNo" json:"accountNo"`
    DoctorName       string            `xml:"DoctorName" json:"doctorName"`
    DischargeSummary string            `xml:"DischargeSummary" json:"dischargeSummary"`
    DrugDiscontinues []DrugDiscontinue `xml:"DrugDiscontinue" json:"drugDiscontinue"`
    VerbalOrders     []VerbalOrder     `xml:"VerbalOrder" json:"verbalOrder"`
}

type DrugDiscontinue struct {
    XMLName xml.Name `xml:"DrugDiscontinue" json:"-"`
    Items   []Item   `xml:"Item" json:"drugItem"`
}

type VerbalOrder struct {
    XMLName xml.Name `xml:"VerbalOrder" json:"-"`
    Items   []Item   `xml:"Item" json:"verbalItem"`
}

type Item struct {
    XMLName           xml.Name `xml:"Item" json:"-"`
    AccessionNo       string   `xml:"AccessionNo" json:"accessionNo"`
    Code              string   `xml:"Code" json:"code"`
    Description       string   `xml:"Description" json:"description"`
    Quantity          string   `xml:"Quantity" json:"quantity"`
    UOM               string   `xml:"UOM" json:"uom"`
    Instruction       string   `xml:"Instruction" json:"instruction"`
    OrderDate         string   `xml:"OrderDate" json:"orderDate"`
    OrderTime         string   `xml:"OrderTime" json:"orderTime"`
    OrderedBy         string   `xml:"OrderedBy" json:"orderedBy"`
    DiscontinueDate   string   `xml:"DiscontinueDate" json:"discontinueDate"`
    DiscontinueTime   string   `xml:"DiscontinueTime" json:"discontinueTime"`
    DiscontinueReason string   `xml:"DiscontinueReason" json:"discontinueReason"`
    DiscontinueBy     string   `xml:"DiscontinueBy" json:"discontinueBy"`
    DispenseQty       string   `xml:"DispenseQty" json:"dispenseQty"`
    LastDispenseTime  string   `xml:"LastDispenseTime" json:"lastDispenseTime"`
    ServedQty         string   `xml:"ServedQty" json:"servedQty"`
    LastServedTime    string   `xml:"LastServedTime" json:"lastServedTime"`
}

type ProcessDoctorToDoAck struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

type ProcessDoctorReviewInvestigationAck struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}