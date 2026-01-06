package service

import (
	"bytes"
	_ "crypto/hmac"
	"encoding/json"
	"net/http"

	"crypto/sha256"

	"database/sql"

	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper"

	"github.com/rest/api/internal/util"
	"github.com/sirupsen/logrus"
)

var _logger = logrus.New()

type APIServiceEndpoint interface {
	Init(config []byte, dbConnection *util.DBConnectionWrapper, verbose bool) error
	AddRouters(base string, router *gin.Engine)
}

// APIResponse returns the service response
type APIResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"serviceMessage"`
	Payload    interface{} `json:"payload,omitempty"`
	ServiceTS  string      `json:"ts"`
	IsSuccess  bool        `json:"isSuccess"`
	Token      *string     `json:"token,omitempty"`
}

func parseInput(c *gin.Context, obj interface{}) bool {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		_logger.Errorf("Error in reading the request body %v", err)
		return false
	}
	if err = json.Unmarshal(bodyBytes, &obj); err != nil {
		_logger.Errorf("Error in parsing request body to %v %v", reflect.TypeOf(obj), err)
		return false
	}
	return true
}

func buildResponse(status int, isOk bool, msg string, payload interface{}) APIResponse {
	return APIResponse{
		StatusCode: status,
		IsSuccess:  isOk,
		Message:    msg,
		Payload:    payload,
		ServiceTS:  time.Now().Format("2006-01-02-15:04:05.000"),
	}
}

func BuildResponse500(msg string, payload interface{}) APIResponse {
	return APIResponse{
		StatusCode: 500,
		IsSuccess:  false,
		Message:    msg,
		Payload:    payload,
		ServiceTS:  time.Now().Format("2006-01-02-15:04:05.000"),
	}
}

func BuildResponse400(msg string) APIResponse {
	return APIResponse{
		StatusCode: 400,
		IsSuccess:  false,
		Message:    msg,
		ServiceTS:  time.Now().Format("2006-01-02-15:04:05.000"),
	}
}

func BuildResponse404(msg string, success bool) APIResponse {
	return APIResponse{
		StatusCode: 404,
		IsSuccess:  success,
		Message:    msg,
		ServiceTS:  time.Now().Format("2006-01-02-15:04:05.000"),
	}
}

func BuildResponse200(msg string, payload interface{}) APIResponse {
	return APIResponse{
		StatusCode: 200,
		IsSuccess:  true,
		Message:    msg,
		Payload:    payload,
		ServiceTS:  time.Now().Format("2006-01-02-15:04:05.000"),
	}
}
func BuildResponse403(msg string, payload interface{}) APIResponse {
	return APIResponse{
		StatusCode: 403,
		IsSuccess:  false,
		Message:    msg,
		Payload:    payload,
		ServiceTS:  time.Now().Format("2006-01-02-15:04:05.000"),
	}
}
func BuildResponse401(msg string) APIResponse {
	return APIResponse{
		StatusCode: 401,
		IsSuccess:  false,
		Message:    "Unauthorized",
		Payload:    msg,
		ServiceTS:  time.Now().Format("2006-01-02-15:04:05.000"),
	}
}

// // get clinet IP accresss
// func GetClientIP(r *http.Request) (string, error) {
// 	//Get IP from the X-REAL-IP header
// 	ip := r.Header.Get("X-REAL-IP")
// 	netIP := net.ParseIP(ip)
// 	if netIP != nil {
// 		return ip, nil
// 	}

// 	//Get IP from X-FORWARDED-FOR header
// 	ips := r.Header.Get("X-FORWARDED-FOR")
// 	splitIps := strings.Split(ips, ",")
// 	for _, ip := range splitIps {
// 		netIP := net.ParseIP(ip)
// 		if netIP != nil {
// 			return ip, nil
// 		}
// 	}

// 	//Get IP from RemoteAddr
// 	ip, _, err := net.SplitHostPort(r.RemoteAddr)
// 	if err != nil {
// 		return "", err
// 	}
// 	netIP = net.ParseIP(ip)
// 	if netIP != nil {
// 		return ip, nil
// 	}
// 	return "", fmt.Errorf("no valid ip found")
// }

func getSQLString(str string) pgtype.Text {
	return pgtype.Text{String: str, Valid: true}
}
func GetSQLInt(val int) sql.NullInt32 {
	return sql.NullInt32{Int32: int32(val), Valid: true}
}
func GetInt(val pgtype.Int4, defaultValue int) int {
	if val.Valid {
		return int(val.Int32)
	}
	return defaultValue
}
func ConvertInt32ToPgInt4(val int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: val,
		Valid: true,
	}
}
func getSQLDate(t *time.Time) pgtype.Date {
	var d pgtype.Date
	if t != nil && !t.IsZero() {
		_ = d.Scan(*t)
	}
	return d
}

func fromSQLDate(d pgtype.Date) *time.Time {
	if d.Valid {
		t := d.Time
		return &t
	}
	return nil
}

func StringToPgDate(dateStr string) pgtype.Date {
	var d pgtype.Date
	if dateStr != "" {
		_ = d.Scan(dateStr) // Handles parsing and nullability
	}
	return d
}

func ToPGTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: true,
	}
}
func ToPGTimestampPtr(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{
		Time:  *t,
		Valid: true,
	}
}

// convert json/interface to byte
func ToJSONBytes(data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data to JSON: %w", err)
	}
	return bytes, nil
}
func ByteToInterface(data []byte) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal byte to interface: %w", err)
	}
	return result, nil
}

func GetString(text pgtype.Text) string {
	if text.Valid {
		return text.String
	}
	return ""
}

// func GenerateHMAC(secret string) string {
// 	location, _ := time.LoadLocation("Asia/Dhaka")
// 	bdTime := time.Now().In(location)
// 	nonce := strconv.Itoa(bdTime.Day()) + bdTime.Month().String() + strconv.Itoa(bdTime.Year())
// 	// fmt.Println("Generated Nonce:", nonce) -> nonce will be like 5March2024
// 	h := hmac.New(sha256.New, []byte(secret))
// 	h.Write([]byte(nonce))
// 	return hex.EncodeToString(h.Sum(nil))
// }

func GetHashOf(password string) string {
	shaBytes := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", shaBytes)
}

func GetExpiryDate(days int) string {
	return time.Now().AddDate(0, 0, days).Format("20060102")
}

func removeSpacesAndSpecialChars(input string) string {
	// Regular expression to match only letters and digits
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	// Replace all characters NOT letters or digits with empty string
	return reg.ReplaceAllString(input, "")
}

func BuildDataMap[T any](objs []T, keySelector func(T) string) map[string][]T {
	result := make(map[string][]T)
	for _, obj := range objs {
		key := keySelector(obj)
		result[key] = append(result[key], obj)
	}
	return result
}

type EmailRequest struct {
	Username  string `json:"username"`
	UserEmail string `json:"user_email"`
	Otp       string `json:"otp"`
	NewUser   bool   `json:"new_user"`
}

func request3rdPartyEmailAPI(email, name, otp string, isNew bool) bool {
	url := "http://118.67.213.45:8339/server/api/email-service/send-email"

	// Build the request model
	payload := EmailRequest{
		Username:  name,
		UserEmail: email,
		Otp:       otp,
		NewUser:   isNew,
	}

	// Serialize to JSON
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return false
	}

	// Create POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return false
	}

	// Headers
	systemKey := viper.GetViper().GetString("system_key")
	apiKey := viper.GetViper().GetString("api_key")

	req.Header.Add("X-SYSTEM-KEY", systemKey)
	req.Header.Add("X-API-KEY", apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
