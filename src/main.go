package main

import (
	"ResManager/src/global"
	"ResManager/src/routes"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var port = fmt.Sprintf(":%d",global.CFG.Settings.Port)
var router = mux.NewRouter()

func indexHandler(w http.ResponseWriter, r *http.Request) {
	routes.RenderHome(w, r, "index")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	routes.RenderLogin(w, r, "login")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	routes.RenderLogout(w, r, "logout")
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	routes.RenderAdmin(w, r, "admin")
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	routes.RenderUploadFile(w, r, "upload")
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	routes.RenderListData(w, r, "list")
}

func init() {
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/logout", logoutHandler)
	router.HandleFunc("/admin", accountHandler)
	router.HandleFunc("/upload", uploadFileHandler)
	router.HandleFunc("/list", listHandler)
	router.PathPrefix("/static/css/").Handler(http.StripPrefix("/static/css/", http.FileServer(http.Dir("static/css/"))))
	router.PathPrefix("/static/js/").Handler(http.StripPrefix("/static/js/", http.FileServer(http.Dir("static/js/"))))
	router.PathPrefix("/static/images/").Handler(http.StripPrefix("/static/images/", http.FileServer(http.Dir("static/images/"))))
	router.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir("files/"))))
}

func main() {
	//启动服务
	log.Printf("Server is running at http://localhost%s/. Press Ctrl+C to stop.", port)

	s := &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    3600 * time.Second,
		WriteTimeout:   3600 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
