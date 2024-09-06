package applog

import (
	"os"
    "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var Logger log.Logger

func ApplogInit(logLevel string){
	Logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	Logger = level.NewFilter(Logger,level.Allow(level.ParseDefault(logLevel, level.DebugValue())))
    Logger = log.With(Logger,
        "time", log.DefaultTimestampUTC,
        "caller", log.DefaultCaller,
    )
}
