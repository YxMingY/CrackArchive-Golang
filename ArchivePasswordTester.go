package main

import (
	"github.com/mholt/archiver/v3"
	"github.com/yeka/zip"
	"strings"
)

func SelectTestFunc(fileName string) func(passwd string, filename string) int {
	suffix := fileName[strings.LastIndex(fileName, ".")+1:]
	switch suffix {
	case "zip":
		return TestZipPasswdCorrect
	case "rar":
		return TestRarPasswdCorrect
	default:
		return nil
	}
}

// TestRarPasswdCorrect
// @return 0:error,1:pass
func TestRarPasswdCorrect(passwd string, filename string) int {
	rar := archiver.Rar{}
	rar.Password = passwd
	return rar.IsPasswdCorrect(filename)
}

// TestZipPasswdCorrect
// @return 0:error,1:pass,2:file error
func TestZipPasswdCorrect(passwd string, filename string) int {
	// 1、使用zip.OpenReader打开zip文件
	archive, err := zip.OpenReader(filename)
	if err != nil {
		return 2
	}
	defer archive.Close()
	for _, f := range archive.File {
		if f.IsEncrypted() {
			f.SetPassword(passwd)
		}
		if f.FileInfo().IsDir() {
			continue
		}
		fileInArchive, err := f.Open()
		if err != nil {
			return 2
		}

		//_, err = io.ReadAll(fileInArchive)
		buf := make([]byte, 2000)
		_, err = fileInArchive.Read(buf)
		if err != nil {
			//fmt.Println(err)
			if strings.ContainsRune(err.Error(), 't') {
				return 0
			}
			//flate: corrupt input before offset 1  no
			//unexpected EOF  no
			// checksum error  yes
		}
		fileInArchive.Close()
		break
	}
	return 1
}
