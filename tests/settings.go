package tests

import (
	"proj/pkg/api/service"
	"strconv"
)

var Port, _ = (strconv.Atoi(service.PortGlobal()))
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = false
var Token = ``
