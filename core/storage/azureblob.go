package storage

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

// AzureStorageClient azure storage client
type AzureStorageClient struct {
	defaultContainer azblob.ContainerURL
}

// NewAzureClient new azure client
func NewAzureClient(accountName string, accountKey string, containerName string) *AzureStorageClient {
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		panic(errors.New("cannot connect azure blob service"))
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))
	containerURL := azblob.NewContainerURL(*URL, p)
	return &AzureStorageClient{defaultContainer: containerURL}
}

// Upload upload azure client
func (s *AzureStorageClient) Upload(file *os.File, fileName string, contentType string) error {
	ctx := context.Background()
	blobURL := s.defaultContainer.NewBlockBlobURL(fileName)
	_, err := azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: contentType,
		},
		Metadata:    map[string]string{},
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})
	if err != nil {
		return err
	}
	return nil
}
