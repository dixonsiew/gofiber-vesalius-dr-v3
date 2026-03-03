package xmlmodel

import "encoding/xml"

type ResultListInvestigation struct {
	XMLName        xml.Name        `xml:"Result" json:"-"`
	Investigations []Investigation `xml:"Investigation"`
	Success        Success         `xml:"Success"`
	Error          Error           `xml:"Error"`
}

type Investigation struct {
	XMLName xml.Name `xml:"Investigation" json:"-"`
	Visits  []Visit  `xml:"Visit" json:"visitList"`
}

type Visit struct {
	XMLName     xml.Name  `xml:"Visit" json:"-"`
	AccountNo   string    `xml:"AccountNo" json:"accountNo"`
	PRN         string    `xml:"PRN" json:"prn"`
	PatientName string    `xml:"patientName" json:"patientName"`
	Reports     []Service `xml:"Reports" json:"reports"`
}

type Service struct {
	XMLName           xml.Name           `xml:"Service" json:"-"`
	ServiceComponents []ServiceComponent `xml:"ServiceComponent"`
}

type ServiceComponent struct {
	XMLName        xml.Name `xml:"ServiceComponent" json:"-"`
	AccessionNo    string   `xml:"AccessionNo" json:"accessionNo"`
	Type           string   `xml:"Type" json:"investigationType"`
	Code           string   `xml:"Code" json:"serviceCode"`
	Description    string   `xml:"Description" json:"serviceDesc"`
	OrderDoctorMcr string   `xml:"OrderDoctorMcr" json:"orderDoctorMcr"`
	OrderDate      string   `xml:"OrderDate" json:"orderDate"`
	ResultDate     string   `xml:"ResultDate" json:"resultDate"`
	ReportType     string   `xml:"ReportType" json:"reportType"`
	Result         string   `xml:"Result" json:"result"`
	Remark         string   `xml:"Remark" json:"remark"`
}