package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			return true, nil
		}
	}
	return false, err
}

func IsEmpty(strs ...string) (isEmpty bool) {
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if str == "" || len(str) == 0 {
			isEmpty = true
			return
		}
	}
	isEmpty = false
	return
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
