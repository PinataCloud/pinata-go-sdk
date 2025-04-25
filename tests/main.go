package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/PinataCloud/pinata-go-sdk/pinata"
	"github.com/PinataCloud/pinata-go-sdk/pinata/files"
	"github.com/PinataCloud/pinata-go-sdk/pinata/upload"
)

func main() {
	// Initialize with your JWT and Gateway
	client := pinata.New(
		os.Getenv("PINATA_JWT"),
		os.Getenv("PINATA_GATEWAY"),
	)

	file, err := os.Open("./pinnie.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create upload options
	opts := &upload.FileOptions{
		FileName: "Pinnie",
	}

	// Upload to public IPFS
	resp, err := client.Upload.Private.File(file, opts)
	if err != nil {
		log.Fatalf("Upload failed: %v", err)
	}

	// Log the entire response by decoding the JSON
	respJSON, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal response: %v", err)
	}
	fmt.Printf("Complete response:\n%s\n", string(respJSON))

	fmt.Printf("File uploaded with CID: %s\n", resp.CID)

	fileData, err := client.Files.Private.List(&files.ListOptions{
		Limit: 1,
		CID: resp.CID,
	})
	if err != nil {
		log.Fatalf("List files faile: %v", err)
	}

	fileDataJson, err := json.MarshalIndent(fileData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal repsonse: %v", err)
	}

	fmt.Printf("List files:\n%s\n", fileDataJson)
}
