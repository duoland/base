package wechat

// See doc at https://work.weixin.qq.com/api/doc/90000/90136/91770

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

// WxWorkRobotMessageAPI is the api to send the robot messages
const WxWorkRobotMessageAPI = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send"

// WxWorkRobotUploadFileAPI is the api to upload file
const WxWorkRobotUploadFileAPI = "https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media"

// WxWorkRobotTimeout is the wxwork robot default timeout
const WxWorkRobotTimeout = time.Second * 30
const WxWorkRobotStatusOK = 0

const (
	WxMessageTypeText     = "text"
	WxMessageTypeMarkdown = "markdown"
	WxMessageTypeImage    = "image"
	WxMessageTypeNews     = "news"
	WxMessageTypeFile     = "file"
)

type WxTextMessage struct {
	MessageType string            `json:"msgtype"`
	MessageBody WxTextMessageBody `json:"text"`
}

type WxTextMessageBody struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list,omitempty"`
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
}

type WxMarkdownMessageBody WxTextMessageBody

type WxMarkdownMessage struct {
	MessageType string                `json:"msgtype"`
	MessageBody WxMarkdownMessageBody `json:"markdown"`
}

type WxImageMessage struct {
	MessageType string             `json:"msgtype"`
	MessageBody WxImageMessageBody `json:"image"`
}
type WxImageMessageBody struct {
	Base64 string `json:"base64"`
	MD5    string `json:"md5"`
}

type WxNewsMessage struct {
	MessageType string            `json:"msgtype"`
	MessageBody WxNewsMessageBody `json:"news"`
}
type WxNewsMessageBody struct {
	Articles []WxNewsMessageArticle `json:"articles"`
}

type WxNewsMessageArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PictureURL  string `json:"picurl"`
}

type WxFileMessage struct {
	MessageType string            `json:"msgtype"`
	MessageBody WxFileMessageBody `json:"file"`
}

type WxFileMessageBody struct {
	MediaID string `json:"media_id"`
}

// WxWorkRobot is a robot to send wxwork messages
type WxWorkRobot struct {
	key    string
	client *http.Client
}

type WxMessageResp struct {
	ErrCode    int    `json:"errcode"`
	ErrMessage string `json:"errmsg"`
}

type WxUploadFileResp struct {
	ErrCode    int    `json:"errcode"`
	ErrMessage string `json:"errmsg"`
	Type       string `json:"type"`
	MediaID    string `json:"media_id"`
	CreatedAt  string `json:"created_at"`
}

// NewWxWorkRobot create a new wxwork robot
func NewWxWorkRobot(key string) *WxWorkRobot {
	return NewWxWorkRobotWithTimeout(key, WxWorkRobotTimeout)
}

// NewWxWorkRobotWithTimeout create a new wxwork robot with timeout
func NewWxWorkRobotWithTimeout(key string, timeout time.Duration) *WxWorkRobot {
	client := http.Client{}
	client.Timeout = timeout
	return &WxWorkRobot{key: key, client: &client}
}

// SendTextMessage send the text message
func (r *WxWorkRobot) SendTextMessage(text string) (err error) {
	textMessage := WxTextMessage{
		MessageType: WxMessageTypeText,
		MessageBody: WxTextMessageBody{
			Content: text,
		},
	}
	return r.sendMessage(&textMessage)
}

// SendTextMessage send the text message with specified mentioned users
func (r *WxWorkRobot) SendTextMessageWithMention(content string, mentionedList []string, mentionedMobileList []string) (err error) {
	textMessage := WxTextMessage{
		MessageType: WxMessageTypeText,
		MessageBody: WxTextMessageBody{
			Content:             content,
			MentionedList:       mentionedList,
			MentionedMobileList: mentionedMobileList,
		},
	}
	return r.sendMessage(&textMessage)
}

// SendMarkdownMessage send the markdown message
func (r *WxWorkRobot) SendMarkdownMessage(content string) (err error) {
	markdownMessage := WxMarkdownMessage{
		MessageType: WxMessageTypeMarkdown,
		MessageBody: WxMarkdownMessageBody{
			Content: content,
		},
	}
	return r.sendMessage(&markdownMessage)
}

// SendMarkdownMessageWithMention send the markdown message with specified mentioned users
func (r *WxWorkRobot) SendMarkdownMessageWithMention(content string, mentionedList []string, mentionedMobileList []string) (err error) {
	markdownMessage := WxMarkdownMessage{
		MessageType: WxMessageTypeMarkdown,
		MessageBody: WxMarkdownMessageBody{
			Content:             content,
			MentionedList:       mentionedList,
			MentionedMobileList: mentionedMobileList,
		},
	}
	return r.sendMessage(&markdownMessage)
}

