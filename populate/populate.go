// package main
package populate

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/google/uuid"
)

// func init() {
// 	initializers.LoadEnv()
// 	initializers.ConnectToDB()
// }

func Slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces with -
	s = strings.ReplaceAll(s, " ", "-")

	// Remove non-word characters except -
	reg := regexp.MustCompile("[^a-zA-Z0-9-]")
	s = reg.ReplaceAllString(s, "")

	// Replace multiple - with single -
	reg = regexp.MustCompile("-+")
	s = reg.ReplaceAllString(s, "-")

	// Remove leading and trailing -
	s = strings.Trim(s, "-")

	return s
}

func ToLowercaseArray(arr []string) []string {
	result := make([]string, len(arr))

	for i, str := range arr {
		result[i] = strings.ToLower(str)
	}

	return result
}

func RandomLinks() []string {
	strings := []string{"https://www.google.com", "https://www.youtube.com", "https://www.facebook.com", "https://www.gmail.com", "https://www.github.com"}

	// Get a random count between 0 and 5
	count := rand.Intn(6)

	rand.Shuffle(len(strings), func(i, j int) { strings[i], strings[j] = strings[j], strings[i] })

	return strings[:count]
}

func getRandomUserID(userIDs []uuid.UUID) uuid.UUID {
	return userIDs[rand.Intn(len(userIDs))]
}

func getRandomProjectID(projectIDs []uuid.UUID) uuid.UUID {
	return projectIDs[rand.Intn(len(projectIDs))]
}

func PopulateProjects() {
	log.Println("----------------Populating Projects----------------")

	jsonFile, err := os.Open("populate/projects.json")
	if err != nil {
		log.Fatalf("Failed to open the JSON file: %v", err)
	}
	defer jsonFile.Close()

	var projects []models.Project
	jsonDecoder := json.NewDecoder(jsonFile)
	if err := jsonDecoder.Decode(&projects); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	var users []models.User
	initializers.DB.Find(&users)

	var userIDs []uuid.UUID
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	for _, project := range projects {
		project.UserID = getRandomUserID(userIDs)
		project.Slug = Slugify(project.Title)
		project.Tags = ToLowercaseArray(project.Tags)
		project.Links = RandomLinks()

		if err := initializers.DB.Create(&project).Error; err != nil {
			log.Printf("Failed to insert project: %v", err)
		} else {
			log.Printf("Added Project: %s", project.Title)
		}
	}
}

func PopulatePosts() {
	log.Println("----------------Populating Posts----------------")

	jsonFile, err := os.Open("populate/posts.json")
	if err != nil {
		log.Fatalf("Failed to open the JSON file: %v", err)
	}
	defer jsonFile.Close()

	var posts []models.Post
	jsonDecoder := json.NewDecoder(jsonFile)
	if err := jsonDecoder.Decode(&posts); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	var users []models.User
	initializers.DB.Find(&users)

	var userIDs []uuid.UUID
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	for _, post := range posts {
		post.UserID = getRandomUserID(userIDs)

		if err := initializers.DB.Create(&post).Error; err != nil {
			log.Printf("Failed to insert post: %v", err)
		}
	}
}

func PopulateOpenings() {
	log.Println("----------------Populating Openings----------------")

	jsonFile, err := os.Open("populate/openings.json")
	if err != nil {
		log.Fatalf("Failed to open the JSON file: %v", err)
	}
	defer jsonFile.Close()

	var openings []models.Opening
	jsonDecoder := json.NewDecoder(jsonFile)
	if err := jsonDecoder.Decode(&openings); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	var projects []models.Project
	initializers.DB.Find(&projects)

	var projectIDs []uuid.UUID
	for _, project := range projects {
		projectIDs = append(projectIDs, project.ID)
	}

	for _, opening := range openings {
		opening.ProjectID = getRandomProjectID(projectIDs)

		var project models.Project
		initializers.DB.First(&project, "id=?", opening.ProjectID)

		opening.UserID = project.UserID

		if err := initializers.DB.Create(&opening).Error; err != nil {
			log.Printf("Failed to insert opening: %v", err)
		} else {
			log.Printf("Added Opening: %s, in Project %s", opening.Title, project.Title)
		}
	}
}

// func main() {
// 	FillDummies()
// }

func FillDummies() {
	PopulateProjects()
	PopulatePosts()
	PopulateOpenings()
}
