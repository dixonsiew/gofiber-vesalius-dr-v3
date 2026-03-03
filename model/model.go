package model

import (
    "database/sql"
)

type DbDoctorAppUser struct {
    Doctor_app_user_id sql.NullInt32
    Username           sql.NullString
    Password           sql.NullString
    Mcr                sql.NullString
    Email              sql.NullString
    Title              sql.NullString
    First_name         sql.NullString
    Middle_name        sql.NullString
    Last_name          sql.NullString
    DOB                sql.NullString
    Sex                sql.NullString
    Contact_number     sql.NullString
    Nationality        sql.NullString
    First_time_login   sql.NullBool
    Role               sql.NullString
    Branch             sql.NullString
    Machine_id         sql.NullString
}

type DoctorAppUser struct {
    Doctor_app_user_id int    `json:"doctor_app_user_id"`
    Username           string `json:"username"`
    Password           string `json:"-"`
    Mcr                string `json:"mcr"`
    Email              string `json:"email"`
    Title              string `json:"title"`
    First_name         string `json:"first_name"`
    Middle_name        string `json:"middle_name"`
    Last_name          string `json:"last_name"`
    DOB                string `json:"dob"`
    Sex                string `json:"sex"`
    Contact_number     string `json:"contact_number"`
    Nationality        string `json:"nationality"`
    First_time_login   bool   `json:"firstTimeLogin"`
    Role               string `json:"role"`
    Branch             string `json:"branch"`
    Machine_id         string `json:"machine_id"`
}

func (o *DoctorAppUser) FromDbModel(m DbDoctorAppUser) {
    o.Doctor_app_user_id = int(m.Doctor_app_user_id.Int32)
    o.Username = m.Username.String
    o.Password = m.Password.String
    o.Mcr = m.Mcr.String
    o.Email = m.Email.String
    o.Title = m.Title.String
    o.First_name = m.First_name.String
    o.Middle_name = m.Middle_name.String
    o.Last_name = m.Last_name.String
    o.DOB = m.DOB.String
    o.Sex = m.Sex.String
    o.Contact_number = m.Contact_number.String
    o.Nationality = m.Nationality.String
    o.First_time_login = m.First_time_login.Bool
    o.Role = m.Role.String
    o.Branch = m.Branch.String
    o.Machine_id = m.Machine_id.String
}

type DbDoctorRequest struct {
    Doctor_request_id   sql.NullInt32
    Accession_no        sql.NullString
    Create_date_time    sql.NullString
    Doctor_request_type sql.NullString
    Mcr                 sql.NullString
    Notification_type   sql.NullString
    Posted_date_time    sql.NullString
    Posted_flag         sql.NullBool
    Posted_message      sql.NullString
    Prn                 sql.NullString
    Remark              sql.NullString
    Review_date         sql.NullString
    Review_doctor       sql.NullString
    Review_time         sql.NullString
}

type DoctorRequest struct {
    Doctor_request_id   int    `json:"doctorRequestId"`
    Accession_no        string `json:"accessionNo"`
    Create_date_time    string `json:"createDateTime"`
    Doctor_request_type string `json:"doctorRequestType"`
    Mcr                 string `json:"mcr"`
    Notification_type   string `json:"notificationType"`
    Posted_date_time    string `json:"postedDateTime"`
    Posted_flag         bool   `json:"postedFlag"`
    Posted_message      string `json:"postedMessage"`
    Prn                 string `json:"prn"`
    Remark              string `json:"remark"`
    Review_date         string `json:"reviewDate"`
    Review_doctor       string `json:"reviewDoctor"`
    Review_time         string `json:"reviewTime"`
}

func (o *DoctorRequest) FromDbModel(m DbDoctorRequest) {
    o.Doctor_request_id = int(m.Doctor_request_id.Int32)
    o.Accession_no = m.Accession_no.String
    o.Create_date_time = m.Create_date_time.String
    o.Doctor_request_type = m.Doctor_request_type.String
    o.Mcr = m.Mcr.String
    o.Notification_type = m.Notification_type.String
    o.Posted_date_time = m.Posted_date_time.String
    o.Posted_flag = m.Posted_flag.Bool
    o.Posted_message = m.Posted_message.String
    o.Prn = m.Prn.String
    o.Remark = m.Remark.String
    o.Review_date = m.Review_date.String
    o.Review_doctor = m.Review_doctor.String
    o.Review_time = m.Review_time.String
}

