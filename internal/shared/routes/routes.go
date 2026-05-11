package routes

import (
	"fmt"
	"net/url"
)

const (
	Static = "/static/*"
)

const (
	Index        = "/"
	Login        = "/login"
	ILoginGoogle = "/i/login/google"
	ILogout      = "/i/logout"
)

var ILoginGoogleCallback = ""

func Load(googleLoginCallbackURL string) error {
	url, err := url.Parse(googleLoginCallbackURL)
	if err != nil {
		return fmt.Errorf("parsing google login callback url: %w", err)
	}
	ILoginGoogleCallback = url.Path
	return nil
}
