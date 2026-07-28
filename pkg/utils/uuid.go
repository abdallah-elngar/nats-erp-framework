package utils

import (
	"crypto/rand"
	"fmt"
)

// UUID يمثل معرفاً فريداً عالمياً
type UUID string

// NewUUID يولد معرفاً فريداً جديداً
func NewUUID() (UUID, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// تعيين الإصدار والمتغير
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4],
		bytes[4:6],
		bytes[6:8],
		bytes[8:10],
		bytes[10:16],
	)

	return UUID(uuid), nil
}

// MustUUID يولد معرفاً فريداً أو يسبب panic
func MustUUID() UUID {
	uuid, err := NewUUID()
	if err != nil {
		panic(err)
	}
	return uuid
}

// String يعيد النص
func (u UUID) String() string {
	return string(u)
}

// IsEmpty يتحقق من أن المعرف فارغ
func (u UUID) IsEmpty() bool {
	return u == ""
}
