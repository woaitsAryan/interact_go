package subscribers

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/Pratham-Mishra04/interact/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ImpressionsDumpSub(sub *redis.PubSub, db *gorm.DB) {
	// Listen for expired events
	if sub == nil {
		return
	}
	ch := sub.Channel() //TODO not listening for events

	for msg := range ch {
		modelName := extractModelDataFromKey(msg.Payload, 1)
		modelID := extractModelDataFromKey(msg.Payload, 2)

		if modelName == "" || modelID == "" {
			return
		}

		model := getModelFromStr(modelName)

		if err := db.Model(model).Where("id = ?", modelID).UpdateColumn("Impressions", gorm.Expr("Impressions + ?", rand.Intn(5)+3)).Error; err != nil {
			fmt.Printf("\n Error dumping impressions of key:%s, error:%e", msg.Payload, err)
		} else {
			fmt.Printf("\nKey %s dumped\n", msg.Payload)
		}
	}

	fmt.Println("Impressions Dump Listener added to Redis")
}

func extractModelDataFromKey(key string, index int) string {
	parts := strings.Split(key, "_")
	if len(parts) >= 3 {
		if parts[0] != "impressions" {
			return ""
		}
		return parts[index]
	}

	return ""
}

func getModelFromStr(modelName string) interface{} {
	switch modelName {
	case "post":
		return models.Post{}
	case "project":
		return models.Project{}
	case "event":
		return models.Event{}
	case "opening":
		return models.Opening{}
	case "user":
		return models.User{}
	}

	return nil
}
