package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type ApiResp struct {
	StatusCode int         `json:"-"`
	OK         bool        `json:"ok"`
	Data       interface{} `json:"data,omitempty"`
	ErrCode    string      `json:"err,omitempty"`
	Message    string      `json:"msg,omitempty"`
	Timestamp  int64       `json:"ts"`
}

func NewSuccessResp(data interface{}) ApiResp {
	return ApiResp{
		StatusCode: http.StatusOK,
		OK:         true,
		Data:       data,
		Timestamp:  time.Now().Unix(),
	}
}

func NewErrorResp(err error) ApiResp {
	var e *Error
	if !errors.As(err, &e) {
		e = NewErrInternalError(err)
	}
	return ApiResp{
		StatusCode: e.StatusCode,
		OK:         false,
		ErrCode:    e.ErrCode,
		Message:    e.Message,
		Timestamp:  time.Now().Unix(),
	}
}

func WriteAPIResp(w http.ResponseWriter, resp ApiResp) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	b, _ := json.Marshal(resp)
	w.Write(b)
}

type analyzeReqBody struct {
	InputText string `json:"input_text"`
	RefText   string `json:"ref_text"`
}

func (rb analyzeReqBody) Validate() error {
	if len(rb.InputText) == 0 {
		return NewErrBadRequest("Source Text cannot be empty!")
	}
	if len(rb.RefText) == 0 {
		return NewErrBadRequest("Target Text cannot be empty!")
	}
	return nil
}

type match struct {
	Input     matchDetails `json:"input"`
	Reference matchDetails `json:"ref"`
}

type matchDetails struct {
	Text     string `json:"text"`
	StartIdx int    `json:"start_idx"`
	EndIdx   int    `json:"end_idx"`
}
