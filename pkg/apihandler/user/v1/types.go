package userv1

// OTPInfo ...
type OTPInfo struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

// GetResponse ...
type GetResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	EMail     string  `json:"email"`
	CreatedAt string  `json:"created_at"`
	OPTInfo   OTPInfo `json:"otp_info"`
}

// ChangePasswordRequest ...
type ChangePasswordRequest struct {
	Password string `json:"password"`
}

// OTPGenerateResponse ...
type OTPGenerateResponse struct {
	QRCodeImage string `json:"qrcode"` // base64 encorded png image data
}

// OTPVerifyRequest ...
type OTPVerifyRequest struct {
	UserCode string `json:"user_code"`
}
