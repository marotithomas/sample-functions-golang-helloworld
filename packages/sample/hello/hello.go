package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

var (
    minioClient *minio.Client
    bucketName  string
)

func main() {
    endpoint := "fra1.digitaloceanspaces.com"
    accessKeyID := "DO00CXZFVFZBJP9FP4C6"
    secretAccessKey := "ixuG36jRgBZjIouJ1w18Uig+nRKYXcdAshqeJqptfPI"
    useSSL := true
    bucketName = "teszt"

    // Initialize minio client object.
    var err error
    minioClient, err = minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
        Secure: useSSL,
    })
    if err != nil {
        log.Fatalln(err)
    }

    http.HandleFunc("/", handleRequest)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    objectKey := r.URL.Query().Get("key")

    if objectKey == "" {
        http.Error(w, "Missing 'key' parameter", http.StatusBadRequest)
        return
    }

    exists, err := checkObjectExists(objectKey)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if exists {
        fmt.Fprint(w, "OK")
    } else {
        fmt.Fprint(w, "FAILED")
    }
}

func checkObjectExists(objectKey string) (bool, error) {
    ctx := context.Background()
    _, err := minioClient.StatObject(ctx, bucketName, objectKey, minio.StatObjectOptions{})
    if err != nil {
        if minio.ToErrorResponse(err).Code == "NoSuchKey" {
            return false, nil
        }
        return false, err
    }
    return true, nil
}
