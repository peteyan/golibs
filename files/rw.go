package files

import (
	"os"
)

// ReadFileToUint8Array 读取文件转化为无符号字节数组
func ReadFileToUint8Array(filePath string) ([]uint8, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	unsignedByteArray := make([]uint8, len(fileContent))
	for i, b := range fileContent {
		unsignedByteArray[i] = b
	}
	return unsignedByteArray, nil
}

// ReadFileToInt8Array 读取文件转化为有符号字节数组，与Java中的处理一致
func ReadFileToInt8Array(filePath string) ([]int8, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	signedByteArray := make([]int8, len(fileContent))
	for i, b := range fileContent {
		signedByteArray[i] = int8(b)
	}
	return signedByteArray, nil
}
