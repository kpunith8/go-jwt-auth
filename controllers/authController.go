package controllers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kpunith8/go-jwt-auth/database"
	"github.com/kpunith8/go-jwt-auth/models"
	"github.com/kpunith8/go-jwt-auth/utils"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func comparePasswords(hashedPassword string, plainPassword []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// Login - Login - POST request
func Login(c *fiber.Ctx) error {
	var requestData map[string]string

	client := database.DBConnection

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testDatabase := client.Database("test")
	goUserCollection := testDatabase.Collection("go-user")

	if err := c.BodyParser(&requestData); err != nil {
		log.Fatal(err)
	}

	user := models.User{}
	if err := goUserCollection.FindOne(ctx, bson.M{"email": requestData["email"]}).Decode(&user); err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found!",
		})
	}

	if isPasswordMatched := comparePasswords(user.Password, []byte(requestData["password"])); !isPasswordMatched {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect Password!",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    primitive.ObjectID(user.ID).Hex(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := claims.SignedString([]byte(utils.GetEnvVariable("JWT_SECRET")))

	if err != nil {
		log.Fatal(err)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

// User - Get the user based on cookie info - GET request
func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.GetEnvVariable("JWT_SECRET")), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized!",
		})
	}

	// Get the claims and get the issuer from the claim
	claims := token.Claims.(*jwt.StandardClaims)

	client := database.DBConnection

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testDatabase := client.Database("test")
	goUserCollection := testDatabase.Collection("go-user")

	user := models.User{}
	id, _ := primitive.ObjectIDFromHex(claims.Issuer)
	if err := goUserCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found!",
		})
	}

	return c.JSON(user.Email)
}

// Register - Register user - POST request
func Register(c *fiber.Ctx) error {
	// var requestData map[string]string
	user := models.User{}

	// Create a mongo objectID before inserting a document
	user.ID = primitive.NewObjectID()

	if err := c.BodyParser(&user); err != nil {
		log.Fatal(err)
	}

	client := database.DBConnection

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testDatabase := client.Database("test")
	goUserCollection := testDatabase.Collection("go-user")

	// Convert string to []byte
	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	// convert back to string before inserting to the DB
	user.Password = string(password)

	insertedUser, err := goUserCollection.InsertOne(ctx, bson.D{
		{"name", user.Name},
		{"age", user.Age},
		{"email", user.Email},
		{"password", user.Password},
	})

	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(insertedUser)
}

// AllUsers - Get all users - GET request
func AllUsers(c *fiber.Ctx) error {
	client := database.DBConnection

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testDatabase := client.Database("test")
	goUserCollection := testDatabase.Collection("go-user")

	usersCursor, err := goUserCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	modelGoUsers := []models.User{}

	if err = usersCursor.All(context.TODO(), &modelGoUsers); err != nil {
		log.Fatal(err)
	}

	return c.JSON(modelGoUsers)
}
