// Copyright 2018 Chatopera Inc. All rights reserved.
// This software and related documentation are provided under a license agreement containing
// restrictions on use and disclosure and are protected by intellectual property laws.
// Except as expressly permitted in your license agreement or allowed by law, you may not use,
// copy, reproduce, translate, broadcast, modify, license, transmit, distribute, exhibit, perform,
// publish, or display any part, in any form, or by any means. Reverse engineering, disassembly,
// or decompilation of this software, unless required by law for interoperability, is prohibited.
package chatopera

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	SaasAPI = "https://bot.chatopera.com"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type SignatureBody struct {
	AppID     string `json:"appId"`
	Timestamp string `json:"timestamp"`
	Random    string `json:"random"`
	Signature string `json:"signature"`
}

func generate(appID string, secret string, method string, path string) (string, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	random := randStringBytes(24)

	key := []byte(secret)
	mac := hmac.New(sha1.New, key)
	p := appID + timestamp + random + method + path
	mac.Write([]byte(p))
	result := mac.Sum(nil)
	signature := hex.EncodeToString(result)

	body := SignatureBody{
		AppID:     appID,
		Timestamp: timestamp,
		Random:    random,
		Signature: signature,
	}

	data, err := json.Marshal(body)
	if err == nil {
		return base64.StdEncoding.EncodeToString(data), nil
	}
	return "", err
}

type Chatopera struct {
	appID  string
	sercet string
}

// 聊天机器人实例
func Chatbot(appID string, sercet string) *Chatopera {
	result := new(Chatopera)
	result.appID = appID
	result.sercet = sercet
	return result
}

type Res struct {
	RC          int         `json:"rc"`
	Total       int32       `json:"total"`
	CurrentPage int32       `json:"current_page"`
	TotalPage   int32       `json:"total_page"`
	Error       string      `json:"error"`
	Data        interface{} `json:"data"`
}

func httpClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Second}
}

// 发送请求工具方法
func request(appID string, sercet string, method string, path string, body io.Reader, v interface{}) (int32, int32, int32, error) {
	t, err := generate(appID, sercet, method, path)
	if err != nil {
		return 0, 0, 0, err
	}

	req, err := http.NewRequest(method, SaasAPI+path, body)
	if err != nil {
		return 0, 0, 0, err
	}
	req.Header.Add("Authorization", t)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	r, err := httpClient().Do(req)
	if err != nil {
		return 0, 0, 0, err
	}
	defer r.Body.Close()

	res := new(Res)
	res.Data = v

	err = json.NewDecoder(r.Body).Decode(res)
	if err != nil {
		return 0, 0, 0, err
	}

	if res.RC != 0 {
		return 0, 0, 0, errors.New(res.Error)
	}

	return res.Total, res.CurrentPage, res.TotalPage, nil
}

// 机器人详情
type BotDetail struct {
	Name            string `json:"name"`
	Fallback        string `json:"fallback"`
	Description     string `json:"description"`
	Welcome         string `json:"welcome"`
	PrimaryLanguage string `json:"primaryLanguage"`
}

// 获得机器人详情
func (c *Chatopera) Detail() (*BotDetail, error) {
	path := "/api/v1/chatbot/" + c.appID
	result := new(BotDetail)

	_, _, _, err := request(c.appID, c.sercet, "GET", path, nil, result)
	return result, err
}

// 多轮对话回复
type ConversationResult struct {
	State           string      `json:"state"`
	CreatedAt       int64       `json:"createdAt"`
	String          string      `json:"string"`
	TopicName       string      `json:"topicName"`
	SubReplies      interface{} `json:"subReplies"`
	LogicIsFallback bool        `json:"logic_is_fallback"`
	BotName         string      `json:"botName"`
	Service         interface{} `json:"service"`
}

type ConversationBody struct {
	FromUserID  string `json:"fromUserId"`
	TextMessage string `json:"textMessage"`
	IsDebug     bool   `json:"isDebug"`
}

