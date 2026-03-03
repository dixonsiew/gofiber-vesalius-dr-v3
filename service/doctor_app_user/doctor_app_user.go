package doctor_app_user

import (
	"vesaliusdr/database"
	"vesaliusdr/model"
	"vesaliusdr/utils"

	"golang.org/x/crypto/bcrypt"
)

func FindByUserId(userId int) (*model.DoctorAppUser, error) {
    o := model.DbDoctorAppUser{}
    k := model.DoctorAppUser{}
    var x *model.DoctorAppUser
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return x, nil
    }

    q := `select doctor_app_user_id, branch, contact_number, dob, email, 
        first_time_login, first_name, last_name, mcr, middle_name, 
        nationality, password, [role], sex, title, 
        username, machine_id 
        from doctor_app_user where doctor_app_user_id = @p1`
    rows, err := db.Queryx(q, userId)
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

func FindByUsername(username string) (*model.DoctorAppUser, error) {
    o := model.DbDoctorAppUser{}
    k := model.DoctorAppUser{}
    var x *model.DoctorAppUser
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return x, nil
    }

    q := `select doctor_app_user_id, branch, contact_number, dob, email, 
        first_time_login, first_name, last_name, mcr, middle_name, 
        nationality, password, [role], sex, title, 
        username, machine_id 
        from doctor_app_user where username = @p1`
    rows, err := db.Queryx(q, username)
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

func FindByEmail(email string) (*model.DoctorAppUser, error) {
    o := model.DbDoctorAppUser{}
    k := model.DoctorAppUser{}
    var x *model.DoctorAppUser
    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return x, nil
    }

    q := `select doctor_app_user_id, branch, contact_number, dob, email, 
        first_time_login, first_name, last_name, mcr, middle_name, 
        nationality, password, [role], sex, title, 
        username, machine_id 
        from doctor_app_user where email = @p1`
    rows, err := db.Queryx(q, email)
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

func UpdateMachineId(o model.DoctorAppUser) error {
    machine_id := []byte(o.Machine_id)
    hashedMachineId, err := bcrypt.GenerateFromPassword(machine_id, bcrypt.DefaultCost)
    if err != nil {
        utils.LogError(err)
        return err
    }

    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return nil
    }

    q := `update doctor_app_user set machine_id = @p1 where doctor_app_user_id = @p2`
    _, err = db.Exec(q, string(hashedMachineId), o.Doctor_app_user_id)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func FirstTimeChangePassword(o model.DoctorAppUser) error {
    pw := []byte(o.Password)
    pwd, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
    if err != nil {
        utils.LogError(err)
        return err
    }

    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return nil
    }

    q := `update doctor_app_user set password = @p1, first_time_login = @p2 where doctor_app_user_id = @p3`
    _, err = db.Exec(q, string(pwd), o.First_time_login, o.Doctor_app_user_id)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func UpdatePassword(o model.DoctorAppUser) error {
    pw := []byte(o.Password)
    pwd, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
    if err != nil {
        utils.LogError(err)
        return err
    }

    db := database.GetDb()
    if db == nil {
        utils.LogInfo("db is nil")
        return nil
    }

    q := `update doctor_app_user set password = @p1 where doctor_app_user_id = @p2`
    _, err = db.Exec(q, string(pwd), o.Doctor_app_user_id)
    if err != nil {
        utils.LogError(err)
        return err
    }

    return nil
}

func ValidateCredentials(user model.DoctorAppUser, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    return err == nil
}

func ValidateCredentials2(user model.DoctorAppUser, machineId string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(user.Machine_id), []byte(machineId))
    return err == nil
}