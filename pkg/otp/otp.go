package otp

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	qrcode "github.com/skip2/go-qrcode"
)

const (
	registerExpiresIn = 10 * time.Minute
	period            = 30
	digitLen          = 6
)

// Register ...
func Register(projectName string, userID, userName string) (string, *errors.Error) {
	// private key is 20 bytes and base32 encoded
	privateKey := make([]byte, 20)
	rand.Read(privateKey)

	data := model.OTPInfo{
		ID:         uuid.New().String(),
		PrivateKey: base32.StdEncoding.EncodeToString(privateKey),
		Enabled:    false,
	}
	logger.Debug("set OTP data: %v", data)

	// enter to db
	if err := db.GetInst().OTPAdd(projectName, userID, &data); err != nil {
		return "", errors.Append(err, "Failed to register OTP data")
	}

	// return qr code
	content := fmt.Sprintf("otpauth://totp/hekate:%s?secret=%s&digits=%d&issuer=hekate&period=%d", userName, data.PrivateKey, digitLen, period)
	var png []byte
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return "", errors.New("QR Code encoding failed", "Failed to QR encode: %v", err)
	}

	return base64.StdEncoding.EncodeToString(png), nil
}

// Verify ...
func Verify() {
	// calculate value
}
