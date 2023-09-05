package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
)

// 文件信息
type FileInfo struct {
	Name string
	Path string
	Hash string
}

// 默认文件路径
const saveFilesPath = "web/files"

// 文件信息集合
var fileInfos []FileInfo

// 创建Hash表
func CreateHash() {

	// 检查是否存在sha256.json文件
	if _, err := os.ReadFile(saveFilesPath + "/sha256.json"); err != nil {
		// 错误的话是因为没有创建files目录
		os.Mkdir("web/files", 0700)
		data, _ := json.MarshalIndent(fileInfos, "", "	")
		os.WriteFile(saveFilesPath+"/sha256.json", data, 0600)

	}
}

// // 将目录下全部文件写入
// func allWriteHash() {
// 	// 递归查找本地文件，返回值为全局变量fileInfos
// 	RecursionRerurnFiles(saveFilesPath)

// 	// 写入json文件
// 	data, _ := json.MarshalIndent(fileInfos, "", "	")
// 	os.WriteFile(saveFilesPath+"/sha256.json", data, 0600)
// 	fileInfos = []FileInfo{}
// }

// 新增Hash
func addHash(fileName string) {
	// 临时存储变量
	tmpInfo := FileInfo{}

	// 开始赋值
	tmpInfo.Name = fileName
	tmpInfo.Path = "web/files/" + fileName
	data, err := os.ReadFile(tmpInfo.Path)
	ErrprDisplay(err)

	//计算Hash
	tmpInfo.Hash = CountHash(data)

	// 将临时存储信息添加到整体
	fileinfos := tmpInfo.getHash()
	fileinfos = append(fileinfos, tmpInfo)

	// 转换为json格式数据
	filedata, err := json.MarshalIndent(fileinfos, "", "	")
	ErrprDisplay(err)

	// 写入json格式
	err = os.WriteFile(saveFilesPath+"/sha256.json", filedata, 0600)
	ErrprDisplay(err)
}

// 获取Hash
func (f FileInfo) getHash() []FileInfo {
	data, err := os.ReadFile(saveFilesPath + "/sha256.json")
	ErrprDisplay(err)
	fileinfos := []FileInfo{}
	err = json.Unmarshal(data, &fileInfos)
	ErrprDisplay(err)
	return fileinfos
}

// 递归返回多个文件信息
// 使用完该函数需要初始化 fileInfos=nil
func RecursionRerurnFiles(dirName string) {
	// 临时存储文件信息
	tmpSave := FileInfo{}

	// 读取文件目录
	dir, err := os.ReadDir(dirName)
	ErrprDisplay(err)

	// 递归开始
	for _, v := range dir {
		// 判断是不是目录
		if v.IsDir() {
			// 是目录递归执行
			RecursionRerurnFiles(dirName + "/" + v.Name())
		} else {
			// 获取文件数据
			data, _ := os.ReadFile(dirName + "/" + v.Name())
			ErrprDisplay(err)

			// 保存文件三要素
			tmpSave.Hash = CountHash(data)
			tmpSave.Name = v.Name()
			tmpSave.Path = dirName + "/" + v.Name()

			// 将临时存储的信息存入文件信息集合
			fileInfos = append(fileInfos, tmpSave)
		}
	}
}

// 计算Hash
func CountHash(data []byte) (hashString string) {
	hashByte := sha256.Sum256(data)
	hashString = hex.EncodeToString(hashByte[:])
	return hashString
}

// 检查Hash是否真确
func CheckingHash() (success []FileInfo, fail []FileInfo) {
	// 保存解析sha256.json的数据
	shaSaveData := []FileInfo{}

	// 读取Hash文件
	data, err := os.ReadFile(saveFilesPath + "/sha256.json")
	ErrprDisplay(err)

	// 开始解析
	err = json.Unmarshal(data, &shaSaveData)
	ErrprDisplay(err)

	// 开始验证
	success, fail = shaVerify(shaSaveData)

	return success, fail
}

// 从sha256.json验证是否被篡改
func shaVerify(shaSaveData []FileInfo) (success []FileInfo, fail []FileInfo) {
	for _, v := range shaSaveData {
		data, err := os.ReadFile(v.Path)
		if err != nil {
			v.Path = err.Error()
			fail = append(fail, v)
		}

		if v.Hash == CountHash(data) {
			success = append(success, v)
		} else {
			if v.Name != "sha256.json" {
				fail = append(fail, v)
			}
		}
	}

	return success, fail
}
