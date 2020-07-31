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
const WxWorkRobotTimeout = time.Second * 10
const WxWorkRobotStatusOK = 0

const (
	WxWorkRobotMessageTypeText     = "text"
	WxWorkRobotMessageTypeMarkdown = "markdown"
	WxWorkRobotMessageTypeImage    = "image"
	WxWorkRobotMessageTypeNews     = "news"
	WxWorkRobotMessageTypeFile     = "file"
)

type WxWorkRobotTextMessage struct {
	MessageType string                     `json:"msgtype"`
	MessageBody WxWorkRobotTextMessageBody `json:"text"`
}

type WxWorkRobotTextMessageBody struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list,omitempty"`
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
}

type WxWorkRobotMarkdownMessageBody WxWorkRobotTextMessageBody

type WxWorkRobotMarkdownMessage struct {
	MessageType string                         `json:"msgtype"`
	MessageBody WxWorkRobotMarkdownMessageBody `json:"markdown"`
}

type WxWorkRobotImagMessage struct {
	MessageType string                     `json:"msgtype"`
	MessageBody WxWorkRobotImagMessageBody `json:"image"`
}
type WxWorkRobotImagMessageBody struct {
	Base64 string `json:"base64"`
	MD5    string `json:"md5"`
}

type WxWorkRobotNewsMessage struct {
	MessageType string                     `json:"msgtype"`
	MessageBody WxWorkRobotNewsMessageBody `json:"news"`
}
type WxWorkRobotNewsMessageBody struct {
	Articles []WxWorkRobotNewsMessageArticle `json:"articles"`
}

type WxWorkRobotNewsMessageArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PictureURL  string `json:"picurl"`
}

type WxWorkRobotFileMessage struct {
	MessageType string                     `json:"msgtype"`
	MessageBody WxWorkRobotFileMessageBody `json:"file"`
}

type WxWorkRobotFileMessageBody struct {
	MediaID string `json:"media_id"`
}

// WxWorkRobot is a robot to send wxwork messages
type WxWorkRobot struct {
	client *http.Client
}

type WxWorkRobotMessageResp struct {
	ErrCode    int    `json:"errcode"`
	ErrMessage string `json:"errmsg"`
}

type WxWorkRobotUploadFileResp struct {
	ErrCode    int    `json:"errcode"`
	ErrMessage string `json:"errmsg"`
	Type       string `json:"type"`
	MediaID    string `json:"media_id"`
	CreatedAt  string `json:"created_at"`
}

// NewWxWorkRobot create a new wxwork robot
func NewWxWorkRobot() *WxWorkRobot {
	return NewWxWorkRobotWithTimeout(WxWorkRobotTimeout)
}

// NewWxWorkRobotWithTimeout create a new wxwork robot with timeout
func NewWxWorkRobotWithTimeout(timeout time.Duration) *WxWorkRobot {
	client := http.Client{}
	client.Timeout = timeout
	return &WxWorkRobot{client: &client}
}

// NewWxWorkRobotWithClient create a new wxwork robot with http.Client
func NewWxWorkRobotWithClient(client *http.Client) *WxWorkRobot {
	return &WxWorkRobot{client: client}
}

// SendTextMessage send the text message
func (r *WxWorkRobot) SendTextMessage(key, text string) (err error) {
	textMessage := WxWorkRobotTextMessage{
		MessageType: WxWorkRobotMessageTypeText,
		MessageBody: WxWorkRobotTextMessageBody{
			Content: text,
		},
	}
	return r.sendMessage(key, &textMessage)
}

// SendTextMessage send the text message with specified mentioned users
func (r *WxWorkRobot) SendTextMessageWithMention(key, content string, mentionedList []string, mentionedMobileList []string) (err error) {
	textMessage := WxWorkRobotTextMessage{
		MessageType: WxWorkRobotMessageTypeText,
		MessageBody: WxWorkRobotTextMessageBody{
			Content:             content,
			MentionedList:       mentionedList,
			MentionedMobileList: mentionedMobileList,
		},
	}
	return r.sendMessage(key, &textMessage)
}

