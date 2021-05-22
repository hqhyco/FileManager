package routes

import (
	"ResManager/src/db"
	"ResManager/src/global"
	"ResManager/src/models"
	"ResManager/src/session"
	"ResManager/src/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var storage = "/files" //上传文件目录
var curFilePath = ""   //当前目录地址

var pn = global.CFG.Settings.Pn //每页数量
var ps = 1 //当前页数

//https://www.cnblogs.com/merrynuts/p/13157381.html

func init() {

}

func RenderHome(w http.ResponseWriter, r *http.Request, title string) {

	query := r.URL.Query()

	_, _ = utils.PathExists("." + storage)
	files, pages := findFile(query, r)

	data := models.Page{}
	user, _ := session.GetSession(w, r).GetAttr("user")

	if user != nil {
		data.Title = title
		data.UserID = user.(*models.User).ID
		data.Username = user.(*models.User).Username
	}
	data.File = files
	data.PageList = pages

	t, err := template.ParseFiles("templates/index.html")
	utils.CheckError(err)

	err = t.Execute(w, data)
	utils.CheckError(err)
}

func RenderLogin(w http.ResponseWriter, r *http.Request, title string) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("templates/login.html")
		utils.CheckError(err)
		err = t.Execute(w, nil)
		utils.CheckError(err)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Println("login", username, password)
	if utils.IsEmpty(username, password) {
		message(w, r, "输入不能为空")
		return
	}

	user := db.FindUserByUsernameAndPassword(username, password)
	if user == nil {
		message(w, r, "用户名或密码错误，登录失败！")
		return
	}
	// 登陆成功
	sess := session.GetSession(w, r)
	sess.SetAttr("user", user)
	http.Redirect(w, r, getDirectLink(ps), 302)
}

func RenderLogout(w http.ResponseWriter, r *http.Request, title string) {
	sess := session.GetSession(w, r)
	sess.DelAttr("user")
	http.Redirect(w, r, getDirectLink(ps), 302)
}

func RenderAdmin(w http.ResponseWriter, r *http.Request, title string) {
	sess := session.GetSession(w, r)
	user, exist := sess.GetAttr("user")
	if !exist {
		http.Redirect(w, r, getDirectLink(ps), 302)
		return
	}

	if r.Method == "GET" {
		t, err := template.ParseFiles("templates/userinfo.html")
		utils.CheckError(err)
		err = t.Execute(w, user)
		utils.CheckError(err)
		return
	}

	// POST 更新用户信息
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if utils.IsEmpty(username, password, email) {
		message(w, r, "字段不能为空")
		return
	}

	switch user := user.(type) {
	case *models.User:
		user.Username = username
		user.Password = password
		user.Email = email
		db.UpdateUser(user)
	default:
		log.Println(":userinfo:user.(type)", user)
	}
	http.Redirect(w, r, "/admin", 302)
}

func RenderUploadFile(w http.ResponseWriter, r *http.Request, title string) {
	//在使用r.MultipartForm前必须先调用ParseMultipartForm方法，参数为最大缓存
	_ = r.ParseMultipartForm(32 << 20)
	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		//进入存储目录
		//os.Chdir(storage)
		//退出存储目录
		//defer os.Chdir(currentDirectory)
		//获取所有上传文件信息
		fhs := r.MultipartForm.File["userfile"]
		num := len(fhs)
		fmt.Printf("总文件数：%d 个文件", num)
		//循环对每个文件进行处理
		for n, fheader := range fhs {
			//设置文件名
			//newFileName := strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(fheader.Filename) //十进制
			dir := "." + storage + fmt.Sprintf("/%d%02d/", time.Now().Year(), time.Now().Month())
			//fmt.Println(dir)
			exist, _ := utils.PathExists(dir)
			if exist == false {
				fmt.Println("文件夹创建失败")
				return
			}
			newFileName := dir + fheader.Filename
			//打开上传文件
			uploadFile, err := fheader.Open()
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer uploadFile.Close()
			//保存文件
			saveFile, err := os.Create(newFileName)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer saveFile.Close()
			_, _ = io.Copy(saveFile, uploadFile)

			//获取文件状态信息
			fileStat, _ := saveFile.Stat()
			//打印接收信息
			fmt.Printf("%s  NO.: %d  Size: %d KB  Name：%s\n", time.Now().Format("2006-01-02 15:04:05"), n, fileStat.Size()/1024, newFileName)

		}
		fmt.Println("上传成功")
	}
}

