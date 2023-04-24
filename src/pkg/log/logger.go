package log

import (
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func init() {
	// TODO:realize log rotate and write to file
	log.SetOutput(io.MultiWriter(os.Stdout))
	log.SetReportCaller(true)

}