// SendMarkdownMessage send the markdown message
func (r *WxWorkRobot) SendMarkdownMessage(key, content string) (err error) {
	markdownMessage := WxWorkRobotMarkdownMessage{
		MessageType: WxWorkRobotMessageTypeMarkdown,
		MessageBody: WxWorkRobotMarkdownMessageBody{
			Content: content,
		},
	}
	return r.sendMessage(key, &markdownMessage)
}

// SendMarkdownMessageWithMention send the markdown message with specified mentioned users
func (r *WxWorkRobot) SendMarkdownMessageWithMention(key, content string, mentionedList []string, mentionedMobileList []string) (err error) {
	markdownMessage := WxWorkRobotMarkdownMessage{
		MessageType: WxWorkRobotMessageTypeMarkdown,
		MessageBody: WxWorkRobotMarkdownMessageBody{
			Content:             content,
			MentionedList:       mentionedList,
			MentionedMobileList: mentionedMobileList,
		},
	}
	return r.sendMessage(key, &markdownMessage)
}

// SendImageMessage send the markdown message
func (r *WxWorkRobot) SendImageMessage(key string, imageData []byte) (err error) {
	imageHash := md5.Sum(imageData)
	imageMessage := WxWorkRobotImagMessage{
		MessageType: WxWorkRobotMessageTypeImage,
		MessageBody: WxWorkRobotImagMessageBody{
			Base64: base64.StdEncoding.EncodeToString(imageData),
			MD5:    fmt.Sprintf("%x", imageHash),
		},
	}
	return r.sendMessage(key, &imageMessage)
}

// SendNewsMessage send the news message
func (r *WxWorkRobot) SendNewsMessage(key string, articles []WxWorkRobotNewsMessageArticle) (err error) {
	newsMessage := WxWorkRobotNewsMessage{
		MessageType: WxWorkRobotMessageTypeNews,
		MessageBody: WxWorkRobotNewsMessageBody{
			Articles: articles,
		},
	}
	return r.sendMessage(key, &newsMessage)
}

// SendFileMessage send the file message
func (r *WxWorkRobot) SendFileMessage(key, mediaID string) (err error) {
	fileMessage := WxWorkRobotFileMessage{
		MessageType: WxWorkRobotMessageTypeFile,
		MessageBody: WxWorkRobotFileMessageBody{
			MediaID: mediaID,
		},
	}
	return r.sendMessage(key, &fileMessage)
}

func (r *WxWorkRobot) sendMessage(key string, messageObj interface{}) (err error) {
	reqURL := fmt.Sprintf("%s?key=%s", WxWorkRobotMessageAPI, key)
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
	var wxMessageResp WxWorkRobotMessageResp
	if decodeErr := decoder.Decode(&wxMessageResp); decodeErr != nil {
		err = fmt.Errorf("parse response error, %s", decodeErr.Error())
		return
	}
	if wxMessageResp.ErrCode != WxWorkRobotStatusOK {
		err = fmt.Errorf("call wxwork robot api error, %d %s", wxMessageResp.ErrCode, wxMessageResp.ErrMessage)
		return
	}
	return
}

// UploadFile upload the media file
func (r *WxWorkRobot) UploadFile(key string, fileBody []byte, fileName string) (mediaID string, createdAt int64, err error) {
	respBodyBuffer := bytes.NewBuffer(nil)
	defer respBodyBuffer.Reset()
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

	reqURL := fmt.Sprintf("%s?key=%s&type=%s", WxWorkRobotUploadFileAPI, key, WxWorkRobotMessageTypeFile)
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
		err = fmt.Errorf("wxwork robot request error, %s", resp.Status)
		io.Copy(ioutil.Discard, resp.Body)
		return
	}
	// parse response body
	decoder := json.NewDecoder(resp.Body)
	var wxUploadFileResp WxWorkRobotUploadFileResp
	if decodeErr := decoder.Decode(&wxUploadFileResp); decodeErr != nil {
		err = fmt.Errorf("parse response error, %s", decodeErr.Error())
		return
	}
	if wxUploadFileResp.ErrCode != WxWorkRobotStatusOK {
		err = fmt.Errorf("call wxwork robot api error, %d %s", wxUploadFileResp.ErrCode, wxUploadFileResp.ErrMessage)
		return
	}

	// set fields
	mediaID = wxUploadFileResp.MediaID
	createdAt, _ = strconv.ParseInt(wxUploadFileResp.CreatedAt, 10, 64)
	return
}