type DbDoctorToDoNotification struct {
    Notification_id      sql.NullInt32
    Accession_no         sql.NullString
    Account_no           sql.NullString
    Acknowledgement_flag sql.NullBool
    Discharge_summary    sql.NullString
    Discontinue_by       sql.NullString
    Discontinue_date     sql.NullString
    Discontinue_reason   sql.NullString
    Discontinue_time     sql.NullString
    Dispense_quantity    sql.NullString
    Doctor_name          sql.NullString
    Instruction          sql.NullString
    Item_code            sql.NullString
    Item_desc            sql.NullString
    Item_quantity        sql.NullString
    Last_dispense_time   sql.NullString
    Last_served_time     sql.NullString
    Last_update_date     sql.NullString
    Mcr                  sql.NullString
    Notification_type    sql.NullString
    Order_by             sql.NullString
    Order_date           sql.NullString
    Order_time           sql.NullString
    Patient_name         sql.NullString
    Prn                  sql.NullString
    Served_quantity      sql.NullString
    UOM                  sql.NullString
}

type DoctorToDoNotification struct {
    Notification_id      int    `json:"notification_id"`
    Accession_no         string `json:"accessionNo"`
    Account_no           string `json:"accountNo"`
    Acknowledgement_flag bool   `json:"acknowledgementFlag"`
    Discharge_summary    string `json:"dischargeSummary"`
    Discontinue_by       string `json:"discontinueBy"`
    Discontinue_date     string `json:"discontinueDate"`
    Discontinue_reason   string `json:"discontinueReason"`
    Discontinue_time     string `json:"discontinueTime"`
    Dispense_quantity    string `json:"dispenseQuantity"`
    Doctor_name          string `json:"doctorName"`
    Instruction          string `json:"instruction"`
    Item_code            string `json:"itemCode"`
    Item_desc            string `json:"itemDesc"`
    Item_quantity        string `json:"itemQuantity"`
    Last_dispense_time   string `json:"lastDispenseTime"`
    Last_served_time     string `json:"lastServedTime"`
    Last_update_date     string `json:"lastUpdateDate"`
    Mcr                  string `json:"mcr"`
    Notification_type    string `json:"notificationType"`
    Order_by             string `json:"orderBy"`
    Order_date           string `json:"orderDate"`
    Order_time           string `json:"orderTime"`
    Patient_name         string `json:"patientName"`
    Prn                  string `json:"prn"`
    Served_quantity      string `json:"servedQuantity"`
    UOM                  string `json:"uom"`
}

func (o *DoctorToDoNotification) FromDbModel(m DbDoctorToDoNotification) {
    o.Notification_id = int(m.Notification_id.Int32)
    o.Accession_no = m.Accession_no.String
    o.Account_no = m.Account_no.String
    o.Acknowledgement_flag = m.Acknowledgement_flag.Bool
    o.Discharge_summary = m.Discharge_summary.String
    o.Discontinue_by = m.Discontinue_by.String
    o.Discontinue_date = m.Discontinue_date.String
    o.Discontinue_reason = m.Discontinue_reason.String
    o.Discontinue_time = m.Discontinue_time.String
    o.Dispense_quantity = m.Dispense_quantity.String
    o.Doctor_name = m.Doctor_name.String
    o.Instruction = m.Instruction.String
    o.Item_code = m.Item_code.String
    o.Item_desc = m.Item_desc.String
    o.Item_quantity = m.Item_quantity.String
    o.Last_dispense_time = m.Last_dispense_time.String
    o.Last_served_time = m.Last_served_time.String
    o.Last_update_date = m.Last_update_date.String
    o.Mcr = m.Mcr.String
    o.Notification_type = m.Notification_type.String
    o.Order_by = m.Order_by.String
    o.Order_date = m.Order_date.String
    o.Order_time = m.Order_time.String
    o.Patient_name = m.Patient_name.String
    o.Prn = m.Prn.String
    o.Served_quantity = m.Served_quantity.String
    o.UOM = m.UOM.String
}

