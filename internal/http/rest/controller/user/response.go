package user

type response struct {
	Message string `json:"message,omitempty"`
	QRCode  string `json:"qrcode,omitempty"`
	Status  string `json:"status,omitempty"`
}
