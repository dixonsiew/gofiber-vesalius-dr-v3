package patient_info

import (
	"vesaliusdr/database"
	"vesaliusdr/model"
	"vesaliusdr/utils"
)

func FindByPrn(prn string) (*model.PatientInformation, error) {
	o := model.DbPatientInformation{}
    k := model.PatientInformation{}
    var x *model.PatientInformation
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return x, nil
    }

	q := `select patient_info_id, contact_number, dob, document_code, document_desc, 
	    document_value, first_name, home_address1, home_address2, home_address3, 
		home_address4, home_address5, last_name, last_update_date, middle_name, 
		nationality_desc, nationalityid, prn, resident, sex_code, 
		sex_desc, title, email 
		from patient_information where prn = @p1`
	rows, err := db.Queryx(q, prn)
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

func AddPatientInfo(o model.PatientInformation) error {
	db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return nil
    }

	q := `insert into patient_information (patient_info_id, contact_number, dob, document_code, document_desc, 
	    document_value, first_name, home_address1, home_address2, home_address3, 
		home_address4, home_address5, last_name, last_update_date, middle_name, 
		nationality_desc, nationalityid, prn, resident, sex_code, 
		sex_desc, title, email) values(next value for patient_info_id_sequence, @p1, @p2, @p3, @p4, 
		@p5, @p6, @p7, @p8, @p9, 
		@p10, @p11, @p12, @p13, @p14,
		@p15, @p16, @p17, @p18, @p19,
		@p20, @p21, @p22)`
	_, err := db.Exec(q, o.Contact_number, o.DOB, o.Document_code, o.Document_desc,
	    o.Document_value, o.First_name, o.Home_address1, o.Home_address2, o.Home_address3,
	    o.Home_address4, o.Home_address5, o.Last_name, o.Last_update_date, o.Middle_name,
	    o.Nationality_desc, o.Nationalityid, o.Prn, o.Resident, o.Sex_code,
	    o.Sex_desc, o.Title, o.Email)
	if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}