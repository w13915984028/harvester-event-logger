package controller

import (
	"github.com/w13915984028/harvester-event-logger/pkg/config"
	"github.com/w13915984028/harvester-event-logger/pkg/controller/eventlogger"
)

var RegisterFuncList = []config.RegisterFunc{
	eventlogger.Register,
}
