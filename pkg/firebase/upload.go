package firebase

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go"
	"github.com/Akrom0181/Food-Delivery/internal/entity"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/storage/v1"
)

// UploadFiles uploads multiple files to Firebase Storage
func UploadFiles(file *multipart.Form) (*entity.MultipleFileUploadResponse, error) {
	var resp entity.MultipleFileUploadResponse

	filePath := filepath.Join("/app", "serviceAccountKey.json")

	// Initialize Firebase App with service account key
	opt := option.WithCredentialsFile(filePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("Firebase App initialization error:", err)
		return nil, err
	}

	client, err := app.Storage(context.TODO())
	if err != nil {
		log.Println("Firebase Storage client initialization error:", err)
		return nil, err
	}

	bucketHandle, err := client.Bucket("food-8ceb4.appspot.com")
	if err != nil {
		log.Println("Bucket handle error:", err)
		return nil, err
	}

	for _, v := range file.File["file"] {
		id := uuid.New().String()
		imageFile, err := v.Open()
		if err != nil {
			return nil, err
		}
		defer imageFile.Close()

		fileName := v.Filename

		objectHandle := bucketHandle.Object(fileName)
		writer := objectHandle.NewWriter(context.Background())
		writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}

		if _, err := io.Copy(writer, imageFile); err != nil {
			return nil, err
		}
		writer.Close()

		encodedFileName := url.PathEscape(fileName)
		fileURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/food-8ceb4.appspot.com/o/%s?alt=media&token=%s", encodedFileName, id)

		resp.Url = append(resp.Url, &entity.Url{
			Id:  id,
			Url: fileURL,
		})
	}

	return &resp, nil
}

// UploadFile uploads a single file to Firebase Storage
func UploadFile(file *os.File) (*entity.MultipleFileUploadResponse, error) {
	var resp entity.MultipleFileUploadResponse

	id := uuid.New().String()

	// Initialize Firebase App with service account key
	opt := option.WithCredentialsFile("/app/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	client, err := app.Storage(context.TODO())
	if err != nil {
		return nil, err
	}

	bucketHandle, err := client.Bucket("food-8ceb4.appspot.com")
	if err != nil {
		return nil, err
	}

	// Use the original file name
	fileName := file.Name()

	// Create a temporary file
	tempFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return nil, err
	}

	// Open the temporary file for reading
	f, err := os.Open(tempFile.Name())
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name())

	// Upload the file to Firebase Storage
	objectHandle := bucketHandle.Object(fileName)
	writer := objectHandle.NewWriter(context.Background())
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}

	defer writer.Close()

	if _, err := io.Copy(writer, f); err != nil {
		return nil, err
	}

	// Encode the file name to handle spaces and special characters
	encodedFileName := url.PathEscape(fileName)

	// Generate the public URL
	fileURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/food-8ceb4.appspot.com/o/%s?alt=media&token=%s", encodedFileName, id)

	// Add the URL to the response
	resp.Url = append(resp.Url, &entity.Url{
		Id:  id,
		Url: fileURL,
	})

	return &resp, nil
}

// DeleteFile deletes a file from Firebase Storage
func DeleteFile(id string) error {
	// Initialize a context and Google Cloud Storage client
	ctx := context.Background()
	client, err := storage.NewService(ctx, option.WithCredentialsFile("/app/serviceAccountKey.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Bucket name and object path to delete
	bucketName := "food-8ceb4.appspot.com"
	objectPath := id
	// Delete the object
	err = client.Objects.Delete(bucketName, objectPath).Do()
	if err != nil {
		log.Fatalf("Failed to delete object: %v", err)
	}

	fmt.Printf("Object %s deleted successfully from bucket %s\n", objectPath, bucketName)
	return nil
}
