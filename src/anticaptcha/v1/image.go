package anticaptcha

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
)

//the structs for anticapthca api
type CreateTaskRequest struct {
	ClientKey string `json:"clientKey"`
	Task      Task   `json:"task"`
}
type Task struct {
	Type      string `json:"type"`
	Body      string `json:"body"`
	Phrase    bool   `json:"phrase"`
	Case      bool   `json:"case"`
	Numeric   bool   `json:"numeric"`
	Math      int    `json:"math"`
	MinLength int    `json:"minLength"`
	MaxLength int    `json:"maxLength"`
}
type CreateTaskResponse struct {
	ErrorId int `json:"errorId"`
	TaskId  int `json:"taskId"`
}
type TaskResultRequest struct {
	ClientKey string `json:"clientKey"`
	TaskId    int    `json:"taskId"`
}
type TaskResultResponse struct {
	ErrorId          int      `json:"errorId"`
	ErrorDescription string   `json:"errorDescription"`
	Status           string   `json:"status"`
	Solution         Solution `json:"solution"`
}
type Solution struct {
	Text string `json:"text"`
	Url  string `json:url`
}

// Image holds this Node's properties
type Image struct {
	runtime.Node `spec:"id=Robomotion.AntiCaptcha.Image,name=Image Captcha,icon=mdiImage,color=#065c95"`

	//Inputs
	InTimeOut   runtime.InVariable `spec:"title=Timeout,type=string,scope=Custom,name=30,messageScope,customScope"`
	InImagePath runtime.InVariable `spec:"title=Image Path,type=string,scope=Custom,messageScope,customScope"`

	//Options
	OptToken runtime.Credential `spec:"title=Credentials,option"`

	//Outputs
	OutResult runtime.OutVariable `spec:"title=result,type=string,scope=Message,name=result,messageScope"`
}

// OnCreate runs once when a flow starts running
func (n *Image) OnCreate() error {
	return nil
}

// OnMessage runs everytime a message is received
func (n *Image) OnMessage(ctx message.Context) (err error) {
	var (
		result, status string
		taskId         int
	)

	creds, err := n.OptToken.Get()
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.Image", err.Error())
		return err
	}

	inToken := creds["value"].(string)
	if inToken == "" {
		err = runtime.NewError("Robomotion.AntiCaptcha.Image", "Token can not be empty")
		return err
	}

	timeOut, err := n.InTimeOut.GetString(ctx)
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.Image.EmptyField", err.Error())
		return err
	}
	//TODO: Designer dan number gelmiyor. Type ı number yapınca transporting closed hatası alıyorum.Bu yüzden designer dan string alıp int e çevirdim
	inTimeOut, err := strconv.Atoi(timeOut)
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.Image.EmptyField", err.Error())
		return err
	}
	if inTimeOut < 0 {
		err = runtime.NewError("Robomotion.AntiCaptcha.Image.WrongField", "Timout can not be less than zero")
		return err
	}

	inImagePath, err := n.InImagePath.GetString(ctx)
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.Image.EmptyField", err.Error())
		return err
	}
	if inImagePath == "" {
		err = runtime.NewError("Robomotion.AntiCaptcha.Image.EmptyField", "Image Path can not be empty")
		return err
	}

	base64, err := encodeBase64(inImagePath) //The image is encoded as base64
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.Image.EncodeBase64", err.Error())
		return err
	}

	taskId, err = createTask(inToken, base64) //Task is created
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.Image.EncodeBase64", err.Error())
		return err
	}

	for {
		if inTimeOut == 0 {
			err = runtime.NewError("Robomotion.AntiCaptcha.Image.TimedOut", "The captcha could not be solved in given timeout")
			return err
		}
		status, result, err = ControlTask(inToken, taskId)
		if err != nil {
			err = runtime.NewError("Robomotion.AntiCaptcha.Image.EncodeBase64", err.Error())
			return err
		}
		if status == "ready" {
			break
		}

		time.Sleep(time.Second)
		inTimeOut--
	}

	err = n.OutResult.Set(ctx, result)
	return nil
}

// OnClose runs once when a flow stops running
func (n *Image) OnClose() error {
	return nil
}

func encodeBase64(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	// Read entire file into byte slice.
	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded, nil
}

func createTask(token, base64 string) (int, error) {
	const url = "https://api.anti-captcha.com/createTask"
	const method = "POST"

	task := Task{
		Type:      "ImageToTextTask",
		Body:      base64,
		Phrase:    false,
		Case:      false,
		Numeric:   false,
		Math:      0,
		MinLength: 0,
		MaxLength: 0,
	}
	ctReq := &CreateTaskRequest{
		ClientKey: token,
		Task:      task,
	}

	reader, err := json.Marshal(&ctReq)
	if err != nil {
		return -1, err
	}

	payload := strings.NewReader(string(reader))

	t := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		// We use ABSURDLY large keys, and should probably not.
		TLSHandshakeTimeout: 60 * time.Second,
	}

	client := &http.Client{
		Transport: t,
	}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return -1, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()
	var taskResp CreateTaskResponse
	err = json.Unmarshal(result, &taskResp)
	if err != nil {
		return -1, nil
	}

	return taskResp.TaskId, nil
}

func ControlTask(token string, taskId int) (string, string, error) {
	const url = "https://api.anti-captcha.com/getTaskResult"
	const method = "POST"

	ctReq := TaskResultRequest{
		ClientKey: token,
		TaskId:    taskId,
	}
	reader, err := json.Marshal(&ctReq)
	if err != nil {
		return "", "", err
	}
	payload := strings.NewReader(string(reader))

	t := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		// We use ABSURDLY large keys, and should probably not.
		TLSHandshakeTimeout: 60 * time.Second,
	}

	client := &http.Client{
		Transport: t,
	}

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", "", err
	}

	temp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	res, err := ioutil.ReadAll(temp.Body)
	defer temp.Body.Close()
	var body TaskResultResponse
	err = json.Unmarshal(res, &body)
	if err != nil {
		return "", "", err
	}
	errId := body.ErrorId
	if errId != 0 {
		err = runtime.NewError("Robomotion.AntiCaptcha.Image", body.ErrorDescription)
		return "", "", err
	}

	status := body.Status
	text := body.Solution.Text

	return status, text, nil
}
