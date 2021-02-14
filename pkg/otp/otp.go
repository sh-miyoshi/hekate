package otp

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/util"
	qrcode "github.com/skip2/go-qrcode"
)

const (
	registerExpiresIn = 10 * time.Minute
	period            = 30
	digitLen          = 6
)

// Register ...
func Register(projectName string, userName string) (string, *errors.Error) {
	// generate private key
	data := model.OTP{
		ID:              uuid.New().String(),
		ProjectName:     projectName,
		PrivateKey:      util.RandomString(32, util.CharTypeDigit|util.CharTypeUpper),
		InitExpiresDate: time.Now().Add(registerExpiresIn),
	}
	logger.Debug("set OTP data: %v", data)

	// enter to db
	if err := db.GetInst().OTPAdd(projectName, &data); err != nil {
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
