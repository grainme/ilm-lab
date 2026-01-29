package main

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/secrets"
)

func main() {
	// this load the config from ~/.oci/config
	configProvider := common.DefaultConfigProvider()

	secretClient, err := secrets.NewSecretsClientWithConfigurationProvider(configProvider)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	secretOCID := "ocid1.vaultsecret.oc1.region_name.some_id"

	res, err := secretClient.GetSecretBundle(context.Background(),
		secrets.GetSecretBundleRequest{
			SecretId: &secretOCID,
		})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	bundle, ok := res.SecretBundleContent.(secrets.Base64SecretBundleContentDetails)
	if !ok {
		log.Fatal("unexpected secret bundle content type")
	}

	decoded, err := base64.StdEncoding.DecodeString(*bundle.Content)
	if err != nil {
		log.Fatalf("decode error: %v", err)
	}

	log.Println("Secret:", string(decoded))

}
