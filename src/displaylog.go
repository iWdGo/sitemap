package src

// Reading log content
import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
	"time"
)

// Logentry keeps some data for each log version
type logentry struct {
	version  string
	failures int
	total    int
}

func displaylog(r *http.Request) (logentries []logentry) {

	c := appengine.NewContext(r)
	query := &log.Query{
		AppLogs: true,
	}

	f, i := 0, 0
	var sinceT time.Time
	versionL := "0" // Log version
	for results := query.Run(c); ; {
		record, err := results.Next()
		if err == log.Done {
			break
		}
		if err != nil {
			log.Errorf(c, "Failed to retrieve next record: %v", err)
			break
		}
		if record.Status != 200 {
			f++
		}
		// Some key values
		sinceT = record.StartTime // last record provides the start time
		if versionL != record.VersionID {
			if versionL != "0" {
				// log.Infof(c, "log version %s has %d failures on %d", versionL, f, i)
				logentries = append(logentries, logentry{versionL, f, i})
			}
			versionL = record.VersionID
			f, i = 0, 0
			// When run in flex with dev_appserver.py, no instance id is provided and log has NO-VERSION
		}
		i++
	}
	log.Infof(c, "log since %v", sinceT)
	return
}
