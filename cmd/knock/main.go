package main

import (
	"fmt"
	"os"

	"github.com/rotationalio/knock"
)

func main() {
	projectID, ok := os.LookupEnv("GOOGLE_PROJECT_NAME")
	if !ok {
		panic(fmt.Sprintln("project name not found; check that environment variable is set"))
	}

	parent := fmt.Sprintf("projects/%s", projectID)
	secret_id := "test"
	duration := int64(60) // in seconds
	payload := []byte("my super secret data")

	// Knock to check if our permission level will let us connect
	err := knock.Knock(parent)
	if err != nil {
		panic(fmt.Sprintln("Something went wrong; you don't seem to have access to the project: ", err))
	}

	// Get existing secrets (if any)
	secrets, _ := knock.ListSecrets(parent)

	// Delete them and all their versions
	for _, secret := range secrets {
		fmt.Println(secret)
		knock.DeleteSecret(secret)
	}

	// Create a new secret that expires after duration
	secretPath, err := knock.CreateSecret(parent, secret_id, duration)
	if err != nil {
		panic(fmt.Sprintln("Something went wrong; check if secret already exists!", err))
	}

	// Add a new secret version with the payload
	versionPath, err := knock.AddSecretVersion(secretPath, payload)
	if err != nil {
		panic(fmt.Sprintln("Couldn't add secret version; check service account permissions!", err))
	}

	// Retrieve the payload for that version
	retrieved, err := knock.AccessSecretVersion(versionPath)
	if err != nil {
		panic(fmt.Sprintln("Couldn't retrieve secret; check service account permissions!", err))
	}
	// Note - just for demo; don't actually print out secrets in practice!
	fmt.Printf("found your secret: %s\n", retrieved)

	// Retrieve payload for the latest version
	latest, err := knock.AccessSecretVersion(fmt.Sprintf("%s/%s/%s", secretPath, "versions", "latest"))
	if err != nil {
		panic(fmt.Sprintln("Couldn't retrieve secret; check service account permissions!", err))
	}
	// Note - just for demo; don't actually print out secrets in practice!
	fmt.Printf("found your secret: %s\n", latest)
}