// SendImageMessage send the markdown message
func (r *WxWorkRobot) SendImageMessage(imageData []byte) (err error) {
	imageHash := md5.Sum(imageData)
	imageMessage := WxImageMessage{
		MessageType: WxMessageTypeImage,
		MessageBody: WxImageMessageBody{
			Base64: base64.StdEncoding.EncodeToString(imageData),
			MD5:    fmt.Sprintf("%x", imageHash),
		},
	}
	return r.sendMessage(&imageMessage)
}

// SendNewsMessage send the news message
func (r *WxWorkRobot) SendNewsMessage(articles []WxNewsMessageArticle) (err error) {
	newsMessage := WxNewsMessage{
		MessageType: WxMessageTypeNews,
		MessageBody: WxNewsMessageBody{
			Articles: articles,
		},
	}
	return r.sendMessage(&newsMessage)
}

// SendFileMessage send the file message
func (r *WxWorkRobot) SendFileMessage(mediaID string) (err error) {
	fileMessage := WxFileMessage{
		MessageType: WxMessageTypeFile,
		MessageBody: WxFileMessageBody{
			MediaID: mediaID,
		},
	}
	return r.sendMessage(&fileMessage)
}

func (r *WxWorkRobot) sendMessage(messageObj interface{}) (err error) {
	reqURL := fmt.Sprintf("%s?key=%s", WxWorkRobotMessageAPI, r.key)
	reqBody, _ := json.Marshal(messageObj)

	req, newErr := http.NewRequest(http.MethodPost, reqURL, bytes.NewReader(reqBody))
	if newErr != nil {
		err = fmt.Errorf("create request error, %s", newErr.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, getErr := r.client.Do(req)
	if getErr != nil {
		err = fmt.Errorf("get response error, %s", getErr.Error())
		return
	}
	defer resp.Body.Close()
	// check http code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("wxwork request error, %s", resp.Status)
		io.Copy(ioutil.Discard, resp.Body)
		return
	}
	// parse response body
	decoder := json.NewDecoder(resp.Body)
	var wxMessageResp WxMessageResp
	if decodeErr := decoder.Decode(&wxMessageResp); decodeErr != nil {
		err = fmt.Errorf("parse response error, %s", decodeErr.Error())
		return
	}
	if wxMessageResp.ErrCode != WxWorkRobotStatusOK {
		err = fmt.Errorf("call wxwork api error, %d %s", wxMessageResp.ErrCode, wxMessageResp.ErrMessage)
		return
	}
	return
}

// UploadFile upload the media file
func (r *WxWorkRobot) UploadFile(fileBody []byte, fileName string) (mediaID string, createdAt int64, err error) {
	respBodyBuffer := bytes.NewBuffer(nil)
	multipartWriter := multipart.NewWriter(respBodyBuffer)
	// add form data
	formFileWriter, createErr := multipartWriter.CreateFormFile("media", fileName)
	if createErr != nil {
		err = fmt.Errorf("create form file error, %s", createErr.Error())
		return
	}
	if _, writeErr := formFileWriter.Write(fileBody); writeErr != nil {
		err = fmt.Errorf("write form file error, %s", writeErr.Error())
		return
	}
	if closeErr := multipartWriter.Close(); closeErr != nil {
		err = fmt.Errorf("close form file error, %s", closeErr.Error())
		return
	}

	reqURL := fmt.Sprintf("%s?key=%s&type=%s", WxWorkRobotUploadFileAPI, r.key, WxMessageTypeFile)
	req, newErr := http.NewRequest(http.MethodPost, reqURL, respBodyBuffer)
	if newErr != nil {
		err = fmt.Errorf("create request error, %s", newErr.Error())
		return
	}
	// set multi-part header
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	resp, getErr := r.client.Do(req)
	if getErr != nil {
		err = fmt.Errorf("get response error, %s", getErr.Error())
		return
	}
	defer resp.Body.Close()
	// check http code
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("wxwork request error, %s", resp.Status)
		io.Copy(ioutil.Discard, resp.Body)
		return
	}
	// parse response body
	decoder := json.NewDecoder(resp.Body)
	var wxUploadFileResp WxUploadFileResp
	if decodeErr := decoder.Decode(&wxUploadFileResp); decodeErr != nil {
		err = fmt.Errorf("parse response error, %s", decodeErr.Error())
		return
	}
	if wxUploadFileResp.ErrCode != WxWorkRobotStatusOK {
		err = fmt.Errorf("call wxwork api error, %d %s", wxUploadFileResp.ErrCode, wxUploadFileResp.ErrMessage)
		return
	}

	// set fields
	mediaID = wxUploadFileResp.MediaID
	createdAt, _ = strconv.ParseInt(wxUploadFileResp.CreatedAt, 10, 64)
	return
}
