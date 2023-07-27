package helpers

import "github.com/Pratham-Mishra04/interact/initializers"

func LogDatabaseError(customString string, err error, path string) {
	initializers.Logger.Warnw(customString, "Message", err.Error(), "Path", path, "Error", err)
}

func LogServerError(customString string, err error, path string) {
	initializers.Logger.Errorw(customString, "Message", err.Error(), "Path", path, "Error", err)
}
