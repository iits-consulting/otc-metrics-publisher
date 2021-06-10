package main

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/rs/zerolog/log"
)

type config struct {
	// cloud endpoint settings
	// Endpoint  string `default:""`
	AuthEndpoint string `default:"https://iam.eu-de.otc.t-systems.com/v3"`
	ProjectID    string `default:"" split_words:"true"`
	// TrustSSL  bool   `default:"true"`

	// auth
	AuthMethod string `default:"aksk"`
	AccessKey  string `default:"" split_words:"true"`
	SecretKey  string `default:"" split_words:"true"`
	User       string `default:""`
	Password   string `default:""`

	// metrics settings
	NameSpace               string `default:"APP.nextcloud"`
	SendInterval            int    `default:"60" split_words:"true"`
	GrabInterval            int    `default:"10" split_words:"true"`
	ScriptsDir              string `default:"/opt/metric-scripts" split_words:"true"`
	InstanceID              string `default:"" split_words:"true"`
	FileCloudInitInstanceID string `default:"/run/cloud-init/.instance-id" split_words:"true"` // we can read instance ID from that file (if not defined)
	TTL                     int    `default:"86400"`                                           // one day

	// log settings
	LogLevel      string `default:"error" split_words:"true"`
	LogFormat     string `default:"json" split_words:"true"`
	LogTimeFormat string `default:"unix" split_words:"true"`
}

func configHandler() (err error) {

	// chose between AK/SK or userpass auth methods. The default one is AK/SK.
	if len(cfg.AccessKey) == 0 &&
		len(cfg.SecretKey) == 0 {
		if len(cfg.User) > 0 && len(cfg.Password) > 0 {
			cfg.AuthMethod = "userpass"
		} else {
			err = errors.New("Either AK/SK or User and Password should be defined to access API; " +
				"Env vars for AK/SK: " + strings.ToUpper(envVarPrefix+"_Access_Key") +
				" and " + strings.ToUpper(envVarPrefix+"_Secret_Key") + "; " +
				"Env vars for user/pass: " + strings.ToUpper(envVarPrefix+"_User") +
				" and " + strings.ToUpper(envVarPrefix+"_Password"))
			return
		}
	}

	// check instance ID
	if len(cfg.InstanceID) == 0 {
		// default value?
		cfg.InstanceID = "undefined"

		// try to get ID form cloud-init
		if binData, e := ioutil.ReadFile(cfg.FileCloudInitInstanceID); e == nil {
			cfg.InstanceID = string(binData)
		}
	}

	// check NameSpace (PREFIX.name)
	ns := strings.Split(cfg.NameSpace, ".")
	if len(ns) != 2 {
		err = errors.New("wrong namespace format '" + cfg.NameSpace + "': The value must be in the service.item format and can contain 3 to 32 characters. service and item must be a string that starts with a letter and contains only uppercase letters, lowercase letters, digits, and underscores (_). In addition, service cannot start with SYS and AGT.")
		return
	}
	if strings.EqualFold(ns[0], "SYS") == true ||
		strings.EqualFold(ns[0], "AGT") {
		err = errors.New("service cannot start with SYS and AGT")
		return
	}

	log.Trace().Msgf("config: %#v", cfg)

	return
}