type DbInpatientQueueList struct {
    Id                    sql.NullInt32
    Admission_date        sql.NullString
    Admission_time        sql.NullString
    Age                   sql.NullString
    Bed                   sql.NullString
    First_name            sql.NullString
    Has_on_arrival_orders sql.NullString
    Last_name             sql.NullString
    Last_update_date      sql.NullString
    Mcr                   sql.NullString
    Middle_name           sql.NullString
    Nationality           sql.NullString
    Patient_status        sql.NullString
    Patient_type          sql.NullString
    Prn                   sql.NullString
    Queue_criteria        sql.NullString
    Queue_number          sql.NullString
    Registration_date     sql.NullString
    Registration_time     sql.NullString
    Routed_by             sql.NullString
    Sex_code              sql.NullString
    Sex_desc              sql.NullString
    Title                 sql.NullString
    Triage_discriminator  sql.NullString
    Triage_score          sql.NullString
    Vip_flag              sql.NullString
    Visit_number          sql.NullString
    Visit_type            sql.NullString
    Vital_are_available   sql.NullString
    Ward                  sql.NullString
}

type InpatientQueueList struct {
    Id                    int    `json:"id"`
    Admission_date        string `json:"admissionDate"`
    Admission_time        string `json:"admissionTime"`
    Age                   string `json:"age"`
    Bed                   string `json:"bed"`
    First_name            string `json:"firstName"`
    Has_on_arrival_orders string `json:"hasOnArrivalOrders"`
    Last_name             string `json:"lastName"`
    Last_update_date      string `json:"lastUpdateDate"`
    Mcr                   string `json:"mcr"`
    Middle_name           string `json:"middleName"`
    Nationality           string `json:"nationality"`
    Patient_status        string `json:"patientStatus"`
    Patient_type          string `json:"patientType"`
    Prn                   string `json:"prn"`
    Queue_criteria        string `json:"queueCriteria"`
    Queue_number          string `json:"queueNumber"`
    Registration_date     string `json:"registrationDate"`
    Registration_time     string `json:"registrationTime"`
    Routed_by             string `json:"routedBy"`
    Sex_code              string `json:"sexCode"`
    Sex_desc              string `json:"sexDesc"`
    Title                 string `json:"title"`
    Triage_discriminator  string `json:"triageDiscriminator"`
    Triage_score          string `json:"triageScore"`
    Vip_flag              string `json:"vipFlag"`
    Visit_number          string `json:"visitNumber"`
    Visit_type            string `json:"visitType"`
    Vital_are_available   string `json:"vitalAreAvailable"`
    Ward                  string `json:"ward"`
}

func (o *InpatientQueueList) FromDbModel(m DbInpatientQueueList) {
    o.Id = int(m.Id.Int32)
    o.Admission_date = m.Admission_date.String
    o.Admission_time = m.Admission_time.String
    o.Age = m.Age.String
    o.Bed = m.Bed.String
    o.First_name = m.First_name.String
    o.Has_on_arrival_orders = m.Has_on_arrival_orders.String
    o.Last_name = m.Last_name.String
    o.Last_update_date = m.Last_update_date.String
    o.Mcr = m.Mcr.String
    o.Middle_name = m.Middle_name.String
    o.Nationality = m.Nationality.String
    o.Patient_status = m.Patient_status.String
    o.Patient_type = m.Patient_type.String
    o.Prn = m.Prn.String
    o.Queue_criteria = m.Queue_criteria.String
    o.Queue_number = m.Queue_number.String
    o.Registration_date = m.Registration_date.String
    o.Registration_time = m.Registration_time.String
    o.Routed_by = m.Routed_by.String
    o.Sex_code = m.Sex_code.String
    o.Sex_desc = m.Sex_desc.String
    o.Title = m.Title.String
    o.Triage_discriminator = m.Triage_discriminator.String
    o.Triage_score = m.Triage_score.String
    o.Vip_flag = m.Vip_flag.String
    o.Visit_number = m.Visit_number.String
    o.Visit_type = m.Visit_type.String
    o.Vital_are_available = m.Vital_are_available.String
    o.Ward = m.Ward.String
}

