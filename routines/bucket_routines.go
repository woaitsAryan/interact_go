package routines

import (
	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
)

func DeleteFromBucket(client *helpers.BucketClient, path string) {
	if _, found := config.AcceptedDefaultProjectHashes[path]; path == "" || path == "default.jpg" || found {
		return
	}
	err := client.DeleteBucketFile(path)
	if err != nil {
		initializers.Logger.Warnw("Error while deleting file from bucket", "Error", err)
	}
}
