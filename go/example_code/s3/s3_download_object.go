/*
   Copyright 2010-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.

   This file is licensed under the Apache License, Version 2.0 (the "License").
   You may not use this file except in compliance with the License. A copy of
   the License is located at

    http://aws.amazon.com/apache2.0/

   This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
   CONDITIONS OF ANY KIND, either express or implied. See the License for the
   specific language governing permissions and limitations under the License.
*/

package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    
    "fmt"
    "os"
)

// Downloads an item from an S3 Bucket in the region configured in the shared config
// or AWS_REGION environment variable.
//
// Usage:
//    go run s3_download_object.go BUCKET ITEM
func main() {
    if len(os.Args) != 3 {
        exitErrorf("Bucket and item names required\nUsage: %s bucket_name item_name",
            os.Args[0])
    }

    bucket := os.Args[1]
    item := os.Args[2]

    file, err := os.Create(item)
    if err != nil {
        exitErrorf("Unable to open file %q, %v", err)
    }

    defer file.Close()

    // Initialize a session in us-west-2 that the SDK will use to load
    // credentials from the shared credentials file ~/.aws/credentials.
    sess, _ := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
    )

    downloader := s3manager.NewDownloader(sess)

    numBytes, err := downloader.Download(file,
        &s3.GetObjectInput{
            Bucket: aws.String(bucket),
            Key:    aws.String(item),
        })
    if err != nil {
        exitErrorf("Unable to download item %q, %v", item, err)
    }

    fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}

func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}
