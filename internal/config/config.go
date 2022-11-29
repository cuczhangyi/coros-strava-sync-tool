package config

import "github.com/tal-tech/go-zero/rest"

type StravaCfg struct{
	ClientId int64    `json:"client_id,optional"`
	ClientSecret string  `json:"client_secret,optional"`
	StravaTempPort int64 `json:"strava_temp_port,optional"`
	CallBackUrl string `json:"call_back_url,optional"`
}

type Config struct {
	rest.RestConf
	StravaInfo StravaCfg 
}


