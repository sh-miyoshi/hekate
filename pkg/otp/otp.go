package otp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	qrcode "github.com/skip2/go-qrcode"
)

const (
	period   = 30
	digitLen = 6
)

var (
	// ErrNotEnabled ...
	ErrNotEnabled = errors.New("Authenticator Application is not set", "User OTP Info is not enabled")
	// ErrVerifyFailed ...
	ErrVerifyFailed = errors.New("Invalid User Code", "Invalid User Code")
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
func Verify(now time.Time, projectName, userID, userCode string) *errors.Error {
	user, err := db.GetInst().UserGet(projectName, userID)
	if err != nil {
		return errors.Append(err, "Failed to get user OTP info")
	}

	if !user.OTPInfo.Enabled {
		return ErrNotEnabled
	}

	keySrc, e := base32.StdEncoding.DecodeString(user.OTPInfo.PrivateKey)
	if e != nil {
		return errors.New("Internal Server Error", "Failed to decode private key %v", e)
	}

	tSrc := zeroPadding(strconv.FormatInt(now.Unix()/period, 16), 16)
	t := make([]byte, hex.DecodedLen(len(tSrc)))
	hex.Decode(t, []byte(tSrc))
	key := make([]byte, hex.DecodedLen(len(keySrc)))
	hex.Decode(key, keySrc)

	hs := getMAC(t, key)
	if len(hs) != 20 {
		return errors.New("Internal Server Error", "Invalid HMAC-SHA-1 size %d", len(hs))
	}

	expect := zeroPadding(strconv.Itoa(truncate(hs)), digitLen)
	if expect != userCode {
		logger.Debug("Failed to verify user code: expect %s but got %s", expect, userCode)
		return ErrVerifyFailed
	}

	return nil
}

func zeroPadding(d string, length int) string {
	for i := len(d); i < length; i++ {
		d = "0" + d
	}
	return d
}

func getMAC(message, key []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

func truncate(hs []byte) int {
	offset := hs[19] & 0xf
	binCode := (int(hs[offset])&0x7f)<<24 | (int(hs[offset+1])&0xff)<<16 | (int(hs[offset+2])&0xff)<<8 | (int(hs[offset+3]) & 0xff)
	return binCode % 1000000
}
