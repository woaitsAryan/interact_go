package auth_controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Pratham-Mishra04/interact/config"
	"github.com/Pratham-Mishra04/interact/helpers"
	"github.com/Pratham-Mishra04/interact/initializers"
	"github.com/Pratham-Mishra04/interact/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func RedirectToSignUp(c *fiber.Ctx, user models.User) error {
	sign_up_token_claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"crt": time.Now().Unix(),
		"exp": time.Now().Add(config.SIGN_UP_TOKEN_TTL).Unix(),
		"rdt": true,
	})

	sign_up_token, err := sign_up_token_claim.SignedString([]byte(initializers.CONFIG.JWT_SECRET))
	if err != nil {
		return &helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: "Error while decrypting JWT Token", Err: err}
	}

	return c.Redirect(initializers.CONFIG.FRONTEND_URL+"/signup/callback?token="+sign_up_token, fiber.StatusTemporaryRedirect)
}

func RedirectToLogin(c *fiber.Ctx, user models.User) error {
	login_token_claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"crt": time.Now().Unix(),
		"exp": time.Now().Add(config.LOGIN_TOKEN_TTL).Unix(),
		"rdt": true,
	})

	login_token, err := login_token_claim.SignedString([]byte(initializers.CONFIG.JWT_SECRET))
	if err != nil {
		return &helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: "Error while decrypting JWT Token", Err: err}
	}

	loginURL := "login"
	if user.OrganizationStatus {
		loginURL = "organisation/" + loginURL
	}

	return c.Redirect(initializers.CONFIG.FRONTEND_URL+"/"+loginURL+"/callback?token="+login_token, fiber.StatusTemporaryRedirect)
}

func GoogleRedirect(c *fiber.Ctx) error {
	URL, err := url.Parse(config.GoogleOAuthConfig.Endpoint.AuthURL)
	if err != nil {
		return &helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}
	parameters := url.Values{}
	parameters.Add("client_id", config.GoogleOAuthConfig.ClientID)
	parameters.Add("scope", strings.Join(config.GoogleOAuthConfig.Scopes, " "))
	parameters.Add("redirect_uri", config.GoogleOAuthConfig.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", config.GoogleOAuthState)

	// if initializers.CONFIG.ENV == initializers.ProductionENV {
	// 	parameters.Add("hd", config.VALID_DOMAINS[0])
	// }

	URL.RawQuery = parameters.Encode()
	url := URL.String()
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

func GoogleCallback(c *fiber.Ctx) error {
	state := c.FormValue("state")
	if state != config.GoogleOAuthState {
		return &fiber.Error{Code: 403, Message: "Invalid Callback State"}
	}

	code := c.FormValue("code")

	if code == "" {
		reason := c.FormValue("error_reason")
		if reason == "user_denied" {
			return &fiber.Error{Code: 403, Message: "User has denied Permission"}
		}
		return &fiber.Error{Code: 403, Message: "Code Not Found to provide AccessToken"}
		//redirect to /login of frontend
	} else {
		token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
		if err != nil {
			return &helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}

		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
		if err != nil {
			return &helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
		defer resp.Body.Close()

		response, err := io.ReadAll(resp.Body)
		if err != nil {
			return &helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}

		type GoogleUserInfo struct {
			Sub           string `json:"sub"`
			Name          string `json:"name"`
			GivenName     string `json:"given_name"`
			FamilyName    string `json:"family_name"`
			Profile       string `json:"profile"`
			Picture       string `json:"picture"`
			Email         string `json:"email"`
			EmailVerified bool   `json:"email_verified"`
		}

		var userInfo GoogleUserInfo
		err = json.Unmarshal(response, &userInfo)
		if err != nil {
			return &helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}

		// if initializers.CONFIG.ENV == initializers.ProductionENV {
		// 	if err := validators.EmailValidator(userInfo.Email); err != nil {
		// 		return err
		// 	}
		// }

		var user models.User
		if err := initializers.DB.Session(&gorm.Session{SkipHooks: true}).Preload("OAuth").First(&user, "email = ?", userInfo.Email).Error; err != nil {
			if err == gorm.ErrRecordNotFound {

				extractedName := userInfo.Name

				re := regexp.MustCompile(config.NAME_REGEX)
				if re.MatchString(userInfo.Name) {
					extractedName = re.ReplaceAllString(userInfo.Name, "")
				}

				newUser := models.User{
					Name:              extractedName,
					Email:             userInfo.Email,
					Username:          strings.Split(userInfo.Email, "@")[0], //* all emails are unique for a college
					PasswordChangedAt: time.Now(),
				}

				result := initializers.DB.Create(&newUser)
				if result.Error != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
				}

				oauth := models.OAuth{
					UserID:              newUser.ID,
					Provider:            "Google",
					OnBoardingCompleted: false,
				}

				result = initializers.DB.Create(&oauth)
				if result.Error != nil {
					return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
				}

				return RedirectToSignUp(c, newUser)
			} else {
				return &helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}
		}

		if user.OAuth.ID == uuid.Nil || user.OAuth.OnBoardingCompleted {
			return RedirectToLogin(c, user)
		}
		return RedirectToSignUp(c, user)
	}
}
