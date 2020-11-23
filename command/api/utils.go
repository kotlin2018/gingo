package api

var md5 = `package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}
`

var utilsValidator = `package utils

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type Rules map[string][]string

type RulesMap map[string]Rules

var CustomizeMap = make(map[string]Rules)

// 注册自定义规则方案建议在路由初始化层即注册
func RegisterRule(key string, rule Rules) (err error) {
	if CustomizeMap[key] != nil {
		return errors.New(key + "已注册,无法重复注册")
	} else {
		CustomizeMap[key] = rule
		return nil
	}
}

// 非空 不能为其对应类型的0值
func NotEmpty() string {
	return "notEmpty"
}

// 小于入参(<) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
func Lt(mark string) string {
	return "lt=" + mark
}

// 小于等于入参(<=) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
func Le(mark string) string {
	return "le=" + mark
}

// 等于入参(==) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
func Eq(mark string) string {
	return "eq=" + mark
}

// 不等于入参(!=)  如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
func Ne(mark string) string {
	return "ne=" + mark
}

// 大于等于入参(>=) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
func Ge(mark string) string {
	return "ge=" + mark
}

// 大于入参(>) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
func Gt(mark string) string {
	return "gt=" + mark
}

// 校验方法 接收两个参数  入参实例，规则map
func Verify(st interface{}, roleMap Rules) (err error) {
	compareMap := map[string]bool{
		"lt": true,
		"le": true,
		"eq": true,
		"ne": true,
		"ge": true,
		"gt": true,
	}

	typ := reflect.TypeOf(st)
	val := reflect.ValueOf(st) // 获取reflect.Type类型

	kd := val.Kind() // 获取到st对应的类别
	if kd != reflect.Struct {
		return errors.New("expect struct")
	}
	num := val.NumField()
	// 遍历结构体的所有字段
	for i := 0; i < num; i++ {
		tagVal := typ.Field(i)
		val := val.Field(i)
		if len(roleMap[tagVal.Name]) > 0 {
			for _, v := range roleMap[tagVal.Name] {
				switch {
				case v == "notEmpty":
					if isBlank(val) {
						return errors.New(tagVal.Name + "值不能为空")
					}
				case compareMap[strings.Split(v, "=")[0]]:
					if !compareVerify(val, v) {
						return errors.New(tagVal.Name + "长度或值不在合法范围," + v)
					}
				}
			}
		}
	}
	return nil
}

// 长度和数字的校验方法 根据类型自动校验
func compareVerify(value reflect.Value, VerifyStr string) bool {
	switch value.Kind() {
	case reflect.String, reflect.Slice, reflect.Array:
		return compare(value.Len(), VerifyStr)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compare(value.Uint(), VerifyStr)
	case reflect.Float32, reflect.Float64:
		return compare(value.Float(), VerifyStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compare(value.Int(), VerifyStr)
	default:
		return false
	}
}

// 非空校验
func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

func compare(value interface{}, VerifyStr string) bool {
	VerifyStrArr := strings.Split(VerifyStr, "=")
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		VInt, VErr := strconv.ParseInt(VerifyStrArr[1], 10, 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Int() < VInt
		case VerifyStrArr[0] == "le":
			return val.Int() <= VInt
		case VerifyStrArr[0] == "eq":
			return val.Int() == VInt
		case VerifyStrArr[0] == "ne":
			return val.Int() != VInt
		case VerifyStrArr[0] == "ge":
			return val.Int() >= VInt
		case VerifyStrArr[0] == "gt":
			return val.Int() > VInt
		default:
			return false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		VInt, VErr := strconv.Atoi(VerifyStrArr[1])
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Uint() < uint64(VInt)
		case VerifyStrArr[0] == "le":
			return val.Uint() <= uint64(VInt)
		case VerifyStrArr[0] == "eq":
			return val.Uint() == uint64(VInt)
		case VerifyStrArr[0] == "ne":
			return val.Uint() != uint64(VInt)
		case VerifyStrArr[0] == "ge":
			return val.Uint() >= uint64(VInt)
		case VerifyStrArr[0] == "gt":
			return val.Uint() > uint64(VInt)
		default:
			return false
		}
	case reflect.Float32, reflect.Float64:
		VFloat, VErr := strconv.ParseFloat(VerifyStrArr[1], 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Float() < VFloat
		case VerifyStrArr[0] == "le":
			return val.Float() <= VFloat
		case VerifyStrArr[0] == "eq":
			return val.Float() == VFloat
		case VerifyStrArr[0] == "ne":
			return val.Float() != VFloat
		case VerifyStrArr[0] == "ge":
			return val.Float() >= VFloat
		case VerifyStrArr[0] == "gt":
			return val.Float() > VFloat
		default:
			return false
		}
	default:
		return false
	}
}
`

