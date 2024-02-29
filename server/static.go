package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"server/util"
	"strings"
)

type UploadResponse struct {
	Url string `json:"url"`
}

func (server *Server) Upload(ctx context.Context, patten string, w http.ResponseWriter, r *http.Request) {
	var err error = nil
	fileName := util.RandomNumLetterStr(8)
	ss := strings.Split(r.URL.Path, patten)
	if len(ss) != 2 {
		panic("not support file name")
	}
	fileName += "." + ss[1]
	filePath := fmt.Sprintf(server.StaticFileSystemPathRoot, fileName)
	defer func() {
		codeMsgData := &CodeMsgResponse{}
		if err != nil {
			codeMsgData.Code = -1
			codeMsgData.Msg = err.Error()
		} else {
			codeMsgData.Data = &UploadResponse{
				Url: fmt.Sprintf(server.StaticHttpRootPath, fileName),
			}
		}
		responseByte, err := json.Marshal(codeMsgData)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		log.Printf("Upload done 1 %s %s %s", r.URL.Path, string(responseByte), ctx.Value(CTX_ID_STR))
		w.Write(responseByte)
	}()
	reader := r.Body
	buff := make([]byte, 1024)
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Server Upload Create %s %s", fileName, filePath)
		return
	}
	defer file.Close()
	for {
		n, err := reader.Read(buff)
		if err != nil && err != io.EOF {
			log.Printf("Upload Read err")
			return
		}
		if err == io.EOF {
			err = nil
			log.Printf("Upload Read done")
			return
		}
		n, err = file.Write(buff[0:n])
		if err != nil {
			log.Printf("Upload Write err")
			return
		}
	}
}

func (server *Server) Download(ctx context.Context, patten string, w http.ResponseWriter, r *http.Request) {
	//todo 一些权限验证
	ss := strings.Split(r.URL.Path, patten)
	if len(ss) != 2 {
		panic("not support file name")
	}
	filePath := fmt.Sprintf(server.StaticFileSystemPathRoot, ss[1])
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	if stat.IsDir() {
		panic("not support dir")
	}
	buff := make([]byte, 4096)
	for {
		n, err := file.Read(buff)
		if err != nil {
			break
		}
		w.Write(buff[0:n])
	}
}
