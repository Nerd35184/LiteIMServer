package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"server/util"
)

const ()

type EmptyRequestBody struct {
}
type EmptyResponseBody struct {
}

type CodeMsgResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewCodeMsgResponse(code int, msg string, data interface{}) *CodeMsgResponse {
	return &CodeMsgResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

type JsonHttpHandler[RequstBody any, ResponseData any] func(
	ctx context.Context,
	body *RequstBody) (*ResponseData, error)

func ProcessHttpJsonRequest[RequstBody any, ResponseData any](ctx context.Context, w http.ResponseWriter, r *http.Request, handler JsonHttpHandler[RequstBody, ResponseData]) *CodeMsgResponse {
	var body RequstBody
	var bodyInterface interface{} = body
	_, ok := bodyInterface.(EmptyRequestBody)
	if !ok {
		bodyByte, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return NewCodeMsgResponse(-1, err.Error(), nil)
		}
		log.Printf("ProcessHttpJsonRequest %s %s %s", r.URL.Path, string(bodyByte), ctx.Value(CTX_ID_STR))
		err = json.Unmarshal(bodyByte, &body)
		if err != nil {
			return NewCodeMsgResponse(-2, err.Error(), nil)
		}
	}
	responseData, err := handler(ctx, &body)
	if err != nil {
		codeErr, ok := err.(*util.CodeError)
		if ok {
			return NewCodeMsgResponse(codeErr.Code, codeErr.Msg, nil)
		}
		return NewCodeMsgResponse(-3, err.Error(), nil)
	}
	return NewCodeMsgResponse(0, "", responseData)
}

func RegisterJsonHttpHandleFunc[RequstBody any, ResponseData any](server *Server, pattern string, auth bool, handler JsonHttpHandler[RequstBody, ResponseData]) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctxId := util.RandomNumLetterStr(8)
		ctx := context.WithValue(r.Context(), CTX_ID_STR, ctxId)
		if auth {
			token := r.Header.Get(HTTP_HEADER_AUTH)
			userId, ok := server.Token2UserId[token]
			if !ok {
				log.Printf("ProcessHttpJsonRequest token not found %s %s", r.URL.Path, ctx.Value(CTX_ID_STR))
				w.Write([]byte("token not found"))
				return
			}

			ctx = context.WithValue(ctx, CTX_USER_ID_KEY, userId)
		}
		codeMsgResponse := ProcessHttpJsonRequest(ctx, w, r, handler)
		responseByte, err := json.Marshal(codeMsgResponse)
		if err != nil {
			log.Printf("ProcessHttpJsonRequest done 1 %s %s %s", r.URL.Path, string(responseByte), ctx.Value(CTX_ID_STR))
			w.Write([]byte(err.Error()))
			return
		}
		log.Printf("ProcessHttpJsonRequest done 2 %s %s %s", r.URL.Path, string(responseByte), ctx.Value(CTX_ID_STR))
		w.Write(responseByte)
	})
}

type DataHttpHandler func(ctx context.Context, patten string, w http.ResponseWriter, r *http.Request)

func RegisterDataHttpHandleFunc(server *Server, pattern string, auth bool, handler DataHttpHandler) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctxId := util.RandomNumLetterStr(8)
		ctx := context.WithValue(r.Context(), CTX_ID_STR, ctxId)
		log.Printf("ProcessHttpDataRequest %s %s", r.URL.Path, ctx.Value(CTX_ID_STR))
		if auth {
			token := r.Header.Get(HTTP_HEADER_AUTH)
			log.Printf("ProcessHttpDataRequest 1 %s %s", util.ToJsonStr(server.Token2UserId), token)
			userId, ok := server.Token2UserId[token]
			if !ok {
				log.Printf("ProcessHttpJsonRequest token not found %s %s %s", r.URL.Path, ctx.Value(CTX_ID_STR), util.ToJsonStr(server.Token2UserId))
				w.Write([]byte("token not found"))
				return
			}

			ctx = context.WithValue(ctx, CTX_USER_ID_KEY, userId)
		}
		handler(ctx, pattern, w, r)
		log.Printf("ProcessHttpDataRequest done 2 %s %s", r.URL.Path, ctx.Value(CTX_ID_STR))
	})
}
