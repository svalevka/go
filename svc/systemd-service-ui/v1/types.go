package v1

import "net/http"

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) StatusCode() int {
	switch e.Code {
	case "NotFound":
		return http.StatusNotFound

	case "MethodNotAllowed":
		return http.StatusMethodNotAllowed

	default:
		return http.StatusInternalServerError
	}
}

func (e *Error) Error() string {
	return "Error: " + e.Code + ": " + e.Message
}

type ErrorWrapper struct {
	*Error `json:"error"`
}

type Service struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Running     bool   `json:"running"`
}

// Services is an array of Service types that is sortable by name.
type Services []*Service

func (s Services) Len() int           { return len(s) }
func (s Services) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s Services) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type ListServicesRes struct {
	Hostname string   `json:"hostname"`
	Services Services `json:"services"`
}
