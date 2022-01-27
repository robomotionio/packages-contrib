package anticaptcha

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/robomotionio/robomotion-go/message"
	"github.com/robomotionio/robomotion-go/runtime"
)

type (
	CreateRCTaskRequest struct {
		ClientKey string `json:"clientKey"`
		Task      RCTask `json:"task"`
	}

	RCTask struct {
		Type       string `json:"type"`
		WebsiteURL string `json:"websiteURL"`
		WebsiteKey string `json:"websiteKey"`
	}

	RCTaskResultResponse struct {
		ErrorId          int        `json:"errorId"`
		ErrorDescription string     `json:"errorDescription"`
		Status           string     `json:"status"`
		Solution         RCSolution `json:"solution"`
	}

	RCSolution struct {
		Cookies            interface{} `json:"cookies"`
		GRecaptchaResponse string      `json:gRecaptchaResponse`
	}
)

//Create ReCaptcha Task Request
//ReCaptcha Task

// Image holds this Node's properties
type ReCaptcha struct {
	runtime.Node `spec:"id=Robomotion.AntiCaptcha.ReCaptcha,name=ReCaptcha,icon=mdiRobot,color=#FBAD00"`

	//Inputs
	InWebsiteKey runtime.InVariable `spec:"title=Website Key,type=string,scope=Custom,messageScope,customScope"`
	InWebsiteURL runtime.InVariable `spec:"title=Website URL,type=string,scope=Custom,messageScope,customScope"`

	//Outputs
	OutResult runtime.OutVariable `spec:"title=result,type=string,scope=Message,name=result,messageScope"`

	//Options
	OptToken   runtime.Credential  `spec:"title=Credentials,scope=Custom,option,messageScope,customScope"`
	OptTimeout runtime.OptVariable `spec:"title=Timeout,type=string,scope=Custom,name=180,messageScope,customScope"`
}

// OnCreate runs once when a flow starts running
func (n *ReCaptcha) OnCreate() error {
	return nil
}

// OnMessage runs everytime a message is received
func (n *ReCaptcha) OnMessage(ctx message.Context) (err error) {
	var (
		result, status string
		taskId         int
	)

	creds, err := n.OptToken.Get(ctx)
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.ReCaptcha", err.Error())
		return err
	}

	inToken := creds["value"].(string)
	if inToken == "" {
		err = runtime.NewError("Robomotion.AntiCaptcha.ReCaptcha", "Token can not be empty")
		return err
	}

	timeOut, err := n.OptTimeout.GetString(ctx)
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.ReCaptcha.EmptyField", err.Error())
		return err
	}
	//TODO: Designer dan number gelmiyor. Type ı number yapınca transporting closed hatası alıyorum.Bu yüzden designer dan string alıp int e çevirdim
	OptTimeout, err := strconv.Atoi(timeOut)
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.ReCaptcha.EmptyField", err.Error())
		return err
	}
	if OptTimeout < 0 {
		err = runtime.NewError("Robomotion.AntiCaptcha.ReCaptcha.WrongField", "Timout can not be less than zero")
		return err
	}

	inWebsiteUrl, err := n.InWebsiteURL.GetString(ctx)
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.ReCaptcha.EmptyField", err.Error())
		return err
	}
	if inWebsiteUrl == "" {
		err = runtime.NewError("Robomotion.AntiCaptcha.ReCaptcha.EmptyField", "Website URL can not be empty")
		return err
	}

	inWebsiteKey, err := n.InWebsiteKey.GetString(ctx)
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.ReCaptcha.EmptyField", err.Error())
		return err
	}
	if inWebsiteKey == "" {
		err = runtime.NewError("Robomotion.AntiCaptcha.ReCaptcha.EmptyField", "Website Key can not be empty")
		return err
	}
	taskId, err = createReCaptchaTask(inToken, inWebsiteUrl, inWebsiteKey) //Task is created
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.ReCaptcha.EncodeBase64", err.Error())
		return err
	}
	for {
		if OptTimeout == 0 {
			err = runtime.NewError("Robomotion.AntiCaptcha.ReCaptcha.TimedOut", "The captcha could not be solved in given timeout")
			return err
		}
		status, result, err = ControlReCaptchaTask(inToken, taskId)
		if err != nil {
			err = runtime.NewError("Robomotion.AntiCaptcha.ReCaptcha.EncodeBase64", err.Error())
			return err
		}
		if status == "ready" {
			break
		}

		time.Sleep(time.Second)
		OptTimeout--
	}
	err = n.OutResult.Set(ctx, result)
	return nil
}

// OnClose runs once when a flow stops running
func (n *ReCaptcha) OnClose() error {
	return nil
}

func createReCaptchaTask(token, websiteUrl, websiteKey string) (int, error) {
	const url = "https://api.anti-captcha.com/createTask"
	const method = "POST"

	task := RCTask{
		Type:       "NoCaptchaTaskProxyless",
		WebsiteURL: websiteUrl,
		WebsiteKey: websiteKey,
	}
	ctReq := &CreateRCTaskRequest{
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
		TLSHandshakeTimeout: 60 * time.Second,
	}

	t.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS10,
		MaxVersion: tls.VersionTLS12,
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

func ControlReCaptchaTask(token string, taskId int) (string, string, error) {
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
		TLSHandshakeTimeout: 60 * time.Second,
	}

	t.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS10,
		MaxVersion: tls.VersionTLS12,
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
	var body RCTaskResultResponse
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
	gRecaptchaResponse := body.Solution.GRecaptchaResponse

	return status, gRecaptchaResponse, nil
}
