package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/golang-jwt/jwt/v5"
)

func createMailerJWT() (string, error) {
	token_claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "backend",
		"crt": time.Now().Unix(),
		"exp": time.Now().Add(config.ACCESS_TOKEN_TTL).Unix(),
	})

	token, err := token_claim.SignedString([]byte(initializers.CONFIG.MAILER_SECRET))
	if err != nil {
		return "", err
	}

	return token, nil
}

func SendMailReq(email string, emailType int, user *models.User, otp *string, secondaryUser *models.User) error {
	jsonData, err := json.Marshal(map[string]any{
		"email":         email,
		"type":          emailType,
		"user":          user,
		"otp":           otp,
		"secondaryUser": secondaryUser,
	})
	if err != nil {
		initializers.Logger.Errorw("Error calling Mailer", "Message", err.Error(), "Path", "SendMailReq", "Error", err.Error())
		return err
	}

	jwt, err := createMailerJWT()
	if err != nil {
		initializers.Logger.Errorw("Error calling Mailer", "Message", err.Error(), "Path", "SendMailReq", "Error", err.Error())
		return err
	}

	request, err := http.NewRequest("POST", initializers.CONFIG.MAILER_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		initializers.Logger.Errorw("Error calling Mailer", "Message", err.Error(), "Path", "SendMailReq", "Error", err.Error())
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+jwt)
	request.Header.Set("api-token", initializers.CONFIG.MAILER_TOKEN)

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		initializers.Logger.Errorw("Error calling Mailer", "Message", err.Error(), "Path", "SendMailReq", "Error", err.Error())
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		initializers.Logger.Errorw("Error calling Mailer", "Message", fmt.Sprint("Status Code - ", response.StatusCode), "Path", "SendMailReq", "Error", response.Body)
		return fmt.Errorf("error calling mailer")
	}

	return nil
}
