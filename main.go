package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	endpoint := os.Getenv("ENDPOINT")
	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	if endpoint == "" || accessKey == "" || secretKey == "" {
		log.Fatal("Environment vars not configured")
	}
	useSSL := true

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal(err)
	}

	bucketName := "pal-save"
	location := "us-east-1"

	// ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	ctx := context.Background()
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("bucket already created")
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	log.Printf("bucket creation finished\n")

	objectName := "Saved.tar.gz"
	reader, err := minioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err == nil {
		defer reader.Close()
		localFile, err := os.Create("Saved.tar.gz")
		if err != nil {
			log.Fatalln(err)
		}
		defer localFile.Close()

		stat, err := reader.Stat()
		if err != nil {
			log.Fatalln(err)
		}

		if _, err := io.CopyN(localFile, reader, stat.Size); err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Removing existing Saved files\n")
		cmd := exec.Command("rm -rf /palworld/Pal/Saved/*")
		cmd.Run()
		fmt.Printf("Extracting Backup\n")
		cmd = exec.Command("tar -zxf Saved.tar.gz -C /palworld/Pal/Saved")
		cmd.Run()
		fmt.Printf("Extraction finished\n")
	}
}
