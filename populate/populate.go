// package main
package populate

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/Pratham-Mishra04/interact/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// func init() {
// 	initializers.LoadEnv()
// 	initializers.ConnectToDB()
// }

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
	if err := initializers.DB.Find(&users).Error; err != nil {
		return
	} else {
		if len(users) == 0 {
			return
		}
	}

	var userIDs []uuid.UUID
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	coverPics := []string{"default_1.jpg", "default_2.jpg", "default_3.jpg", "default_4.jpg", "default_5.jpg", "default_6.jpg", "default_7.jpg", "default_8.jpg", "default_9.jpg"}

	for _, project := range projects {
		project.UserID = getRandomUserID(userIDs)
		project.Slug = utils.SoftSlugify(project.Title)
		project.Tags = ToLowercaseArray(project.Tags)
		project.Links = RandomLinks()

		// Randomly select a cover picture
		coverPicIndex := rand.Intn(len(coverPics))
		project.CoverPic = coverPics[coverPicIndex]

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
	if err := initializers.DB.Find(&users).Error; err != nil {
		return
	} else {
		if len(users) == 0 {
			return
		}
	}

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
	if err := initializers.DB.Find(&projects).Error; err != nil {
		return
	} else {
		if len(projects) == 0 {
			return
		}
	}

	var projectIDs []uuid.UUID
	for _, project := range projects {
		projectIDs = append(projectIDs, project.ID)
	}

	for _, opening := range openings {
		projectID := getRandomProjectID(projectIDs)
		opening.ProjectID = &projectID

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

func PopulateColleges() {
	log.Println("----------------Populating Colleges----------------")

	jsonFile, err := os.Open("populate/colleges.json")
	if err != nil {
		log.Fatalf("Failed to open the JSON file: %v", err)
	}
	defer jsonFile.Close()

	var colleges []models.College
	jsonDecoder := json.NewDecoder(jsonFile)
	if err := jsonDecoder.Decode(&colleges); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	for _, college := range colleges {
		if err := initializers.DB.Create(&college).Error; err != nil {
			log.Printf("Failed to insert college: %v", err)
		} else {
			log.Printf("Insert college: %s", college.Name)
		}
	}
}

func PopulateUsersAndOrgs() {
	log.Println("----------------Populating Organisations----------------")

	jsonFile, err := os.Open("populate/users.json")
	if err != nil {
		log.Fatalf("Failed to open the JSON file: %v", err)
	}
	defer jsonFile.Close()

	type User struct {
		Name               string         `json:"name"`
		Username           string         `json:"username"`
		Email              string         `json:"email"`
		Password           string         `json:"password"`
		Bio                string         `json:"bio"`
		Title              string         `json:"title"`
		Tagline            string         `json:"tagline"`
		Tags               pq.StringArray `json:"tags"`
		Links              pq.StringArray `json:"links"`
		OrganizationStatus bool           `json:"isOrganization"`
	}

	var users []User
	jsonDecoder := json.NewDecoder(jsonFile)
	if err := jsonDecoder.Decode(&users); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	for _, user := range users {
		log.Println("\nCreating User - " + user.Name)

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			log.Println("Error while hashing Password.", err)
			continue
		}

		newUser := models.User{
			Name:                user.Name,
			Email:               user.Email,
			Password:            string(hash),
			Username:            user.Username,
			Tagline:             user.Tagline,
			PasswordChangedAt:   time.Now(),
			OrganizationStatus:  user.OrganizationStatus,
			Tags:                user.Tags,
			Links:               user.Links,
			Verified:            true,
			OnboardingCompleted: true,
		}

		result := initializers.DB.Create(&newUser)
		if result.Error != nil {
			log.Println("Error while creating Org User.", result.Error)
			continue
		}

		if newUser.OrganizationStatus {
			log.Println("\nCreating Org - " + user.Name)
			organization := models.Organization{
				UserID:            newUser.ID,
				OrganizationTitle: newUser.Name,
				CreatedAt:         time.Now(),
			}

			result = initializers.DB.Create(&organization)
			if result.Error != nil {
				log.Println("Error while creating Org.", result.Error)
				continue

			}
		}

		log.Println("\nCreating Profile - " + user.Name)

		newProfile := models.Profile{
			UserID: newUser.ID,
		}

		result = initializers.DB.Create(&newProfile)
		if result.Error != nil {
			log.Println("Error while creating Profile.", result.Error)
			continue
		}

		log.Println("Successfully created User - " + newUser.Name)
	}
}

func PopulateOrgs() {
	log.Println("----------------Populating Organisations----------------")

	jsonFile, err := os.Open("scripts/organisations.json")
	if err != nil {
		log.Fatalf("Failed to open the JSON file: %v", err)
	}
	defer jsonFile.Close()

	type User struct {
		Name     string `json:"name"`
		Username string `son:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Tagline  string `json:"tagline"`
	}

	var users []User
	jsonDecoder := json.NewDecoder(jsonFile)
	if err := jsonDecoder.Decode(&users); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	for _, user := range users {
		log.Println("\nCreating Org - " + user.Name)

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			log.Println("Error while hashing Password.", err)
			continue
		}

		newOrg := models.User{
			Name:                user.Name,
			Email:               user.Email,
			Password:            string(hash),
			Username:            user.Username,
			Tagline:             user.Tagline,
			PasswordChangedAt:   time.Now(),
			OrganizationStatus:  true,
			Verified:            true,
			OnboardingCompleted: true,
		}

		result := initializers.DB.Create(&newOrg)
		if result.Error != nil {
			log.Println("Error while creating Org User.", result.Error)
			continue
		}

		organization := models.Organization{
			UserID:            newOrg.ID,
			OrganizationTitle: newOrg.Name,
			CreatedAt:         time.Now(),
		}

		result = initializers.DB.Create(&organization)
		if result.Error != nil {
			log.Println("Error while creating Org.", result.Error)
			continue

		}

		newProfile := models.Profile{
			UserID: newOrg.ID,
		}

		result = initializers.DB.Create(&newProfile)
		if result.Error != nil {
			log.Println("Error while creating Org User Profile.", result.Error)
			continue
		}

		log.Println("Successfully created Org - " + newOrg.Name)
	}
}

func FillDummies() {
	PopulateProjects()
	PopulatePosts()
	PopulateOpenings()
}
