package utils

import "fmt"

func GetFileSize(size float64) (size_str string) {
	var str string;
	limit := float64(1 << 20)

	if size > limit{
		str = fmt.Sprintf("%.2f MB", size / limit)
	}else{
		str = fmt.Sprintf("%.2f KB", size / 1024)
	}
	return str;
}
