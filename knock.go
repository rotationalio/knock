package knock

import (
	"context"
	"errors"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"google.golang.org/api/iterator"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// Knock checks to make sure we can create a new client.
// This validates IAM permissions to some extent.
func Knock(parent string) error {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		switch err.Error() {
		case "google: could not find default credentials. See https://developers.google.com/accounts/docs/application-default-credentials for more information.":
			return errors.New("service account doesn't have permissions to create client")
		default:
			return err
		}
	}
	client.Close()
	return nil
}

// CreateSecret creates a new secret in the Google Cloud Manager top-
// level directory, specified as `parent`, using the `secretID` provided
// as the name. The parent should be a path, e.g.
//     "projects/project-name"
// This function returns a string representation of the path where the
// new secret is stored, e.g.
//     "projects/projectID/secrets/secretID"
// and an error if any occurs.
// Note: A secret is a logical wrapper around a collection of secret versions.
// Secret versions hold the actual secret material.
func CreateSecret(parent, secretID string) (string, error) {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		// The most likely causes of the error are:
		//     1 - that google application creds failed
		//     2 - secret already exists
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: secretID,
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
		return "", fmt.Errorf("failed to create secret: %v", err)
	}
	fmt.Printf("created secret: %s\n", result.Name)
	return result.Name, nil
}

// AddSecretVersion adds a new secret version to the given secret path with the
// provided payload. The path should be the full path to the secret, e.g.
//     "projects/projectID/secrets/secretID"
// Returns the path to the secret version, e.g.:
//     "projects/projectID/secrets/secretID/versions/1"
// and an error if one occurs.
func AddSecretVersion(path string, payload []byte) (string, error) {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AddSecretVersionRequest{
		Parent: path,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	// Call the API.
	result, err := client.AddSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to add secret version: %v", err)
	}

	fmt.Printf("added secret version: %s\n", result.Name)
	return result.Name, nil
}

// AccessSecretVersion returns the payload for the given secret version if one
// exists. The `version` is the full path to the secret version, and can be a
// version number as a string (e.g. "5") or an alias (e.g. "latest"), i.e.
//     "projects/projectID/secrets/secretID/versions/latest"
//     "projects/projectID/secrets/secretID/versions/5"
func AccessSecretVersion(version string) ([]byte, error) {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: version,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %v", err)
	}

	fmt.Printf("retrieved payload for: %s\n", result.Name)
	return result.Payload.Data, nil
}

// DeleteSecret deletes the secret with the given `name`, and all of its versions.
// `name` should be the root path to the secret, e.g.:
//     "projects/projectID/secrets/secretID"
// This is an irreversible operation. Any service or workload that attempts to
// access a deleted secret receives a Not Found error.
func DeleteSecret(name string) error {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.DeleteSecretRequest{
		Name: name,
	}

	// Call the API.
	if err := client.DeleteSecret(ctx, req); err != nil {
		return fmt.Errorf("failed to delete secret: %v", err)
	}
	return nil
}

// ListSecrets retrieves the names of all secrets in the project,
// given the `parent`, e.g.:
//     "projects/my-project"
// It returns a slice of strings representing the paths to the retrieved secrets,
// and a matching slice of errors for each failed retrieval.
func ListSecrets(parent string) (secrets []string, errors []error) {

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return secrets, append(errors, err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.ListSecretsRequest{
		Parent: parent,
	}

	// Call the API.
	it := client.ListSecrets(ctx, req)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			errors = append(errors, err)
			secrets = append(secrets, "")
			continue
		}
		secrets = append(secrets, resp.Name)
		errors = append(errors, nil)
	}
	return secrets, errors
}
