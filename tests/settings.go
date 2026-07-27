package tests

import (
	"proj/pkg/api/service"
	"strconv"
)

var Port, _ = (strconv.Atoi(service.PortGlobal()))
var DBFile = service.PathDbManual()
var FullNextDate = false
var Search = false
var Token = ``
