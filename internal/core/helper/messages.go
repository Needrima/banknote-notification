package helper

import (
	"fmt"
	"time"
)

var (
	NoResourceFound = "this resource does not exist"
	NoRecordFound   = "sorry. no record found"
	//NoErrorsFound   = "no errors at the moment"
	ServerStarted = fmt.Sprintf("Started " + Config.ServiceName + " on " + Config.ServiceAddress + ":" +
		Config.ServicePort + " in " + time.Since(time.Now()).String())
)
