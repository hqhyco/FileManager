package models

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

type PageListInfo struct {
	PUrl string
	PNum string
	Active bool
}

//File meta
type FileInfo struct {
	Name  string
	Size  string
	Date  string
	Link  string
	IsDir bool
	Type  string
	Ps    int
}

// page meta
type Page struct {
	UserID   int
	Username string
	Title    string
	Account  interface{}
	CSS      []string
	JS       []string
	Message  []string
	//Files []FileInfo
	File []FileInfo
	PageList []PageListInfo
}
