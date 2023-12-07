package helpers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/gofiber/fiber/v2"
)

type BucketClient struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

var ProjectClient *BucketClient
var EventClient *BucketClient
var PostClient *BucketClient
var ChatClient *BucketClient
var UserProfileClient *BucketClient
var UserCoverClient *BucketClient
var UserResumeBucket *BucketClient

func createNewBucketClient(uploadPath string) *BucketClient {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "bucket-key.json") // FILE PATH
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return &BucketClient{
		cl:         client,
		bucketName: initializers.CONFIG.GCP_BUCKET,
		projectID:  initializers.CONFIG.GCP_PROJECT,
		uploadPath: uploadPath,
	}
}

func InitializeBucketClients() {
	ProjectClient = createNewBucketClient("projects/")
	EventClient = createNewBucketClient("events/")
	PostClient = createNewBucketClient("posts/")
	ChatClient = createNewBucketClient("chats/")
	UserProfileClient = createNewBucketClient("users/profilePics/")
	UserCoverClient = createNewBucketClient("users/coverPics/")
	UserResumeBucket = createNewBucketClient("users/resumes/")
}

func (c *BucketClient) UploadBucketFile(buffer *bytes.Buffer, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	reader := bytes.NewReader(buffer.Bytes())

	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, reader); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

func (c *BucketClient) DeleteBucketFile(fileName string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	if err := c.cl.Bucket(c.bucketName).Object(c.uploadPath + fileName).Delete(ctx); err != nil {
		return &fiber.Error{Code: 500, Message: fmt.Sprintf("Failed to delete file: %v", err)}
	}

	return nil
}
