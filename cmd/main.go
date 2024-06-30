package main

import (
	"context"
	"fmt"
	"log"
	"s3/internal/s3"
	"s3/pkg/api"
)

func main() {
	store := s3.New()

	req := &api.DownloadMediaRequest{
		Bucket: "pawpawchat",
		Key:    "testvideo.mp4",
	}

	info, err := store.Bucket().Object().Download(context.TODO(), req)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(info)
}
