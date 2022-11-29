package svc

import (
	"github.com/cuczhangyi/coros-strava/internal/config"
)


var SCtx  *ServiceContext 

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {

	SCtx  =&ServiceContext{
		Config: c,
	}
	return  SCtx  
}
