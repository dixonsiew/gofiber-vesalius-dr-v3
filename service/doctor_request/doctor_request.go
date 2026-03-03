package doctor_request

import (
	"vesaliusdr/database"
	"vesaliusdr/model"
	"vesaliusdr/utils"
)

func GetDoctorRequestByCreateDateTimeAsc() (*model.DoctorRequest, error) {
	o := model.DbDoctorRequest{}
	k := model.DoctorRequest{}
    var x *model.DoctorRequest
	db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return x, nil
    }

	q := `select top 1 doctor_request_id, accession_no, create_date_time, doctor_request_type, mcr, 
	    notification_type, posted_date_time, posted_flag, posted_message, prn, 
		remark, review_date, review_doctor, review_time
		from doctor_request
		where posted_flag = 0
		order by create_date_time, doctor_request_id asc`
	rows, err := db.Queryx(q)
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

func UpdateDoctorRequestStatus(accessionNo string, postedDateTime string, postedMessage string) error {
	db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return nil
    }

	q := `update doctor_request set posted_flag = 1, posted_date_time = @p1, posted_message = @p2 where accession_no = @p3`
	_, err := db.Exec(q, postedDateTime, postedMessage, accessionNo)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func AddDoctorRequest(o model.DoctorRequest) error {
	db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return nil
    }

	q := `insert into doctor_request (doctor_request_id, accession_no, create_date_time, doctor_request_type, mcr, 
	    notification_type, posted_date_time, posted_flag, posted_message, prn, 
		remark, review_date, review_doctor, review_time) values(
        next value for doctor_request_id_sequence,
        @p1, @p2, @p3, @p4,
		@p5, @p6, @p7, @p8, @p9, 
		@p10, @p11, @p12, @p13)`
	_, err := db.Exec(q, o.Accession_no, o.Create_date_time, o.Doctor_request_type, o.Mcr,
	    o.Notification_type, o.Posted_date_time, o.Posted_flag, o.Posted_message, o.Prn,
	    o.Remark, o.Review_date, o.Review_doctor, o.Review_time)
	if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}