package handler

import (
	"backend1/app/services/upload"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadHandler(t *testing.T) {
	dirPath := "../../upload_test"
	filePath := dirPath + "/test1.txt"
	data := []byte("test")
	clearOrCreateDir(dirPath)
	addFile(filePath, data)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("cant open file")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	us := upload.NewUploadService("../../upload_test")

	uploadHandler := &Router{
		HostAddr: "localhost:8080",
		us:       us,
	}
	uploadHandler.uploadFile(rr, req)

	path := rr.Body.String()
	code := rr.Code
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, path, "localhost:8080/test1.txt\n")

	file.Close()
	removeDir(dirPath)
}

func TestAllFiles(t *testing.T) {
	dirPath := "../../upload_test"
	data := []byte("test")
	testData := []upload.FileData{
		{
			Name:      "test1.txt",
			SizeByte:  int64(binary.Size(data)),
			Extension: ".txt",
		},
		{
			Name:      "test4.md",
			SizeByte:  int64(binary.Size(data)),
			Extension: ".md",
		},
		{
			Name:      "test3.json",
			SizeByte:  int64(binary.Size(data)),
			Extension: ".json",
		},
	}

	clearOrCreateDir(dirPath)
	for _, file := range testData {
		addFile(dirPath+"/"+file.Name, data)
	}

	body := &bytes.Buffer{}
	req, _ := http.NewRequest(http.MethodPost, "/all_files", body)

	rr := httptest.NewRecorder()
	us := upload.NewUploadService(dirPath)

	uploadHandler := &Router{
		HostAddr: "localhost:8000",
		us:       us,
	}
	uploadHandler.getAllFileInfo(rr, req)

	result := []upload.FileData{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		fmt.Println("can't unmarshal result data")
	}

	code := rr.Code
	assert.Equal(t, code, http.StatusOK)

	for _, val := range result {
		assert.Contains(t, testData, val)
	}

	removeDir(dirPath)
}

func TestFilesWithExt(t *testing.T) {
	dirPath := "../../upload_test"
	data := []byte("test")
	testData := []upload.FileData{
		{
			Name:      "test1.txt",
			SizeByte:  int64(binary.Size(data)),
			Extension: ".txt",
		},
		{
			Name:      "test3.txt",
			SizeByte:  int64(binary.Size(data)),
			Extension: ".txt",
		},
	}

	clearOrCreateDir(dirPath)
	for _, file := range testData {
		addFile(dirPath+"/"+file.Name, data)
	}
	addFile(dirPath+"/"+"badfile.md", data)

	body := &bytes.Buffer{}
	req, _ := http.NewRequest(http.MethodPost, "/all_files_ext?ext=.txt", body)

	rr := httptest.NewRecorder()
	us := upload.NewUploadService(dirPath)

	uploadHandler := &Router{
		HostAddr: "localhost:8000",
		us:       us,
	}
	uploadHandler.getFilesWithExt(rr, req)

	result := []upload.FileData{}
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		fmt.Println("can't unmarshal result data")
	}

	code := rr.Code
	assert.Equal(t, code, http.StatusOK)

	for _, val := range result {
		assert.Contains(t, testData, val)
	}

	removeDir(dirPath)
}

func addFile(filePath string, data []byte) {
	err := ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		fmt.Println("can't create file: ", filePath)
	}
}

func clearOrCreateDir(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		errDir := os.Mkdir(dirPath, 0777)
		if errDir != nil {
			fmt.Println("cant create dir")
		}
		return
	}

	err := os.RemoveAll(dirPath)
	if err != nil {
		fmt.Println("cant remove files")
	}
}

func removeDir(dirPath string) {
	err := os.RemoveAll(dirPath)
	if err != nil {
		fmt.Println("cant remove files")
	}

	os.Remove(dirPath)
	if err != nil {
		fmt.Println("cant remove dir")
	}
}
