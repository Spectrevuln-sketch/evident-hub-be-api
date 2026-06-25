package envs

import (
	"github.com/kelseyhightower/envconfig"
)

type appConfig struct {
	Host       string `envconfig:"HOST"`
	Port       string `envconfig:"PORT"`
	UploadPath string `envconfig:"UPLOAD_PATH"`

	LogoUrl string `envconfig:"LOGO_URL"`

	FileKey       string `envconfig:"FILE_KEY"`
	UsernameAdmin string `envconfig:"USERNAME_ADMIN"`
	EmailAdmin    string `envconfig:"EMAIL_ADMIN"`
	PassAdmin     string `envconfig:"PASS_ADMIN"`

	JwtSecret    string `envconfig:"JWT_SECRET"`
	JwtExpiredIn string `envconfig:"JWT_EXPIRED_IN"`
	JwtMaxAge    int    `envconfig:"JWT_MAXAGE"`

	PusdatinBasicAuthUsername string `envconfig:"PUSDATIN_BASIC_AUTH_USERNAME" default:"pusdatin"`
	PusdatinBasicAuthPassword string `envconfig:"PUSDATIN_BASIC_AUTH_PASSWORD" default:"pusdatin"`
}

func newAppConfig() *appConfig {
	var appCfg appConfig
	envconfig.MustProcess("", &appCfg)
	return &appCfg
}