// 检索多轮对话
func (c *Chatopera) Conversation(userID string, textMessage string) (*ConversationResult, error) {
	path := "/api/v1/chatbot/" + c.appID + "/conversation/query"
	result := new(ConversationResult)

	body := ConversationBody{
		FromUserID:  userID,
		TextMessage: textMessage,
		IsDebug:     false,
	}

	json, _ := json.Marshal(body)

	_, _, _, err := request(c.appID, c.sercet, "POST", path, bytes.NewBuffer(json), result)
	return result, err
}

// 知识库回复
type FaqResult struct {
	ID    string  `json:"id"`
	Score float32 `json:"score"`
	Post  string  `json:"post"`
	Reply string  `json:"reply"`
}

type FaqBody struct {
	FromUserID string `json:"fromUserId"`
	Query      string `json:"query"`
}

// 检索知识库
func (c *Chatopera) Faq(userID string, query string) (*[]FaqResult, error) {
	path := "/api/v1/chatbot/" + c.appID + "/faq/query"
	result := new([]FaqResult)

	body := FaqBody{
		FromUserID: userID,
		Query:      query,
	}

	json, _ := json.Marshal(body)

	_, _, _, err := request(c.appID, c.sercet, "POST", path, bytes.NewBuffer(json), result)
	return result, err
}

// 用户信息
type UsersResult struct {
	UserID   string `json:"userId"`
	Lasttime string `json:"lasttime"`
	Created  string `json:"created"`
}

// 获得用户列表
func (c *Chatopera) Users(limit int, page int, sortby string) (int32, int32, int32, *[]UsersResult, error) {
	path := "/api/v1/chatbot/" + c.appID + "/users?page=" + strconv.Itoa(page) + "&limit=" + strconv.Itoa(limit) + "&sortby=" + sortby
	result := new([]UsersResult)

	total, currentPage, totalPage, err := request(c.appID, c.sercet, "GET", path, nil, result)
	return total, currentPage, totalPage, result, err
}

// 聊天历史
type ChatsResult struct {
	UserID      string `json:"userId"`
	TextMessage string `json:"textMessage"`
	Direction   string `json:"direction"`
	Created     string `json:"created"`
}

// 获得聊天历史
func (c *Chatopera) Chats(userID string, limit int, page int, sortby string) (int32, int32, int32, *[]ChatsResult, error) {
	path := "/api/v1/chatbot/" + c.appID + "/users/" + userID + "/chats?page=" + strconv.Itoa(page) + "&limit=" + strconv.Itoa(limit) + "&sortby=" + sortby
	result := new([]ChatsResult)

	total, currentPage, totalPage, err := request(c.appID, c.sercet, "GET", path, nil, result)
	return total, currentPage, totalPage, result, err
}

// 屏蔽用户
func (c *Chatopera) Mute(userID string) error {
	path := "/api/v1/chatbot/" + c.appID + "/users/" + userID + "/mute"
	_, _, _, err := request(c.appID, c.sercet, "POST", path, nil, nil)
	return err
}

// 取消屏蔽
func (c *Chatopera) Unmute(userID string) error {
	path := "/api/v1/chatbot/" + c.appID + "/users/" + userID + "/unmute"
	_, _, _, err := request(c.appID, c.sercet, "POST", path, nil, nil)
	return err
}

type IsmuteResult struct {
	Mute bool `json:"mute"`
}

// 用户是否被屏蔽
func (c *Chatopera) Ismute(userID string) (bool, error) {
	path := "/api/v1/chatbot/" + c.appID + "/users/" + userID + "/ismute"
	result := new(IsmuteResult)
	_, _, _, err := request(c.appID, c.sercet, "POST", path, nil, result)
	return result.Mute, err
}

// 用户画像
type UserResult struct {
	UserID   string `json:"userId"`
	Name     string `json:"name"`
	Lasttime string `json:"lasttime"`
	Mute     bool   `json:"mute"`
}

func (c *Chatopera) User(userID string) (*UserResult, error) {
	path := "/api/v1/chatbot/" + c.appID + "/users/" + userID + "/profile"
	result := new(UserResult)
	_, _, _, err := request(c.appID, c.sercet, "GET", path, nil, result)
	return result, err
}
