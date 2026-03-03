package nova_patient_alert

import (
	"vesaliusdr/database"
	"vesaliusdr/model"
	"vesaliusdr/utils"
)

func FindPatientActiveAlertByPrn(prn string) ([]model.NovaPatientAlert, error) {
    lx := make([]model.NovaPatientAlert, 0)
    db := database.GetDbrs()
    if db == nil {
        utils.LogInfo("db is nil")
        return lx, nil
    }

    q := `SELECT ALERT_REF_NO, PRN, ALERT_TYPE, ALLERGY_TYPE, DESCRIPTION,
        SYSTEM, ROUTE, PROBABILITY, REACTION, CREATED_BY,
        CREATION_DATE, INACTIVE_USER, INACTIVE_DATE_TIME, INACTIVE_REASON
        FROM NOVA_PATIENT_ALERT WHERE PRN = :prn AND INACTIVE_DATE_TIME IS NULL 
        ORDER BY ALERT_TYPE`
    rows, err := db.Queryx(q, prn)
    if err != nil {
        utils.LogError(err)
        return lx, err
    }

    defer rows.Close()

    for rows.Next() {
        o := model.DbNovaPatientAlert{}
        err := rows.StructScan(&o)

        if err != nil {
            utils.LogError(err)
            return lx, err
        }

        x := model.NovaPatientAlert{}
        x.FromDbModel(o)
        lx = append(lx, x)
    }

    return lx, nil
}