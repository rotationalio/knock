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
	payload := []byte("my super secret data")

	// Get existing secrets (if any)
	secrets, _ := knock.ListSecrets(parent)

	// Delete them and all their versions
	for _, secret := range secrets {
		fmt.Println(secret)
		knock.DeleteSecret(secret)
	}

	// Create a new secret
	secretPath, err := knock.CreateSecret(parent, secret_id)
	if err != nil {
		panic(fmt.Sprintln("Something went wrong; check credentials/if secret already exists!", err))
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
}