type DbOutpatientQueueList struct {
    Id                    sql.NullInt32
    Age                   sql.NullString
    Appointment_date      sql.NullString
    Appointment_time      sql.NullString
    First_name            sql.NullString
    Has_on_arrival_orders sql.NullString
    Last_name             sql.NullString
    Last_update_date      sql.NullString
    Mcr                   sql.NullString
    Middle_name           sql.NullString
    Nationality           sql.NullString
    Patient_status        sql.NullString
    Patient_type          sql.NullString
    Prn                   sql.NullString
    Queue_criteria        sql.NullString
    Queue_number          sql.NullString
    Registration_date     sql.NullString
    Registration_time     sql.NullString
    Routed_by             sql.NullString
    Sex_code              sql.NullString
    Sex_desc              sql.NullString
    Title                 sql.NullString
    Triage_discriminator  sql.NullString
    Triage_score          sql.NullString
    Vip_flag              sql.NullString
    Visit_number          sql.NullString
    Visit_type            sql.NullString
    Vital_are_available   sql.NullString
}

type OutpatientQueueList struct {
    Id                    int    `json:"id"`
    Age                   string `json:"age"`
    Appointment_date      string `json:"appointmentDate"`
    Appointment_time      string `json:"appointmentTime"`
    First_name            string `json:"firstName"`
    Has_on_arrival_orders string `json:"hasOnArrivalOrders"`
    Last_name             string `json:"lastName"`
    Last_update_date      string `json:"lastUpdateDate"`
    Mcr                   string `json:"mcr"`
    Middle_name           string `json:"middleName"`
    Nationality           string `json:"nationality"`
    Patient_status        string `json:"patientStatus"`
    Patient_type          string `json:"patientType"`
    Prn                   string `json:"prn"`
    Queue_criteria        string `json:"queueCriteria"`
    Queue_number          string `json:"queueNumber"`
    Registration_date     string `json:"registrationDate"`
    Registration_time     string `json:"registrationTime"`
    Routed_by             string `json:"routedBy"`
    Sex_code              string `json:"sexCode"`
    Sex_desc              string `json:"sexDesc"`
    Title                 string `json:"title"`
    Triage_discriminator  string `json:"triageDiscriminator"`
    Triage_score          string `json:"triageScore"`
    Vip_flag              string `json:"vipFlag"`
    Visit_number          string `json:"visitNumber"`
    Visit_type            string `json:"visitType"`
    Vital_are_available   string `json:"vitalAreAvailable"`
}

func (o *OutpatientQueueList) FromDbModel(m DbOutpatientQueueList) {
    o.Id = int(m.Id.Int32)
    o.Age = m.Age.String
    o.Appointment_date = m.Appointment_date.String
    o.Appointment_time = m.Appointment_time.String
    o.First_name = m.First_name.String
    o.Has_on_arrival_orders = m.Has_on_arrival_orders.String
    o.Last_name = m.Last_name.String
    o.Last_update_date = m.Last_update_date.String
    o.Mcr = m.Mcr.String
    o.Middle_name = m.Middle_name.String
    o.Nationality = m.Nationality.String
    o.Patient_status = m.Patient_status.String
    o.Patient_type = m.Patient_type.String
    o.Prn = m.Prn.String
    o.Queue_criteria = m.Queue_criteria.String
    o.Queue_number = m.Queue_number.String
    o.Registration_date = m.Registration_date.String
    o.Registration_time = m.Registration_time.String
    o.Routed_by = m.Routed_by.String
    o.Sex_code = m.Sex_code.String
    o.Sex_desc = m.Sex_desc.String
    o.Title = m.Title.String
    o.Triage_discriminator = m.Triage_discriminator.String
    o.Triage_score = m.Triage_score.String
    o.Vip_flag = m.Vip_flag.String
    o.Visit_number = m.Visit_number.String
    o.Visit_type = m.Visit_type.String
    o.Vital_are_available = m.Vital_are_available.String
}

