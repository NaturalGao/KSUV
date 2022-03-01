package main

import (
	"KsUploadVideo/api"
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math/big"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/eiannone/keyboard"

	"github.com/go-resty/resty/v2"
)

const FILE_URL string = "./config.json"
const MAX_UPLOAD_VIDEO = 100

var chSignal chan os.Signal

var Settings Config

var Titles []string

var VideoFiles []fs.FileInfo

var ClientHttp *resty.Client

var inputReader *bufio.Reader

var Ao *api.ApiObject

var Commodities []interface{}

var SelectCommodIndex int

type UserConfig struct {
	Cookie   string
	WebApiPh string
}

type Config struct {
	Version         string
	Name            string
	Authors         string
	UserConfig      *UserConfig
	TitleFileUrl    string
	VideoFileUrl    string
	SecondDomain    string
	UseSerialNumber bool
}

type ErrorD struct {
	err string
}

func (e *ErrorD) Error() string {
	return e.err
}

func main() {

	fmt.Printf("************************** 欢迎使用 %s/%s **************************\n", Settings.Name, Settings.Version)

	fmt.Printf("程序已启动就绪：\n1. enter 开始执行程序\n2. Esc 终止程序\n")

	isStart := isStart()

	if !isStart {
		fmt.Println("程序即将终止....")
		os.Exit(1)
	}

	fmt.Println("程序开始执行.......")

	Ao = api.New(Settings.UserConfig.Cookie)

	isSelected := selectCommoditie()

	if !isSelected {
		fmt.Println("程序即将终止....")
		os.Exit(1)
	}

	uploadRun()

	//合建chan
	chSignal = make(chan os.Signal, 1)
	//监听指定信号 ctrl+c kill
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)

	<-chSignal

	fmt.Println("************************** 程序执行完毕 **************************")
}

// 上传任务开始
func uploadRun() {
	s_c, e_c := 0, 0

	tl := int64(len(Titles))

	var wg sync.WaitGroup
	ch := make(chan struct{}, 10)

	for i, file := range VideoFiles {
		ch <- struct{}{}
		wg.Add(1)

		go func(file fs.FileInfo, i int) {

			defer wg.Done()

			fileName := file.Name()
			fileLength := file.Size()

			upToken, err := UploadToken()

			if err != nil {
				e_c++
				<-ch
				return
			}

			isUpload := UploadMultipart(upToken, fileName, fileLength)

			if !isUpload {
				e_c++
				<-ch
				return
			}

			fileInfo, err := UploadFinish(upToken, fileName, fileLength)

			if err != nil {
				e_c++
				<-ch
				return
			}

			isSubmit, isEnd := SubmitVideo(fileInfo, fileName, i, tl)

			if isEnd {
				os.Exit(1)
			}

			if !isSubmit {
				e_c++
				<-ch
				return
			}

			s_c++
			<-ch
		}(file, i)
	}

	wg.Wait()

	fmt.Println("视频发布完成!")
	fmt.Printf("success: %d  error: %d\n", s_c, e_c)
}

// 获取视频上传Token
func UploadToken() (token string, err error) {
	fmt.Println("获取视频上传凭证中...")

	body := map[string]interface{}{
		"kuaishou.web.cp.api_ph": Settings.UserConfig.WebApiPh,
		"uploadType":             1,
	}

	resp, err := Ao.UploadToken(body)

	if err != nil {
		fmt.Println("获取上传视频凭证失败!")
		return token, err
	}

	if resp["result"] != float64(1) {
		fmt.Println("获取上传视频凭证失败!")
		return token, err
	}

	data := (resp["data"]).(map[string]interface{})

	return (data["token"]).(string), nil
}

// 上传视频
func UploadMultipart(upToekn string, fileName string, fileLength int64) bool {

	fmt.Printf("上传文件【%s】中,文件大小 %d\n", fileName, fileLength)

	fileBytes, _ := ioutil.ReadFile(Settings.VideoFileUrl + fileName)
	resp, err := Ao.UploadMultipart(upToekn, fileName, fileBytes)

	if err != nil {
		fmt.Printf("文件【%s】上传失败！原因：%s\n", fileName, err)
		return false
	}

	if resp["result"] != float64(1) {
		fmt.Printf("文件【%s】上传失败！原因：%s\n", fileName, resp["message"])
		return false
	}

	fmt.Printf("文件【%s】上传成功！\n", fileName)

	return true
}

