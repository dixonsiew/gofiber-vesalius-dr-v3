package dto

type LoginDto struct {
    Username      string `json:"username" validate:"required,min=1" default:"nova.doctor"`
    Password      string `json:"password" validate:"required,min=1" default:"password"`
    FromBiometric int    `json:"fromBiometric"`
}

type PostChangePasswordDto struct {
    OldPassword string `json:"oldPassword" validate:"required,min=1"`
    NewPassword string `json:"newPassword" validate:"required,min=1"`
}

type PostMachineInfo struct {
    MachineId string `json:"machineId" validate:"required,min=1"`
}

type PostProcessDoctorToDoAcknowledgeDto struct {
    NotificationType string `json:"notificationType"`
    AccessionNo      string `json:"accessionNo"`
    Remark           string `json:"remark"`
}

type PostProcessDoctorReviewInvestigationDto struct {
    PRN          string `json:"prn"`
    AccessionNo  string `json:"accessionNo"`
    ReviewDoctor string `json:"reviewDoctor"`
    ReviewDate   string `json:"reviewDate"`
    ReviewTime   string `json:"reviewTime"`
}