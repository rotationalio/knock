package knock

import (
	"context"
	"fmt"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// CreateSecret creates a new secret with the given name. A secret is a logical
// wrapper around a collection of secret versions. Secret versions hold the
// actual secret material.
func CreateSecret(parent, id string) error {
	// parent := "projects/my-project"
	// id := "my-secret"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: id,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	// Call the API.
	result, err := client.CreateSecret(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create secret: %v", err)
	}
	fmt.Printf("Created secret: %s\n", result.Name)
	return nil
}

// AddSecretVersion adds a new secret version to the given secret with the
// provided payload.
func AddSecretVersion(w io.Writer, parent string) error {
	// parent := "projects/my-project/secrets/my-secret"

	// Declare the payload to store.
	payload := []byte("my super secret data")

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AddSecretVersionRequest{
		Parent: parent,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	// Call the API.
	result, err := client.AddSecretVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to add secret version: %v", err)
	}
	fmt.Fprintf(w, "Added secret version: %s\n", result.Name)
	return nil
}

// AccessSecretVersion accesses the payload for the given secret version if one
// exists. The version can be a version number as a string (e.g. "5") or an
// alias (e.g. "latest").
func AccessSecretVersion(w io.Writer, name string) error {
	// name := "projects/my-project/secrets/my-secret/versions/5"
	// name := "projects/my-project/secrets/my-secret/versions/latest"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to access secret version: %v", err)
	}

	// WARNING: Do not print the secret in a production environment - this snippet
	// is showing how to access the secret material.
	fmt.Fprintf(w, "Plaintext: %s\n", string(result.Payload.Data))
	return nil
}