var zipFiles = `package utils

import (
	"archive/zip"
	"io"
	"os"
	"strings"
)

func ZipFiles(filename string, files []string, oldform, newform string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 把files添加到zip中
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// 获取file的基础信息
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 使用上面的FileInforHeader() 就可以把文件保存的路径替换成我们自己想要的了，如下面
		header.Name = strings.Replace(file, oldform, newform, -1)

		// 优化压缩
		// 更多参考see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}
`

var structToMap = `package utils

import "reflect"

// 利用反射将结构体转化为map
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}
`

var breakpoint = `package utils

import (
	"io/ioutil"
	"os"
	"strconv"
)

// 前端传来文件片与当前片为什么文件的第几片
// 后端拿到以后比较次分片是否上传 或者是否为不完全片
// 前端发送每片多大
// 前端告知是否为最后一片且是否完成

const breakpointDir = "./breakpointDir/"
const finishDir = "./fileDir/"

func BreakPointContinue(content []byte, fileName string, contentNumber int, contentTotal int, fileMd5 string) (error, string) {
	path := breakpointDir + fileMd5 + "/"
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err, path
	}
	err, pathc := makeFileContent(content, fileName, path, contentNumber)
	return err, pathc

}

func CheckMd5(content []byte, chunkMd5 string) (CanUpload bool) {
	fileMd5 := MD5V(content)
	if fileMd5 == chunkMd5 {
		return true // "可以继续上传"
	} else {
		return false // "切片不完整，废弃"
	}
}

func makeFileContent(content []byte, fileName string, FileDir string, contentNumber int) (error, string) {
	path := FileDir + fileName + "_" + strconv.Itoa(contentNumber)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err, path
	} else {
		_, err = f.Write(content)
		if err != nil {
			return err, path
		}
	}
	return nil, path
}

func MakeFile(fileName string, FileMd5 string) (error, string) {
	rd, err := ioutil.ReadDir(breakpointDir + FileMd5)
	if err != nil {
		return err, finishDir + fileName
	}
	_ = os.MkdirAll(finishDir, os.ModePerm)
	fd, _ := os.OpenFile(finishDir+fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	for k := range rd {
		content, _ := ioutil.ReadFile(breakpointDir + FileMd5 + "/" + fileName + "_" + strconv.Itoa(k))
		_, err = fd.Write(content)
		if err != nil {
			_ = os.Remove(finishDir + fileName)
			return err, finishDir + fileName
		}
	}
	defer fd.Close()
	return nil, finishDir + fileName
}

func RemoveChunk(FileMd5 string) error {
	err := os.RemoveAll(breakpointDir + FileMd5)
	return err
}
`

