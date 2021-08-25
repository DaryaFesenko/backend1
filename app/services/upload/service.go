package upload

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"path/filepath"
)

const (
	fileMode = 0777
)

type UploadService struct {
	UploadDir string
}

type FileData struct {
	Name      string
	SizeByte  int64
	Extension string
}

func NewUploadService(uploadDir string) *UploadService {
	return &UploadService{
		UploadDir: uploadDir,
	}
}

func (*UploadService) WriteFile(filePath string, file multipart.File) error {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Unable to read file")
	}

	err = ioutil.WriteFile(filePath, data, fileMode)
	if err != nil {
		return fmt.Errorf("Unable to save file")
	}

	return nil
}

func (us *UploadService) AllFiles() ([]FileData, error) {
	fileData := make([]FileData, 0)

	files, err := ioutil.ReadDir(us.UploadDir)
	if err != nil {
		return fileData, err
	}
	for _, f := range files {
		f := FileData{
			Name:      f.Name(),
			SizeByte:  f.Size(),
			Extension: filepath.Ext(f.Name()),
		}

		fileData = append(fileData, f)
	}

	return fileData, nil
}

func (us *UploadService) AllFilesFilterExt(ext string) ([]FileData, error) {
	fileData := make([]FileData, 0)

	files, err := ioutil.ReadDir(us.UploadDir)
	if err != nil {
		return fileData, err
	}
	for _, f := range files {
		extFile := filepath.Ext(f.Name())
		if extFile == ext {
			f := FileData{
				Name:      f.Name(),
				SizeByte:  f.Size(),
				Extension: extFile,
			}

			fileData = append(fileData, f)
		}
	}

	return fileData, nil
}
