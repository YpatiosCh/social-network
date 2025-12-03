package handlers

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"log"
// 	"mime/multipart"
// 	"os"
// 	"time"

// 	"github.com/minio/minio-go/v7"
// 	"github.com/minio/minio-go/v7/pkg/credentials"
// )

// func NewMinIOConn() *minio.Client {
// 	var minioClient *minio.Client
// 	var err error
// 	for range 10 {
// 		minioClient, err = minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
// 			Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), ""),
// 			Secure: false,
// 		})
// 		if err == nil {
// 			break
// 		}
// 		log.Println("MinIO not ready, retrying in 2s...")
// 		time.Sleep(2 * time.Second)
// 	}
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	// Ensure bucket exists
// 	bucket := "images"
// 	ctx := context.Background()
// 	exists, errBucketExists := minioClient.BucketExists(ctx, bucket)
// 	if errBucketExists != nil {
// 		log.Fatalln(errBucketExists)
// 	}
// 	if !exists {
// 		err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 		fmt.Println("Created bucket:", bucket)
// 	}
// 	return minioClient
// }

// func (h *Handlers) uploadToMinIO(
// 	ctx context.Context,
// 	file multipart.File,
// 	header *multipart.FileHeader,
// 	bucket string,
// 	contentType string,
// ) (minio.UploadInfo, error) {

// 	tmp, err := os.CreateTemp("", "upload-*")
// 	if err != nil {
// 		return minio.UploadInfo{}, err
// 	}
// 	defer os.Remove(tmp.Name())
// 	defer tmp.Close()

// 	if _, err = io.Copy(tmp, file); err != nil {
// 		return minio.UploadInfo{}, err
// 	}

// 	info, err := h.MinIOClient.FPutObject(
// 		ctx,
// 		bucket,
// 		header.Filename,
// 		tmp.Name(),
// 		minio.PutObjectOptions{
// 			ContentType: contentType,
// 		},
// 	)
// 	if err != nil {
// 		return minio.UploadInfo{}, err
// 	}

// 	return info, nil
// }
