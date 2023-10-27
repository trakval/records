package fs

import (
	"errors"
	"fmt"
	"os"

	"github.com/gernest/front"
	"gopkg.in/yaml.v2"
)

type FsRecords struct {
	path string
}

func NewFsRecords(path string) FsRecords {
	return FsRecords{
		path,
	}
}

func (fsr FsRecords) getVaultDirPath() string {
	return fsr.path + "/vault"
}

func (fsr FsRecords) getFilePath(key string) string {
	return fsr.getVaultDirPath() + "/" + key + ".md"
}

func (fsr FsRecords) writeRecord(key string, record map[string]interface{}) (string, error) {
	fm, ok := record["frontmatter"].(map[string]interface{})
	if !ok {
		return key, errors.New("records/fs: error extracting 'frontmatter' in supplied record")
	}
	body, ok := record["body"].(string)
	if !ok {
		return key, errors.New("records/fs: error extracting 'body' in supplied record")
	}

	frontmatter, err := yaml.Marshal(fm)
	if err != nil {
		return key, nil
	}

	content := fmt.Sprintf(`---
%s
---
%s
`, frontmatter, body)

	return key, os.WriteFile(fsr.getFilePath(key), []byte(content), 0666)
}

func (fsr FsRecords) Connect() error {
	err := os.MkdirAll(fsr.getVaultDirPath(), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (fsr FsRecords) Close() error {
	return nil
}

func (fsr FsRecords) CreateRecord(key string, record map[string]interface{}) (string, error) {
	return fsr.writeRecord(key, record)
}

func (fsr FsRecords) UpdateRecord(key string, record map[string]interface{}) (string, error) {
	return fsr.writeRecord(key, record)
}

func (fsr FsRecords) DeleteRecord(key string) (string, error) {
	return key, os.Remove(fsr.getFilePath(key))
}

func (fsr FsRecords) FetchRecord(key string) (string, map[string]interface{}, error) {
	file, err := os.Open(fsr.getFilePath(key))
	if err != nil {
		return key, nil, err
	}

	fmp := front.NewMatter()
	fmp.Handle("---", front.YAMLHandler)
	frontmatter, body, err := fmp.Parse(file)
	if err != nil {
		return key, nil, err
	}

	record := map[string]interface{}{}
	record["frontmatter"] = frontmatter
	record["body"] = body

	return key, record, nil
}