// 获取远程文件信息
func UploadFinish(upToken string, fileName string, fileLength int64) (api.ResultMp, error) {
	fmt.Printf("获取远程视频文件【%s】信息中,文件大小 %d\n", fileName, fileLength)
	body := map[string]interface{}{
		"token":                  upToken,
		"kuaishou.web.cp.api_ph": Settings.UserConfig.WebApiPh,
		"fileName":               fileName,
		"fileLength":             fileLength,
		"fileType":               "video/mp4",
	}

	resp, err := Ao.UploadFinish(body)

	if err != nil {
		fmt.Printf("获取远程文件【%s】信息失败！原因：%s\n", fileName, err)
		return make(map[string]interface{}), err
	}

	if resp["result"] != float64(1) {
		fmt.Printf("获取远程文件【%s】信息失败！原因：%s\n", fileName, resp["message"])
		e_d := new(ErrorD)
		e_d.err = (resp["message"]).(string)
		return make(map[string]interface{}), e_d
	}

	fileInfo := (resp["data"]).(map[string]interface{})

	return fileInfo, nil
}

// 发布视频
func SubmitVideo(fileInfo api.ResultMp, fileName string, s_n int, tl int64) (isSubmit bool, isEnd bool) {
	fmt.Printf("发布视频【%s】中....\n", fileName)

	commodity := Commodities[SelectCommodIndex]

	commodity_map := (commodity).(map[string]interface{})

	ti, _ := rand.Int(rand.Reader, big.NewInt(tl))

	caption := Titles[ti.Int64()]

	if Settings.UseSerialNumber {
		caption += " " + strconv.Itoa(s_n+1)
	}

	recTagIdList := []int64{33532345,
		22154940,
		16859352,
		13623510,
		22916352,
		29267290,
		14746816,
		15831701,
		2411963,
		18887932}

	body := map[string]interface{}{
		"fileId":                 int((fileInfo["fileId"]).(float64)),
		"coverKey":               fileInfo["coverKey"],
		"kuaishou.web.cp.api_ph": Settings.UserConfig.WebApiPh,
		"onvideoDuration":        fileInfo["duration"],
		"associateTaskId":        commodity_map["associateTaskId"],
		"caption":                caption, //标题,
		"domain":                 "其它",
		"secondDomain":           Settings.SecondDomain, //其它描述
		"coverTimeStamp":         -1,
		"photoStatus":            1,
		"coverType":              1,
		"coverTitle":             "",
		"photoType":              0,
		"collectionId":           "",
		"publishTime":            0,
		"longitude":              "",
		"latitude":               "",
		"notifyResult":           0,
		"coverCropped":           false,
		"pkCoverKey":             "",
		"downloadType":           1,
		"disableNearbyShow":      false,
		"allowSameFrame":         true,
		"movieId":                "",
		"projectId":              "",
		"videoInfoMeta":          "",
		"needDeleteKey":          []string{},
		"triggerH265":            false,
		"recTagIdList":           recTagIdList,
	}

	resp, err := Ao.SubmitVideo(body)

	if err != nil {
		fmt.Printf("发布视频【%s】失败！原因：%s\n", fileName, err)
		return false, false
	}

	if resp["result"] != float64(1) {
		fmt.Printf("发布视频【%s】失败！原因：%s\n", fileName, resp["message"])
		if resp["result"] == float64(30006) {
			return false, true
		}
		return false, false
	}

	fmt.Printf("发布视频【%s】成功！\n", fileName)
	return true, false
}

func isStart() bool {
	if err := keyboard.Open(); err != nil {
		fmt.Println(err)
		return false
	}

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println(err)
			return false
		}

		if key == keyboard.KeyEsc {
			keyboard.Close()
			return false
		} else if key == keyboard.KeyEnter {
			// 程序开始
			keyboard.Close()
			break
		}
	}

	return true
}

