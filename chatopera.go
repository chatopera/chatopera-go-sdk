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
	"io"
	"math/rand"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const SaasAPI = "https://bot.chatopera.com"

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
	appID   string
	sercet  string
	saasAPI string
}

// 聊天机器人实例
func Chatbot(args ...interface{}) *Chatopera {
	result := new(Chatopera)
	result.saasAPI = SaasAPI

	for i, p := range args {
		switch i {
		case 0:
			param, ok := p.(string)
			if !ok {
				panic("1st parameter not type string.")
			}
			result.appID = param

		case 1:
			param, ok := p.(string)
			if !ok {
				panic("2nd parameter not type string.")
			}
			result.sercet = param

		case 2:
			param, ok := p.(string)
			if !ok {
				panic("2nd parameter not type string.")
			}
			result.saasAPI = param

		default:
			panic("Wrong parameter count.")
		}
	}

	return result
}

type Res struct {
	RC          int         `json:"rc"`
	Total       int32       `json:"total"`
	CurrentPage int32       `json:"current_page"`
	TotalPage   int32       `json:"total_page"`
	Error       string      `json:"error"`
	Msg         string      `json:"msg"`
	Status      interface{} `json:"status"`
	Data        interface{} `json:"data"`
}

func httpClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

// 发送请求工具方法
func (c *Chatopera) request(method string, path string, body io.Reader) (*Res, error) {
	t, err := generate(c.appID, c.sercet, method, path)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	req, err := http.NewRequest(method, c.saasAPI+path, body)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	req.Header.Add("Authorization", t)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	r, err := httpClient().Do(req)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	defer r.Body.Close()

	res := new(Res)

	err = json.NewDecoder(r.Body).Decode(res)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	return res, nil
}

func (c *Chatopera) command(method string, path string, payloads ...interface{}) (*Res, error) {
	path = "/api/v1/chatbot/" + c.appID + path

	if strings.Contains(path, "?") {
		path += "&sdklang=go"
	} else {
		path += "?sdklang=go"
	}

	var reader io.Reader = nil
	if payloads != nil {
		json, _ := json.Marshal(payloads[0])
		reader = bytes.NewBuffer(json)
	}
	res, err := c.request(method, path, reader)

	return res, err
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
func (c *Chatopera) Detail() (*Res, error) {
	return c.command("GET", "/")
}

type ConversationBody struct {
	FromUserID  string `json:"fromUserId"`
	TextMessage string `json:"textMessage"`
	IsDebug     bool   `json:"isDebug"`
}

// 检索多轮对话
func (c *Chatopera) Conversation(userID string, textMessage string) (*Res, error) {
	body := ConversationBody{
		FromUserID:  userID,
		TextMessage: textMessage,
		IsDebug:     false,
	}

	return c.command("POST", "/conversation/query", body)
}

type FaqBody struct {
	FromUserID string `json:"fromUserId"`
	Query      string `json:"query"`
}

// 检索知识库
func (c *Chatopera) Faq(userID string, query string) (*Res, error) {
	body := FaqBody{
		FromUserID: userID,
		Query:      query,
	}

	return c.command("POST", "/faq/query", body)
}

// 获得用户列表
func (c *Chatopera) Users(limit int, page int, sortby string) (*Res, error) {
	return c.command("GET", "/users?page="+strconv.Itoa(page)+"&limit="+strconv.Itoa(limit)+"&sortby="+sortby)
}

// 获得聊天历史
func (c *Chatopera) Chats(userID string, limit int, page int, sortby string) (*Res, error) {
	return c.command("GET", "/users/"+userID+"/chats?page="+strconv.Itoa(page)+"&limit="+strconv.Itoa(limit)+"&sortby="+sortby)
}

// 屏蔽用户
func (c *Chatopera) Mute(userID string) (*Res, error) {
	return c.command("POST", "/users/"+userID+"/mute")
}

// 取消屏蔽
func (c *Chatopera) Unmute(userID string) (*Res, error) {
	return c.command("POST", "/users/"+userID+"/unmute")
}

// 用户是否被屏蔽
func (c *Chatopera) Ismute(userID string) (*Res, error) {
	return c.command("POST", "/users/"+userID+"/ismute")
}

func (c *Chatopera) User(userID string) (*Res, error) {
	return c.command("GET", "/users/"+userID+"/profile")
}

type IntentSessionBody struct {
	UID     string `json:"uid"`
	Channel string `json:"channel"`
}

func (c *Chatopera) IntentSession(uid string, channel string) (*Res, error) {
	body := IntentSessionBody{
		UID:     uid,
		Channel: channel,
	}
	return c.command("POST", "/clause/prover/session", body)
}

func (c *Chatopera) IntentSessionDetail(sessionId string) (*Res, error) {
	return c.command("GET", "/clause/prover/session/"+sessionId)
}

type IntentChatBodySession struct {
	ID string `json:"id"`
}

type IntentChatBodyMessage struct {
	TextMessage string `json:"textMessage"`
}

type IntentChatBody struct {
	FromUserID string                `json:"fromUserId"`
	Session    IntentChatBodySession `json:"session"`
	Message    IntentChatBodyMessage `json:"message"`
}

type IntentMessage struct {
	Receiver    string `json:"receiver"`
	IsProactive bool   `json:"is_proactive"`
	IsFallback  bool   `json:"is_fallback"`
	TextMessage string `json:"textMessage"`
}

func (c *Chatopera) IntentChat(sessionId string, uid string, textMessage string) (*Res, error) {
	body := IntentChatBody{
		FromUserID: uid,
		Session: IntentChatBodySession{
			ID: sessionId,
		},
		Message: IntentChatBodyMessage{
			TextMessage: textMessage,
		},
	}
	return c.command("POST", "/clause/prover/chat", body)
}
