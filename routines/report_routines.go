package routines

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/models"
)

func LogReport(report *models.Report) {
	logFilePath := "log/report.log"
	absLogFilePath, err := filepath.Abs(logFilePath)
	if err != nil {
		helpers.LogServerError("Error while logging a report-LogReport", err, "go_routine")
		return
	}

	// Open the log file in append mode or create it if it doesn't exist
	logFile, err := os.OpenFile(absLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		helpers.LogServerError("Error while logging a report-LogReport", err, "go_routine")
		return
	}
	defer logFile.Close()

	// Create a logger with timestamps in ISO 8601 format and a custom logging format
	logger := log.New(logFile, "", 0) // No predefined flags, as we'll format the timestamp ourselves

	logger.Printf("Timestamp: %s, Report ID: %s, Reported by: %s, Report Type: %d\n", time.Now().UTC().Format(time.RFC3339), report.ID, report.ReporterID, report.ReportType)
}
