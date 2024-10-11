package util

import "os"

func CheckDirIsExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 目录不存在, 创建目录
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