var uploadAvatar = `package utils

import (
	"{{.Appname}}/global"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

func UploadAvatarLocal(file *multipart.FileHeader) (err error, localPath string, key string) {
	// 读取文件后缀
	ext := path.Ext(file.Filename)
	// 读取文件名并加密
	fileName := strings.TrimSuffix(file.Filename, ext)
	fileName = MD5V([]byte(fileName))
	// 拼接新文件名
	lastName := fileName + "_" + time.Now().Format("20060102150405") + ext
	// 读取全局变量的定义路径
	savePath := global.GVA_CONFIG.LocalUpload.AvatarPath
	// 尝试创建此路径
	err = os.MkdirAll(savePath, os.ModePerm)
	if err != nil{
		global.GVA_LOG.Error("upload local file fail:", err)
		return err, "", ""
	}
	// 拼接路径和文件名
	dst := savePath + "/" + lastName
	// 下面为上传逻辑
	// 打开文件 defer 关闭
	src, err := file.Open()
	if err != nil {
		global.GVA_LOG.Error("upload local file fail:", err)
		return err, "", ""
	}
	defer src.Close()
	// 创建文件 defer 关闭
	out, err := os.Create(dst)
	if err != nil {
		global.GVA_LOG.Error("upload local file fail:", err)
		return err, "", ""
	}
	defer out.Close()
	// 传输（拷贝）文件
	_, err = io.Copy(out, src)
	if err != nil {
		global.GVA_LOG.Error("upload local file fail:", err)
		return err, "", ""
	}
	return nil, dst, lastName
}
`
var uploadFile = `package utils

import (
	"{{.Appname}}/global"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

func UploadFileLocal(file *multipart.FileHeader) (err error, localPath string, key string) {
	// 读取文件后缀
	ext := path.Ext(file.Filename)
	// 读取文件名并加密
	fileName := strings.TrimSuffix(file.Filename, ext)
	fileName = MD5V([]byte(fileName))
	// 拼接新文件名
	lastName := fileName + "_" + time.Now().Format("20060102150405") + ext
	// 读取全局变量的定义路径
	savePath := global.GVA_CONFIG.LocalUpload.FilePath
	// 尝试创建此路径
	err = os.MkdirAll(savePath, os.ModePerm)
	if err != nil{
		global.GVA_LOG.Error("upload local file fail:", err)
		return err, "", ""
	}
	// 拼接路径和文件名
	dst := savePath + "/" + lastName
	// 下面为上传逻辑
	// 打开文件 defer 关闭
	src, err := file.Open()
	if err != nil {
		global.GVA_LOG.Error("upload local file fail:", err)
		return err, "", ""
	}
	defer src.Close()
	// 创建文件 defer 关闭
	out, err := os.Create(dst)
	if err != nil {
		global.GVA_LOG.Error("upload local file fail:", err)
		return err, "", ""
	}
	defer out.Close()
	// 传输（拷贝）文件
	_, err = io.Copy(out, src)
	if err != nil {
		global.GVA_LOG.Error("upload local file fail:", err)
		return err, "", ""
	}
	return nil, dst, lastName
}
`

var uploadRemote = `package utils

import (
	"context"
	"fmt"
	"{{.Appname}}/global"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"mime/multipart"
	"time"
)

// 接收两个参数 一个文件流 一个 bucket 你的七牛云标准空间的名字
func UploadRemote(file *multipart.FileHeader) (err error, path string, key string) {
	putPolicy := storage.PutPolicy{
		Scope: global.GVA_CONFIG.Qiniu.Bucket,
	}
	mac := qbox.NewMac(global.GVA_CONFIG.Qiniu.AccessKey, global.GVA_CONFIG.Qiniu.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	f, e := file.Open()
	if e != nil {
		fmt.Println(e)
		return e, "", ""
	}
	dataLen := file.Size
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename) // 文件名格式 自己可以改 建议保证唯一性
	err = formUploader.Put(context.Background(), &ret, upToken, fileKey, f, dataLen, &putExtra)
	if err != nil {
		global.GVA_LOG.Error("upload file fail:", err)
		return err, "", ""
	}
	return err, global.GVA_CONFIG.Qiniu.ImgPath + "/" + ret.Key, ret.Key
}

func DeleteFile(key string) error {

	mac := qbox.NewMac(global.GVA_CONFIG.Qiniu.AccessKey, global.GVA_CONFIG.Qiniu.SecretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	// cfg.Zone=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(mac, &cfg)
	err := bucketManager.Delete(global.GVA_CONFIG.Qiniu.Bucket, key)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
`

var directory = `package utils

import (
	"{{.Appname}}/global"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			global.GVA_LOG.Debug("create directory ", v)
			err = os.MkdirAll(v, os.ModePerm)
			if err != nil {
				global.GVA_LOG.Error("create directory", v, " error:", err)
			}
		}
	}
	return err
}
`

var des = `package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

func padding(src []byte, blocksize int) []byte {
	n := len(src)
	padnum := blocksize - n%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	dst := append(src, pad...)
	return dst
}

func unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	dst := src[:n-unpadnum]
	return dst
}

func EncryptDES(src []byte) []byte {
	key := []byte("qimiao66")
	block, _ := des.NewCipher(key)
	src = padding(src, block.BlockSize())
	blockmode := cipher.NewCBCEncrypter(block, key)
	blockmode.CryptBlocks(src, src)
	return src
}

func DecryptDES(src []byte) []byte {
	key := []byte("qimiao66")
	block, _ := des.NewCipher(key)
	blockmode := cipher.NewCBCDecrypter(block, key)
	blockmode.CryptBlocks(src, src)
	src = unpadding(src)
	return src
}
`
var array = `package utils

import (
	"fmt"
	"strings"
)

func ArrayToString(array []interface{}) string {
	return strings.Replace(strings.Trim(fmt.Sprint(array), "[]"), " ", ",", -1)
}
`

