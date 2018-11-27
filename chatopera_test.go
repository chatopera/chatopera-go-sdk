// Copyright 2018 Chatopera Inc. All rights reserved.
// This software and related documentation are provided under a license agreement containing
// restrictions on use and disclosure and are protected by intellectual property laws.
// Except as expressly permitted in your license agreement or allowed by law, you may not use,
// copy, reproduce, translate, broadcast, modify, license, transmit, distribute, exhibit, perform,
// publish, or display any part, in any form, or by any means. Reverse engineering, disassembly,
// or decompilation of this software, unless required by law for interoperability, is prohibited.

package chatopera_test

import (
	"github.com/chatopera/chatopera-go-sdk"
	"testing"
)

var bot = chatopera.Chatbot("5bf27e4d6f80090017b404b7", "e4cbc6a65708c011ec0da73b0f5db7a1")

func TestConversion(t *testing.T) {
	reply, err := bot.Conversation("xiao", "你好")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("TestConversion reply:" + reply.String)
	}
}
func TestFaq(t *testing.T) {
	_, err := bot.Faq("xiao", "机器人怎么购买")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("TestFaq reply")
	}
}
func TestUsers(t *testing.T) {
	_, _, _, _, err := bot.Users(1, 1, "-lasttime")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Pass")
	}
}
func TestChats(t *testing.T) {
	_, _, _, _, err := bot.Chats("xiao", 1, 1, "-lasttime")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Pass")
	}
}
func TestMute(t *testing.T) {
	err := bot.Mute("xiao")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Pass")
	}
}
func TestUnmute(t *testing.T) {
	err := bot.Unmute("xiao")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Pass")
	}
}
func TestIsmute(t *testing.T) {
	_, err := bot.Ismute("xiao")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Pass")
	}
}
func TestUser(t *testing.T) {
	_, err := bot.User("xiao")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Pass")
	}
}
