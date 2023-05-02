package main

import (
	"context"
	"fmt"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	// Set the project ID, key ring name, and key name
	projectID := "projects/kubebuilder-try-1/locations/global"
	keyRingName := "my-key-ring"
	keyName := "my-key"

	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		fmt.Printf("failed to create kms client: %v\n", err)
		return
	}
	defer client.Close()

	keyRingPath := fmt.Sprintf("%s/keyRings/%s", projectID, keyRingName)

	// Check if the KeyRing exists
	_, err = client.GetKeyRing(ctx, &kmspb.GetKeyRingRequest{Name: keyRingPath})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// KeyRing doesn't exist, create it
			req := &kmspb.CreateKeyRingRequest{
				Parent:    projectID,
				KeyRingId: keyRingName,
			}

			result, err := client.CreateKeyRing(ctx, req)
			if err != nil {
				fmt.Printf("failed to create key ring: %v\n", err)
				return
			}
			fmt.Printf("Created key ring: %s\n", result.Name)
		} else {
			fmt.Printf("failed to check key ring: %v\n", err)
			return
		}
	} else {
		fmt.Printf("Key ring already exists: %s\n", keyRingPath)
	}

	// Check if the Key exists
	keyPath := fmt.Sprintf("%s/cryptoKeys/%s", keyRingPath, keyName)
	_, err = client.GetCryptoKey(ctx, &kmspb.GetCryptoKeyRequest{Name: keyPath})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Key doesn't exist, create it
			req := &kmspb.CreateCryptoKeyRequest{
				Parent:      keyRingPath,
				CryptoKeyId: keyName,
				CryptoKey: &kmspb.CryptoKey{
					Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
				},
			}

			result, err := client.CreateCryptoKey(ctx, req)
			if err != nil {
				fmt.Printf("failed to create key: %v\n", err)
			} else {
				fmt.Printf("Created key: %s\n", result.Name)
			}
		} else {
			fmt.Printf("failed to check key: %v\n", err)
		}
	} else {
		fmt.Printf("Key already exists: %s\n", keyPath)
	}
}