type DbPatientInformation struct {
    Patient_info_id  sql.NullInt32
    Contact_number   sql.NullString
    DOB              sql.NullString
    Document_code    sql.NullString
    Document_desc    sql.NullString
    Document_value   sql.NullString
    First_name       sql.NullString
    Home_address1    sql.NullString
    Home_address2    sql.NullString
    Home_address3    sql.NullString
    Home_address4    sql.NullString
    Home_address5    sql.NullString
    Last_name        sql.NullString
    Last_update_date sql.NullString
    Middle_name      sql.NullString
    Nationality_desc sql.NullString
    Nationalityid    sql.NullString
    Prn              sql.NullString
    Resident         sql.NullString
    Sex_code         sql.NullString
    Sex_desc         sql.NullString
    Title            sql.NullString
    Email            sql.NullString
}

type PatientInformation struct {
    Patient_info_id  int    `json:"patient_info_id"`
    Contact_number   string `json:"contactNumber"`
    DOB              string `json:"dob"`
    Document_code    string `json:"documentCode"`
    Document_desc    string `json:"documentDesc"`
    Document_value   string `json:"documentValue"`
    First_name       string `json:"firstName"`
    Home_address1    string `json:"homeAddress1"`
    Home_address2    string `json:"homeAddress2"`
    Home_address3    string `json:"homeAddress3"`
    Home_address4    string `json:"homeAddress4"`
    Home_address5    string `json:"homeAddress5"`
    Last_name        string `json:"lastName"`
    Last_update_date string `json:"lastUpdateDate"`
    Middle_name      string `json:"middleName"`
    Nationality_desc string `json:"nationalityDesc"`
    Nationalityid    string `json:"nationalityID"`
    Prn              string `json:"prn"`
    Resident         string `json:"resident"`
    Sex_code         string `json:"sexCode"`
    Sex_desc         string `json:"sexDesc"`
    Title            string `json:"title"`
    Email            string `json:"email"`
}

func (o *PatientInformation) FromDbModel(m DbPatientInformation) {
    o.Patient_info_id = int(m.Patient_info_id.Int32)
    o.Contact_number = m.Contact_number.String
    o.DOB = m.DOB.String
    o.Document_code = m.Document_code.String
    o.Document_desc = m.Document_desc.String
    o.Document_value = m.Document_value.String
    o.First_name = m.First_name.String
    o.Home_address1 = m.Home_address1.String
    o.Home_address2 = m.Home_address2.String
    o.Home_address3 = m.Home_address3.String
    o.Home_address4 = m.Home_address4.String
    o.Home_address5 = m.Home_address5.String
    o.Last_name = m.Last_name.String
    o.Last_update_date = m.Last_update_date.String
    o.Middle_name = m.Middle_name.String
    o.Nationality_desc = m.Nationality_desc.String
    o.Nationalityid = m.Nationalityid.String
    o.Prn = m.Prn.String
    o.Resident = m.Resident.String
    o.Sex_code = m.Sex_code.String
    o.Sex_desc = m.Sex_desc.String
    o.Title = m.Title.String
    o.Email = m.Email.String
}

type DbInvestigationReport struct {
    Id                        sql.NullInt32
    Accession_no              sql.NullString
    Account_no                sql.NullString
    Investigation_type        sql.NullString
    Last_update_date          sql.NullString
    Order_date                sql.NullString
    Order_doctor_mcr          sql.NullString
    Patient_name              sql.NullString
    Prn                       sql.NullString
    Remark                    sql.NullString
    Report_type               sql.NullString
    Result                    sql.NullString
    Result_date               sql.NullString
    Review_investigation_flag sql.NullBool
    Service_code              sql.NullString
    Service_desc              sql.NullString
}

