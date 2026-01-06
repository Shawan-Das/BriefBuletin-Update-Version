package util

import (
	"fmt"
	"io"
	"math/rand"
	"regexp"

	a "crypto/rand"
	"strconv"
	"strings"
	"time"
	_ "unicode"

	"github.com/jackc/pgx/v5/pgtype"
)

func GenerateSixDigits() string {
	rand.NewSource(time.Now().UnixNano())
	// rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	otpInt := rand.Intn(max-min+1) + min
	otp := strconv.Itoa(otpInt)

	return otp
}

// GetFinancialYearFromDateTime returns financial year from input date
func GetFinancialYearFromDateTime(date *time.Time) string {

	y, m, d := date.Date()
	if m >= time.July && d >= 1 {
		return fmt.Sprintf("%d-%d", y, y+1)
	}
	return fmt.Sprintf("%d-%d", y-1, y)

}

func GetFinancialYearFromDate(date string) (string, error) {
	dt, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", fmt.Errorf("invalid date provided")
	}
	y, m, d := dt.Date()
	if m >= time.July && d >= 1 {
		return fmt.Sprintf("%d-%d", y, y+1), nil
	}
	return fmt.Sprintf("%d-%d", y-1, y), nil
}

func GenerateUniqueID(sec bool, microSec bool, miliSec bool) string {
	// Generate a unique ID using the current timestamp and a random number
	if sec {
		return fmt.Sprint(time.Now().Unix())
	} else if miliSec {
		return fmt.Sprint(time.Now().UnixMilli())
	} else if microSec {
		return fmt.Sprint(time.Now().UnixMicro())
	}
	return fmt.Sprint(time.Now().UnixNano())
}

func MustStringToPgDate(dateStr string) pgtype.Date {
	t, _ := time.Parse("2006-01-02", dateStr)

	var pgDate pgtype.Date
	_ = pgDate.Scan(t)

	return pgDate
}

func MustStringToPgTimestamp(dateStr string) pgtype.Timestamp {
	t, _ := time.Parse("2006-01-02 15:04:05", dateStr)

	var pgTimestamp pgtype.Timestamp
	_ = pgTimestamp.Scan(t)

	return pgTimestamp
}

func TrimString(str string) string {
	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(str, " "))
}

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func IsValidName(name string) bool {
	name = TrimString(name)
	// Regex: Starts with a letter, allows letters, spaces, apostrophes, and hyphens
	// Length: 2-50 characters
	regex := `^[A-Za-z][A-Za-z\s'\-\.]{1,500}$`
	re := regexp.MustCompile(regex)

	return re.MatchString(name)
}

func GetCurrentBDTime() time.Time {
	loc, err := time.LoadLocation("Asia/Dhaka")
	if err != nil {
		// fallback to UTC+6 if timezone file is missing
		return time.Now().UTC().Add(6 * time.Hour)
	}
	return time.Now().In(loc)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(a.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}