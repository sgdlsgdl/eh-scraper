package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

var (
	client *http.Client
	once   sync.Once
)

func batchSend(proxy, token, chatId string, items ItemList) {
	if token == "" || chatId == "" {
		return
	}

	once.Do(func() {
		client = &http.Client{}
		if proxy != "" {
			proxyURL, _ := url.Parse(proxy)
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
			client.Transport = transport
		}
	})

	for _, item := range items {
		sendSingle(token, chatId, item)
	}
}

func sendSingle(token string, chatId string, item Item) {
	sendMessageUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	err := post(sendMessageUrl, map[string]string{
		"chat_id": chatId,
		"text":    item.Name,
	}, "", "")
	if err != nil {
		log.Printf("sendMessage failed %v", err)
	}

	if strings.Contains(item.Image, "http") {
		err = post(sendMessageUrl, map[string]string{
			"chat_id": chatId,
			"text":    item.Image,
		}, "", "")
		if err != nil {
			log.Printf("sendMessage failed %v", err)
		}
		return
	}

	sendPhotoUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", token)
	err = post(sendPhotoUrl, map[string]string{
		"chat_id": chatId,
	}, "photo", basePath+item.Image)
	if err != nil {
		log.Printf("sendPhoto failed %v", err)
	}
}

func post(url string, requestData map[string]string, fileKey, filePath string) error {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	for k, v := range requestData {
		_ = writer.WriteField(k, v)
	}
	if fileKey != "" && filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("os.Open failed: %w", err)
		}
		defer file.Close()
		ww, _ := writer.CreateFormFile(fileKey, file.Name())
		_, _ = io.Copy(ww, file)
	}
	_ = writer.Close()

	req, _ := http.NewRequest("POST", url, buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("client.Do failed: %w", err)
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll failed: %w", err)
	}

	var respData Response
	if err = json.Unmarshal(respBytes, &respData); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	if !respData.Ok {
		return fmt.Errorf("!respData.Ok")
	}
	return nil
}

type Response struct {
	Ok bool `json:"ok"`
}