type InvestigationReport struct {
    Id                        int     `json:"id"`
    Accession_no              string  `json:"accessionNo"`
    Account_no                string  `json:"accountNo"`
    Investigation_type        string  `json:"investigationType"`
    Last_update_date          string  `json:"lastUpdateDate"`
    Order_date                string  `json:"orderDate"`
    Order_doctor_mcr          string  `json:"orderDoctorMcr"`
    Patient_name              string  `json:"patientName"`
    Prn                       string  `json:"prn"`
    Remark                    string  `json:"remark"`
    Report_type               string  `json:"reportType"`
    Result                    *string `json:"result"`
    Result_date               string  `json:"resultDate"`
    Review_investigation_flag bool    `json:"reviewInvestigationFlag"`
    Service_code              string  `json:"serviceCode"`
    Service_desc              string  `json:"serviceDesc"`
}

func (o *InvestigationReport) FromDbModel(m DbInvestigationReport) {
    o.Id = int(m.Id.Int32)
    o.Accession_no = m.Accession_no.String
    o.Account_no = m.Account_no.String
    o.Investigation_type = m.Investigation_type.String
    o.Last_update_date = m.Last_update_date.String
    o.Order_date = m.Order_date.String
    o.Order_doctor_mcr = m.Order_doctor_mcr.String
    o.Patient_name = m.Patient_name.String
    o.Prn = m.Prn.String
    o.Remark = m.Remark.String
    o.Report_type = m.Report_type.String
    o.Result = &m.Result.String
    o.Result_date = m.Result_date.String
    o.Review_investigation_flag = m.Review_investigation_flag.Bool
    o.Service_code = m.Service_code.String
    o.Service_desc = m.Service_desc.String
}

type DbNovaPatientAlert struct {
    ALERT_REF_NO       sql.NullInt32
    PRN                sql.NullString
    ALERT_TYPE         sql.NullString
    ALLERGY_TYPE       sql.NullString
    DESCRIPTION        sql.NullString
    SYSTEM             sql.NullString
    ROUTE              sql.NullString
    PROBABILITY        sql.NullString
    REACTION           sql.NullString
    CREATED_BY         sql.NullString
    CREATION_DATE      sql.NullString
    INACTIVE_USER      sql.NullString
    INACTIVE_DATE_TIME sql.NullString
    INACTIVE_REASON    sql.NullString
}

type NovaPatientAlert struct {
    ALERT_REF_NO       int    `json:"ialertRefNod"`
    PRN                string `json:"prn"`
    ALERT_TYPE         string `json:"alertType"`
    ALLERGY_TYPE       string `json:"allergyType"`
    DESCRIPTION        string `json:"description"`
    SYSTEM             string `json:"system"`
    ROUTE              string `json:"route"`
    PROBABILITY        string `json:"probability"`
    REACTION           string `json:"reaction"`
    CREATED_BY         string `json:"createdBy"`
    CREATION_DATE      string `json:"creationDate"`
    INACTIVE_USER      string `json:"inactiveUser"`
    INACTIVE_DATE_TIME string `json:"inactiveDateTime"`
    INACTIVE_REASON    string `json:"inactiveReason"`
}

func (o *NovaPatientAlert) FromDbModel(m DbNovaPatientAlert) {
    o.ALERT_REF_NO = int(m.ALERT_REF_NO.Int32)
    o.PRN = m.PRN.String
    o.ALERT_TYPE = m.ALERT_TYPE.String
    o.ALLERGY_TYPE = m.ALLERGY_TYPE.String
    o.DESCRIPTION = m.DESCRIPTION.String
    o.SYSTEM = m.SYSTEM.String
    o.ROUTE = m.ROUTE.String
    o.PROBABILITY = m.PROBABILITY.String
    o.REACTION = m.REACTION.String
    o.CREATED_BY = m.CREATED_BY.String
    o.CREATION_DATE = m.CREATION_DATE.String
    o.INACTIVE_USER = m.INACTIVE_USER.String
    o.INACTIVE_DATE_TIME = m.INACTIVE_DATE_TIME.String
    o.INACTIVE_REASON = m.INACTIVE_REASON.String
}

type OutpatientQueueSummary struct {
    QueueCriteria      string `json:"queueCriteria"`
    QueueCount         int    `json:"queueCount"`
    ImageName          string `json:"imageName"`
    LinkName           string `json:"linkName"`
    LastUpdateDateTime string `json:"lastUpdateDateTime"`
}