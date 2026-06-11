package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/nyaruka/phonenumbers"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func comparePassword(password string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

func generateFormattedThaiPhone(phone string) (string, error) {
	num, err := phonenumbers.Parse(phone, "TH")
	if err != nil {
		return "", err
	}
	return phonenumbers.Format(num, phonenumbers.NATIONAL), nil
}

// GenerateTicketCode สร้างรหัสตั๋วหน้างานในรูปแบบ: TK-ปีเดือนวัน-รหัสแคชเชียร์-อักษรสุ่ม n หลัก
func generateTicketCode(userID uint) string {
	now := time.Now()
	dateStr := now.Format("060102") // ได้ฟอร์แมต YYMMDD (เช่น 260524)
	randomStr := generateSecureRandomString(8)
	
	return fmt.Sprintf("TK-%s-%d-%s", dateStr, userID, randomStr)
}

func generateSecureRandomString(length int) string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" 
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))
	
	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, charsetLen)
		result[i] = charset[num.Int64()]
	}
	return string(result)
}