func init() {
	fmt.Println("初始化配置中，请稍后......")

	init_status := loadConfig()
	if !init_status {
		os.Exit(1)
	}

	// 加载标题文件``
	title_status := loadTitle()

	if !title_status {
		os.Exit(1)
	}

	load_video_status := loadVideoPath()

	if !load_video_status {
		os.Exit(1)
	}

	fmt.Println("配置文件加载成功，程序启动中,请稍后....")
}

// 初始化配置
func loadConfig() bool {

	data, err := ioutil.ReadFile(FILE_URL)

	if err != nil {
		fmt.Println(FILE_URL, "读取文件失败，程序即将终止....")
		return false
	}

	err = json.Unmarshal(data, &Settings)

	if err != nil {
		fmt.Println("配置文件解析失败，程序即将终止....")
		return false
	}

	if Settings.Name == "" ||
		Settings.Version == "" ||
		Settings.TitleFileUrl == "" ||
		Settings.VideoFileUrl == "" ||
		Settings.SecondDomain == "" ||
		(*Settings.UserConfig) == (UserConfig{}) ||
		(*Settings.UserConfig).Cookie == "" ||
		(*Settings.UserConfig).WebApiPh == "" {
		fmt.Println("配置文件确实必要参数，程序即将终止...")
		return false
	}

	return true
}

func loadTitle() bool {
	data, err := ioutil.ReadFile(Settings.TitleFileUrl)

	if err != nil {
		fmt.Println(Settings.TitleFileUrl, "文件读取失败，程序即将终止....")
		return false
	}

	if len(data) == 0 {
		fmt.Println(Settings.TitleFileUrl, "标题文件为空，程序即将终止....")
		return false
	}

	Titles = strings.Split(strings.Trim(string(data), "\n"), "\n")

	fmt.Println("检测到有", len(Titles), "个标题")

	return true
}

// 加载视频路径
func loadVideoPath() bool {
	files, err := ioutil.ReadDir(Settings.VideoFileUrl)

	if err != nil {
		fmt.Println(Settings.VideoFileUrl, "目录不正确，程序即将终止....")
		return false
	}

	for _, f := range files {
		fileName := f.Name()
		filesNameString := strings.Split(fileName, ".")
		suffix := strings.ToLower(filesNameString[len(filesNameString)-1])
		if suffix == "mp4" {
			VideoFiles = append(VideoFiles, f)
		}

	}

	v_l := len(VideoFiles)

	if v_l == 0 {
		fmt.Println(Settings.VideoFileUrl, "此目录下暂未发现 MP4 格式文件,程序即将终止....")
		return false
	}

	if v_l > MAX_UPLOAD_VIDEO {
		VideoFiles = VideoFiles[0:MAX_UPLOAD_VIDEO]
	}

	fmt.Println("检测到有", v_l, "个视频可上传!上传视频上限", MAX_UPLOAD_VIDEO, "个视频")

	return true
}

// 绑定关联商品
func selectCommoditie() bool {
	body := map[string]interface{}{
		"kuaishou.web.cp.api_ph": Settings.UserConfig.WebApiPh,
		"type":                   1,
		"cursor":                 "",
	}

	r, err := Ao.RelationList(body)

	if err != nil {
		fmt.Println(err)
		return false
	}

	if r["result"] != float64(1) {
		fmt.Println(r["message"])
		return false
	}

	l := (r["data"]).(map[string]interface{})["list"]

	Commodities = l.([]interface{})

	fmt.Println("发现有", len(Commodities), "种商品可选择,请输入对应的编号关联上传商品：")

	titles := []string{}

	for i, v := range Commodities {
		ob := v.(map[string]interface{})

		t := ob["title"]

		titles = append(titles, t.(string))

		fmt.Printf("[%d] %s\n", i, t)
	}

	isInputRead := listenInputReader()

	if isInputRead {
		fmt.Println("您最终选择关联的商品：")
		fmt.Println(titles[SelectCommodIndex])
	}

	return isInputRead
}

func listenInputReader() bool {
	inputReader = bufio.NewReader(os.Stdin)

RESTART:
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return false
	}

	input = strings.Trim(input, "\n\r")

	SelectCommodIndex, err = strconv.Atoi(input)

	if err != nil || SelectCommodIndex < 0 || SelectCommodIndex > len(Commodities)-1 {
		fmt.Println("输入了非法值")
		goto RESTART
	}

	return true
}
