package outpatient_queue_list

import (
	"database/sql"
	"vesaliusdr/database"
	"vesaliusdr/model"
	"vesaliusdr/utils"

	"github.com/gofiber/fiber/v2"
)

func GetOutpatientDetailByPrn(mcr string, prn string) (*model.OutpatientQueueList, error) {
    o := model.DbOutpatientQueueList{}
	k := model.OutpatientQueueList{}
    var x *model.OutpatientQueueList
	db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return x, nil
    }

    q := `select id, age, appointment_date, appointment_time, first_name, 
        has_on_arrival_orders, last_name, last_update_date, mcr, middle_name, 
        nationality, patient_status, patient_type, prn, queue_criteria, 
        queue_number, registration_date, registration_time, routed_by, sex_code, 
        sex_desc, title, triage_discriminator, triage_score, vip_flag, 
        visit_number, visit_type, vital_are_available
        from outpatient_queue_list where mcr = @p1 and prn = @p2`
    rows, err := db.Queryx(q, mcr, prn)
    if err != nil {
        utils.LogError(err)
        return x, err
    }

    defer rows.Close()

    if rows.Next() {
        err := rows.StructScan(&o)

        if err != nil {
            utils.LogError(err)
            return x, err
        }

        k.FromDbModel(o)
        x = &k
    }

    return x, nil
}

func GetOutpatientQueueDetailListByMcr(mcr string, queueCriteria string) ([]model.OutpatientQueueList, error) {
    lx := make([]model.OutpatientQueueList, 0)
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return lx, nil
    }

    q := `select id, age, appointment_date, appointment_time, first_name, 
        has_on_arrival_orders, last_name, last_update_date, mcr, middle_name, 
        nationality, patient_status, patient_type, prn, queue_criteria, 
        queue_number, registration_date, registration_time, routed_by, sex_code, 
        sex_desc, title, triage_discriminator, triage_score, vip_flag, 
        visit_number, visit_type, vital_are_available
        from outpatient_queue_list where mcr = @p1 and queue_criteria = @p2
        order by prn asc`
    rows, err := db.Queryx(q, mcr, queueCriteria)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.DbOutpatientQueueList{}
        err := rows.StructScan(&o)

        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        x := model.OutpatientQueueList{}
        x.FromDbModel(o)
        lx = append(lx, x)
    }

    return lx, nil
}

func GetLastUpdateDateByMCR(mcr string) (string, error) {
    s := ""
    db := database.GetDb()
	if db == nil {
		utils.LogInfo("db is nil")
		return s, nil
	}

    q := `select distinct top 1 last_update_date from outpatient_queue_list where mcr = @p1 order by last_update_date asc`
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

func GetOutpatientQueueSummaryByMCR(mcr string) ([]fiber.Map, error) {
    lx := make([]fiber.Map, 0)
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return lx, nil
    }

    q := `select queue_criteria as queueCriteria, count(*) as queueCount from outpatient_queue_list where mcr = @p1 group by queue_criteria`
    rows, err := db.Query(q, mcr)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        var queueCriteria sql.NullString
        var queueCount sql.NullInt32
        err := rows.Scan(&queueCriteria, &queueCount)

        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        x := fiber.Map{
            "queueCriteria": queueCriteria.String,
            "queueCount": int(queueCount.Int32),
        }
        lx = append(lx, x)
    }

    return lx, nil
}

func DeleteOutpatientQueueListByMcr(mcr string) error {
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return nil
    }

    q := `delete from outpatient_queue_list where mcr = @p1`
    _, err := db.Exec(q, mcr)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func AddOutpatientQueueList(o model.OutpatientQueueList) error {
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return nil
    }

    q := `insert into outpatient_queue_list (id, age, appointment_date, appointment_time, first_name, 
        has_on_arrival_orders, last_name, last_update_date, mcr, middle_name, 
        nationality, patient_status, patient_type, prn, queue_criteria, 
        queue_number, registration_date, registration_time, routed_by, sex_code, 
        sex_desc, title, triage_discriminator, triage_score, vip_flag, 
        visit_number, visit_type, vital_are_available) values(
        (select isnull(max(id) + 1, 0) from outpatient_queue_list with(serializable, updlock)),
        @p1, @p2, @p3, @p4,
        @p5, @p6, @p7, @p8, @p9,
        @p10, @p11, @p12, @p13, @p14,
        @p15, @p16, @p17, @p18, @p19,
        @p20, @p21, @p22, @p23, @p24,
        @p25, @p26, @p27)`
    _, err := db.Exec(q, o.Age, o.Appointment_date, o.Appointment_time, o.First_name,
        o.Has_on_arrival_orders, o.Last_name, o.Last_update_date, o.Mcr, o.Middle_name,
        o.Nationality, o.Patient_status, o.Patient_type, o.Prn, o.Queue_criteria,
        o.Queue_number, o.Registration_date, o.Registration_time, o.Routed_by, o.Sex_code,
        o.Sex_desc, o.Title, o.Triage_discriminator, o.Triage_score, o.Vip_flag,
        o.Visit_number, o.Visit_type, o.Vital_are_available)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}