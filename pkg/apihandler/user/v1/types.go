package userv1

// ChangePasswordRequest ...
type ChangePasswordRequest struct {
	Password string `json:"password"`
}

// OTPGenerateResponse ...
type OTPGenerateResponse struct {
	QRCodeImage string `json:"qrcode"` // base64 encorded png image data
}
