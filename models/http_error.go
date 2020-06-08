package models

type HTTPError struct {
	Error HTTPErrorMessage `json:"error"`
}

type HTTPErrorMessage struct {
	Message string `json:"message"`
}
