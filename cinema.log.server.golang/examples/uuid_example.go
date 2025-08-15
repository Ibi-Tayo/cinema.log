package main

import (
	"fmt"
	"log"

	"cinema.log.server.golang/internal/domain"
	"cinema.log.server.golang/internal/utils"
)

func main() {
	// Example of creating a user with UUID
	user := &domain.User{
		GithubID:      123456,
		Name:          "John Doe",
		Username:      "johndoe",
		ProfilePicURL: "https://example.com/avatar.jpg",
	}

	// Generate a UUID for the user
	user.ID = utils.GenerateUUID()
	
	fmt.Printf("Generated User:\n")
	fmt.Printf("  ID: %s\n", user.ID.String())
	fmt.Printf("  Github ID: %s\n", user.GithubID)
	fmt.Printf("  Name: %s\n", user.Name)
	fmt.Printf("  Username: %s\n", user.Username)
	
	// Example of parsing UUID from string
	uuidStr := user.ID.String()
	parsedUUID, err := utils.ParseUUID(uuidStr)
	if err != nil {
		log.Fatal("Failed to parse UUID:", err)
	}
	
	fmt.Printf("\nParsed UUID: %s\n", parsedUUID.String())
	fmt.Printf("UUIDs match: %t\n", user.ID == parsedUUID)
	
	// Example of validating UUID
	fmt.Printf("\nUUID Validation:\n")
	fmt.Printf("Valid UUID '%s': %t\n", uuidStr, utils.ValidateUUID(uuidStr))
	fmt.Printf("Valid UUID 'invalid-uuid': %t\n", utils.ValidateUUID("invalid-uuid"))
	
	// Check for nil UUID
	var nilUUID domain.User
	fmt.Printf("Is nil UUID: %t\n", utils.IsNilUUID(nilUUID.ID))
}
