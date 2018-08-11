package gzip

import (
	"bytes"
	"compress/gzip"
	"log"
	"time"
	"testing"
)

func TestGzipCompress(t *testing.T) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	// 设置标题字段是可选的。
	zw.Name = "a-new-hope.txt"
	zw.Comment = "an epic space opera by George Lucas"
	zw.ModTime = time.Date(1977, time.May, 25, 0, 0, 0, 0, time.UTC)

	_, err := zw.Write([]byte("A long time ago in a galaxy far, far away..."))
	if err != nil {
		log.Fatal(err)
	}

	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}

	//zr, err := gzip.NewReader(&buf)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Printf("Name: %s\nComment: %s\nModTime: %s\n\n", zr.Name, zr.Comment, zr.ModTime.UTC())
	//
	//if _, err := io.Copy(os.Stdout, zr); err != nil {
	//	log.Fatal(err)
	//}
	//
	//if err := zr.Close(); err != nil {
	//	log.Fatal(err)
	//}

}
