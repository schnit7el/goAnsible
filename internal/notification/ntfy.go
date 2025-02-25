package notification

import (
	"net/http"
	"strings"
)

type NtfyController struct {
	Auth string
	Url  string
}

func NewNtfyController(ntfy_auth string, ntfy_url string) *NtfyController {
	return &NtfyController{
		Auth: ntfy_auth,
		Url:  ntfy_url,
	}
}

func (c *NtfyController) SendNtfy(title string, content string) (string, error) {
	// godotenv.Load()
	// auth := os.Getenv("NTFY_AUTH")
	// url := os.Getenv("NTFY_URL")

	req, _ := http.NewRequest(
		"POST",
		c.Url,
		strings.NewReader(content),
	)

	req.Header.Set("Title", title)
	req.Header.Set("Authorization", c.Auth)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "nil", err
	}

	return res.Status, nil
}
