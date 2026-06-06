package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// GenerateTicketCode สร้างรหัสตั๋วหน้างานในรูปแบบ: TK-ปีเดือนวัน-รหัสแคชเชียร์-อักษรสุ่ม4หลัก
func GenerateTicketCode(cashierID uint) string {
	now := time.Now()
	dateStr := now.Format("060102") // ได้ฟอร์แมต YYMMDD (เช่น 260524)
	randomStr := generateSecureRandomString(4)
	
	return fmt.Sprintf("TK-%s-%d-%s", dateStr, cashierID, randomStr)
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