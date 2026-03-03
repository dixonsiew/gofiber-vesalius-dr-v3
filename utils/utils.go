package utils

import (
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
    "vesaliusdr/model"

    "github.com/go-playground/validator/v10"
    "github.com/go-resty/resty/v2"
    "github.com/gofiber/fiber/v3"
    "github.com/golang-jwt/jwt/v5"
    "github.com/rs/zerolog"
    "github.com/ztrue/tracerr"
)

type StructValidator struct {
    Xvalidate *validator.Validate
}

func (v *StructValidator) Validate(out any) error {
    return v.Xvalidate.Struct(out)
}

var (
    Logger  zerolog.Logger
    iLogger zerolog.Logger
    client  *resty.Client
)

func SetClient() {
    client = resty.New()
    client.SetTimeout(time.Minute * 5)
}

func SetLogger(runLogFile *os.File) {
    multi := zerolog.MultiLevelWriter(os.Stdout, runLogFile)
    Logger = zerolog.New(multi).Level(zerolog.ErrorLevel).With().Timestamp().Caller().Logger()

    iLogger = zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Timestamp().Logger()
}

func GetValidationErrors(errs validator.ValidationErrors) error {
    if len(errs) > 0 {
        errMsgs := make([]string, 0)
        for _, err := range errs {
            errMsgs = append(errMsgs, fmt.Sprintf(
                "[%s]: '%v' | Needs to implement '%s'",
                err.Field(),
                err.Value(),
                err.Tag(),
            ))
        }

        return &fiber.Error{
            Code:    fiber.ErrBadRequest.Code,
            Message: strings.Join(errMsgs, " and "),
        }
    }

    return nil
}

func GenerateToken(user model.DoctorAppUser) (string, error) {
    claims := jwt.MapClaims{
        "username": user.Username,
        "subject":  fmt.Sprintf("%d", user.Doctor_app_user_id),
        "exp":      time.Now().Add(time.Hour * 87600).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    t, err := token.SignedString([]byte(JWT_SECRET))
    if err != nil {
        LogError(err)
        return "", err
    }

    return t, nil
}

func DecodeToken(c fiber.Ctx) (string, int, error) {
    tokenStr := c.Get("Authorization")
    tokenStr = strings.ReplaceAll(tokenStr, "Bearer ", "")
    token, err := jwt.ParseWithClaims(tokenStr, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(JWT_SECRET), nil
    })

    if err != nil {
        LogError(err)
        return "", 0, err
    }

    claims, ok := token.Claims.(*jwt.MapClaims)
    if !ok {
        return "", 0, fmt.Errorf("could not parse claims")
    }

    sub := (*claims)["subject"].(string)
    username := (*claims)["username"].(string)
    id, _ := strconv.Atoi(sub)
    return username, id, nil
}

func GenerateRefreshToken(user model.DoctorAppUser) (string, error) {
    claims := jwt.MapClaims{
        "username": user.Username,
        "subject":  fmt.Sprintf("%d", user.Doctor_app_user_id),
        "exp":      time.Now().Add(time.Hour * 87600).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    t, err := token.SignedString([]byte(JWT_SECRET))
    if err != nil {
        LogError(err)
        return "", err
    }

    return t, nil
}

func GetErrors(errs []error) string {
    ls := []string{}
    for _, err := range errs {
        ls = append(ls, err.Error())
    }

    return strings.Join(ls, "|")
}

func GetR(action string) *resty.Request {
    return client.R().
        SetHeader("Content-Type", "text/xml; charset=utf-8").
        SetHeader("SOAPAction", fmt.Sprintf("urn:%s", action))
}

func GetXmlResult(body string) []byte {
    i := strings.Index(body, "<ns:return>")
    j := strings.Index(body, "</ns:return>")
    content := body[11+i : j]
    content = strings.ReplaceAll(content, "&lt;", "<")
    content = strings.ReplaceAll(content, "&gt;", ">")
    content = strings.ReplaceAll(content, `encoding="UTF-8">`, `encoding="UTF-8"?>`)
    bx := []byte(content)
    return bx
}

func CatchPanic(funcName string) {
    if err := recover(); err != nil {
        LogError(fmt.Errorf("recovered from panic -%s:%v", funcName, err))
    }
}

func LogError(err error) {
    ex := tracerr.Wrap(err)
    Logger.Err(err).Msg(tracerr.Sprint(ex))
}

func LogInfo(s string) {
    iLogger.Info().Msg(s)
}
