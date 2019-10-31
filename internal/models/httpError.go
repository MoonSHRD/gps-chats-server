package models

type HTTPError struct {
	Ok      bool   `json:"ok"`
	ErrCode int    `json:"errCode"`
	ErrText string `json:"errText"`
}