func RenderListData(w http.ResponseWriter, r *http.Request, title string) {
	if r.Method == "GET" {
		query := r.URL.Query()

		if len(query) <= 0 {
			message(w, r, "未发现记录")
			return
		}

		path, ok1 := query["fp"]
		file, ok2 := query["fn"]

		if ok1 && ok2 {
			message(w, r, "参数错误")
			return
		}

		if ok1 && len(path) == 1 {
			//打开文件夹
			if curFilePath == path[0] {
				http.Redirect(w, r, "/", 302)
				return
			}
			curFilePath = path[0]
			http.Redirect(w, r, getDirectLink(ps), 302)
		}

		if ok2 && len(file) == 1 {
			//查看文件

		}
	} else if r.Method == "POST" {

		fmt.Println(r.Header.Get("Content-Type"))

		var post_data map[string]interface{}
		body, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(body, &post_data)

		fname := fmt.Sprintf("./%s/%s/%s", storage, curFilePath, post_data["del"])

		if delFile(fname) {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}

}

func message(w http.ResponseWriter, r *http.Request, message string) {
	t, err := template.ParseFiles("templates/message.html")
	utils.CheckError(err)

	err = t.Execute(w, map[string]string{"Message": message})
	utils.CheckError(err)
}

/**
获取文件列表
*/
func findFile(v url.Values, r *http.Request) (files []models.FileInfo, pages []models.PageListInfo) {

	fp := v["fp"]
	_ps := v["ps"]

	if len(fp) > 0 && fp[0] != "" {
		curFilePath = fp[0]
	} else {
		curFilePath = ""
	}

	if len(_ps) > 0 && _ps[0] != ""{
		ps, _ = strconv.Atoi(_ps[0])
		if ps < 1{
			ps = 1
		}
	}

	if curFilePath == "" {
		//创建目录数据
		mDirs, _ := listDir("." + storage)

		for _, f := range mDirs {
			file := models.FileInfo{
				Name:  f.Name(),
				IsDir: f.IsDir(),
				Size:  "",
				Date:  f.ModTime().Format("2006-01-02 15:04"),
				Link:  "-",
				Type:  "fp",
				Ps: 1,
			}
			files = append(files, file)
		}
	} else {
		//查找文件夹是否存在
		joinPath := "." + storage + "/" + curFilePath
		if !isDir(joinPath) {
			return nil, nil
		}
		//遍历文件
		_, mFiles := listDir(joinPath)

		//添加父文件夹
		files = append(files, models.FileInfo{
			Name:  curFilePath,
			IsDir: true,
			Size:  "",
			Date:  "-",
			Type:  "fp",
			Ps: ps,
		})

		for _, f := range mFiles {
			file := models.FileInfo{
				Name:  f.Name(),
				IsDir: f.IsDir(),
				Size:  utils.GetFileSize(float64(f.Size())),
				Date:  f.ModTime().Format("2006-01-02 15:04"),
				Link:  getWebUrl(r) + storage + "/" + curFilePath + "/" + f.Name(),
				Type:  "fn",
			}
			files = append(files, file)
		}

		//获取文件列表
		low := (ps-1) * pn
		high := ps * pn
		file_len := len(files)
		if high > file_len{high = file_len}
		if low > file_len{low=file_len-1}
		files = files[low:high]

		pages = getPageList(file_len)
	}

	return files, pages
}

func listDir(path string) (dirs []os.FileInfo, files []os.FileInfo) {
	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirs = append(dirs, info)
		} else {
			files = append(files, info)
		}
		return nil
	})

	return dirs, files
}

func delFile(fname string) (state bool) {
	if !utils.CheckFileIsExist(fname) {
		return false
	}

	//删除
	err := os.Remove(fname)

	if err != nil {
		return false
	} else {
		return true
	}
}

func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func isFile(path string) bool {
	return !isDir(path)
}

func getWebUrl(r *http.Request) (url string) {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	return strings.Join([]string{scheme, r.Host}, "")
}

func getDirectLink(p int) (url string) {
	return fmt.Sprintf("/?fp=%s&ps=%d", curFilePath, p)
}

func getPageList(file_len int) (list []models.PageListInfo) {
	list_count := 5

	pss := math.Ceil(float64(file_len)/ float64(pn))
	ps = int(math.Ceil(math.Min(float64(ps), pss)))

	//计算分页信息
	var start int
	var stop int
	if ps > list_count {
		start = int(math.Max(float64((ps) - 2), 1))
		stop = int(math.Min(float64(ps + 2), pss))
	}else{
		start = 1
		stop = int(math.Min(float64(list_count), pss))
	}

	//添加前一页
	prev := int(math.Max(float64(ps - 1), 1))
	if prev >= list_count{
		info_prev := models.PageListInfo{
			PUrl: getDirectLink(1),  //首页
			PNum: "<<",
		}
		list = append(list, info_prev)
	}
	//添加中间页列表
	for i:=start;i<=stop;i++{
		info := models.PageListInfo{
			PUrl: getDirectLink(i),
			PNum: fmt.Sprintf("%d", i),
			Active: i == ps,
		}
		list = append(list, info)
	}
	//添加后一页
	next := int(math.Min(float64(ps + 1), pss))
	if ps < int(pss)-1{
		info_next := models.PageListInfo{
			PUrl: getDirectLink(next),
			PNum: ">>",
		}
		list = append(list, info_next)
	}

	//fmt.Println(list)
	return list
}
