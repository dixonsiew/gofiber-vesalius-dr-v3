package doctor_todo_notification

import (
	"database/sql"
	"vesaliusdr/database"
	"vesaliusdr/model"
	"vesaliusdr/utils"
)

func GetLastUpdateDateByMcr(mcr string) (string, error) {
	s := ""
	db := database.GetDb()
	if db == nil {
		utils.LogInfo("db is nil")
		return s, nil
	}

    q := `select distinct top 1 last_update_date from doctor_todo_notification where mcr = @p1 order by last_update_date asc`
    rows, err := db.Query(q, mcr)
    if err != nil {
        utils.LogError(err)
        return s, err
    }

    defer rows.Close()

    if rows.Next() {
        var o sql.NullString
        err := rows.Scan(&o)

        if err != nil {
            utils.LogError(err)
            return s, err
        }

        s = o.String
    }

    return s, nil
}

func DeleteDoctorToDoNotificationListByMcr(mcr string) error {
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return nil
    }

    q := `delete from doctor_todo_notification where mcr = @p1`
    _, err := db.Exec(q, mcr)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func GetDoctorToDoNotificationListByMcr(mcr string) ([]model.DoctorToDoNotification, error) {
    lx := make([]model.DoctorToDoNotification, 0)
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return lx, nil
    }

    q := `select notification_id, accession_no, account_no, acknowledgement_flag, discharge_summary, 
        discontinue_by, discontinue_date, discontinue_reason, discontinue_time, dispense_quantity, 
        doctor_name, instruction, item_code, item_desc, item_quantity, 
        last_dispense_time, last_served_time, last_update_date, mcr, notification_type, 
        order_by, order_date, order_time, patient_name, prn, 
        served_quantity, uom 
        from doctor_todo_notification where mcr = @p1 and acknowledgement_flag = 0`
    rows, err := db.Queryx(q, mcr)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.DbDoctorToDoNotification{}
        err := rows.StructScan(&o)

        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        x := model.DoctorToDoNotification{}
        x.FromDbModel(o)
        lx = append(lx, x)
    }

    return lx, nil
}

func GetDoctorToDoNotificationListByMcrPrnNotificationType(mcr string, prn string) ([]model.DoctorToDoNotification, error) {
    lx := make([]model.DoctorToDoNotification, 0)
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return lx, nil
    }

    q := `select notification_id, accession_no, account_no, acknowledgement_flag, discharge_summary, 
        discontinue_by, discontinue_date, discontinue_reason, discontinue_time, dispense_quantity, 
        doctor_name, instruction, item_code, item_desc, item_quantity, 
        last_dispense_time, last_served_time, last_update_date, mcr, notification_type, 
        order_by, order_date, order_time, patient_name, prn, 
        served_quantity, uom 
        from doctor_todo_notification where mcr = @p1 and prn = @p2 and notification_type <> 'Discharge Summary' and acknowledgement_flag = 0`
    rows, err := db.Queryx(q, mcr, prn)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.DbDoctorToDoNotification{}
        err := rows.StructScan(&o)

        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        x := model.DoctorToDoNotification{}
        x.FromDbModel(o)
        lx = append(lx, x)
    }

    return lx, nil
}

func UpdateDoctorToDoNotificationAcknowledgementFlag(mcr string, accessionNo string) error {
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return nil
    }

	q := `update doctor_todo_notification set acknowledgement_flag = 1 where mcr = @p1 and accession_no = @p2`
	_, err := db.Exec(q, mcr, accessionNo)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func AddDoctorToDoNotification(o model.DoctorToDoNotification) error {
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return nil
    }

    q := `insert into doctor_todo_notification (notification_id, accession_no, account_no, acknowledgement_flag, discharge_summary, 
        discontinue_by, discontinue_date, discontinue_reason, discontinue_time, dispense_quantity, 
        doctor_name, instruction, item_code, item_desc, item_quantity, 
        last_dispense_time, last_served_time, last_update_date, mcr, notification_type, 
        order_by, order_date, order_time, patient_name, prn, 
        served_quantity, uom) values(
        (select isnull(max(notification_id) + 1, 0) from doctor_todo_notification with(serializable, updlock)),
        @p1, @p2, 0, @p3,
        @p4, @p5, @p6, @p7, @p8,
        @p9, @p10, @p11, @p12, @p13,
        @p14, @p15, @p16, @p17, @p18,
        @p19, @p20, @p21, @p22, @p23,
        @p24, @p25)`
    _, err := db.Exec(q, o.Accession_no, o.Account_no, o.Discharge_summary,
	    o.Discontinue_by, o.Discontinue_date, o.Discontinue_reason, o.Discontinue_time, o.Dispense_quantity,
	    o.Doctor_name, o.Instruction, o.Item_code, o.Item_desc, o.Item_quantity,
        o.Last_dispense_time, o.Last_served_time, o.Last_update_date, o.Mcr, o.Notification_type,
        o.Order_by, o.Order_date, o.Order_time, o.Patient_name, o.Prn,
        o.Served_quantity, o.UOM)
	if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}