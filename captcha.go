package d

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Captcha interface, implement at least the following methods to facilitate internal calls in the devtool library
type InterfaceCaptcha interface {
	Init()
	VerifyToken(token string) error
}

const (
	ConfigPathCaptchaSecret = "captcha.secret"
	ConfigPathCaptchaUrl    = "captcha.url"
)

var (
	captcha InterfaceCaptcha // Global variable, stores the initialized interface, if not initialized, it is nil
)

// Captcha library unified access entry
type Captcha[T InterfaceCaptcha] struct{}

// Initialization
func (c Captcha[T]) Init(conf T) {
	captcha = conf
}

// Get the initialized interface. If it is not initialized, Turnstile library is used by default.
func (c Captcha[T]) Get() T {
	if captcha == nil {
		LibraryTurnstile{}.Init()
	}
	return captcha.(T)
}

// the variable of Turnstile library
var (
	ErrCaptchaTurnstileEmptySecret = errors.New("the secret of Turnstile cannot be empty")
)

// Turnstile library
type LibraryTurnstile struct {
	Secret string // Requried, the secret key of Turnstile
	Url    string
}

// Initialization
func (l LibraryTurnstile) Init() {
	Captcha[LibraryTurnstile]{}.Init(LibraryTurnstile{
		Secret: Config[InterfaceConfig]{}.Get().GetString(ConfigPathCaptchaSecret),
		Url:    Config[InterfaceConfig]{}.Get().GetString(ConfigPathCaptchaUrl),
	})
}

// https://developers.cloudflare.com/turnstile/get-started/server-side-validation/
// curl 'https://challenges.cloudflare.com/turnstile/v0/siteverify' --data 'secret=verysecret&response=<RESPONSE>'
func (t LibraryTurnstile) VerifyToken(token string) error {
	// Check if secret is empty
	capt := Captcha[LibraryTurnstile]{}.Get()
	if capt.Secret == "" {
		return ErrCaptchaTurnstileEmptySecret
	}
	// Set defaut URL
	if t.Url == "" {
		t.Url = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	}

	data := url.Values{}
	data.Set("secret", t.Secret)
	data.Set("response", token)

	resp, err := http.PostForm(t.Url, data)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if success, ok := result["success"].(bool); ok {
		if !success {
			errorCodes, ok := result["error-codes"].([]interface{})
			if !ok {
				return fmt.Errorf("%v", result["error-codes"])
			}

			var errorCodesStr string
			for _, code := range errorCodes {
				errorCodesStr += fmt.Sprintf("%s, ", code)
			}
			// Remove the last comma and space
			errorCodesStr = errorCodesStr[:len(errorCodesStr)-2]

			return errors.New(errorCodesStr)
		}
	} else {
		return errors.New("invalid response format")
	}

	return nil
}
