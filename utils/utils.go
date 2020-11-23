package utils

import "os"

// 将content写入fileName 文件中
// 将文本内容写入文件中
func WriteToFile(fileName,content string){
	f, err := os.Create(fileName)
	if err !=nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err !=nil {
		panic(err)
	}
}

// 一个文件或者目录是否存在
func IsExist(path string)bool{
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}