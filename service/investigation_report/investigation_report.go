package investigation_report

import (
    "database/sql"
    "vesaliusdr/database"
    "vesaliusdr/model"
    "vesaliusdr/utils"
)

func GetLastUpdateDateByOrderDate(orderDate string) (string, error) {
    s := ""
    db := database.GetDb()
	if db == nil {
		utils.LogInfo("db is nil")
		return s, nil
	}

    q := `select distinct top 1 last_update_date from investigation_report where order_date = @p1 order by last_update_date asc`
    rows, err := db.Query(q, orderDate)
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

func DeleteInvestigationListByOrderDate(orderDate string) error {
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return nil
    }

    q := `delete from investigation_report where order_date = @p1`
    _, err := db.Exec(q, orderDate)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func GetInvestigationListByOrderDate(mcr string, orderDate string) ([]model.InvestigationReport, error) {
    lx := make([]model.InvestigationReport, 0)
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return lx, nil
    }

    q := `select id, accession_no, account_no, investigation_type, last_update_date, 
        order_date, order_doctor_mcr, patient_name, prn, remark, 
        report_type, [result], result_date, review_investigation_flag, service_code, 
        service_desc 
        from investigation_report 
        where order_doctor_mcr = @p1 and order_date = @p2 and review_investigation_flag = '0' 
        order by prn asc`
    rows, err := db.Queryx(q, mcr, orderDate)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.DbInvestigationReport{}
        err := rows.StructScan(&o)

        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        x := model.InvestigationReport{}
        x.FromDbModel(o)
        lx = append(lx, x)
    }

    return lx, nil
}

func GetInvestigationPDFResultByAccessionNo(mcr string, accessionNo string) (string, error) {
    s := ""
    db := database.GetDb()
	if db == nil {
		utils.LogInfo("db is nil")
		return s, nil
	}

    q := `select result from investigation_report where order_doctor_mcr = @p1 and accession_no = @p2 and report_type = 'PDF' and review_investigation_flag = '0'`
    rows, err := db.Query(q, mcr, accessionNo)
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

func UpdateReviewInvestigationFlag(accessionNo string) error {
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return nil
    }

    q := `update investigation_report set review_investigation_flag = 1 where accession_no = @p1`
    _, err := db.Exec(q, accessionNo)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func AddInvestigationReport(o model.InvestigationReport) error {
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return nil
    }

    q := `insert into investigation_report (id, accession_no, account_no, investigation_type, last_update_date, 
        order_date, order_doctor_mcr, patient_name, prn, remark, 
        report_type, [result], result_date, review_investigation_flag, service_code, 
        service_desc) values(
        (select isnull(max(id) + 1, 0) from investigation_report with(serializable, updlock)),
        @p1, @p2, @p3, @p4,
        @p5, @p6, @p7, @p8, @p9,
        @p10, @p11, @p12, @p13, @p14,
        @p15)`
    _, err := db.Exec(q, o.Accession_no, o.Account_no, o.Investigation_type, o.Last_update_date,
        o.Order_date, o.Order_doctor_mcr, o.Patient_name, o.Prn, o.Remark,
        o.Report_type, *o.Result, o.Result_date, o.Review_investigation_flag, o.Service_code,
        o.Service_desc)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}