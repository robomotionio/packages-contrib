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
	CreateHCTaskRequest struct {
		CreateRCTaskRequest
	}

	HCTask struct {
		RCTask
	}

	HCTaskResultResponse struct {
		RCTaskResultResponse
	}

	HCSolution struct {
		RCSolution
	}
)

type HCaptcha struct {
	runtime.Node `spec:"id=Robomotion.AntiCaptcha.HCaptcha,name=HCaptcha,icon=mdiRobotIndustrial,color=#FBAD00"`

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
func (n *HCaptcha) OnCreate() error {
	return nil
}

// OnMessage runs everytime a message is received
func (n *HCaptcha) OnMessage(ctx message.Context) (err error) {
	var (
		result, status string
		taskId         int
	)

	creds, err := n.OptToken.Get(ctx)
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.HCaptcha", err.Error())
		return err
	}

	inToken := creds["value"].(string)
	if inToken == "" {
		err = runtime.NewError("Robomotion.AntiCaptcha.HCaptcha", "Token can not be empty")
		return err
	}

	timeOut, err := n.OptTimeout.GetString(ctx)
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.HCaptcha.EmptyField", err.Error())
		return err
	}
	//TODO: Designer dan number gelmiyor. Type ı number yapınca transporting closed hatası alıyorum.Bu yüzden designer dan string alıp int e çevirdim
	OptTimeout, err := strconv.Atoi(timeOut)
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.HCaptcha.EmptyField", err.Error())
		return err
	}
	if OptTimeout < 0 {
		err = runtime.NewError("Robomotion.AntiCaptcha.HCaptcha.WrongField", "Timout can not be less than zero")
		return err
	}

	inWebsiteUrl, err := n.InWebsiteURL.GetString(ctx)
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.HCaptcha.EmptyField", err.Error())
		return err
	}
	if inWebsiteUrl == "" {
		err = runtime.NewError("Robomotion.AntiCaptcha.HCaptcha.EmptyField", "Website URL can not be empty")
		return err
	}

	inWebsiteKey, err := n.InWebsiteKey.GetString(ctx)
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.HCaptcha.EmptyField", err.Error())
		return err
	}
	if inWebsiteKey == "" {
		err = runtime.NewError("Robomotion.AntiCaptcha.HCaptcha.EmptyField", "Website Key can not be empty")
		return err
	}
	taskId, err = createHCaptchaTask(inToken, inWebsiteUrl, inWebsiteKey) //Task is created
	if err != nil {
		err = runtime.NewError("Robomotion.AntiCaptcha.HCaptcha.EncodeBase64", err.Error())
		return err
	}
	for {
		if OptTimeout == 0 {
			err = runtime.NewError("Robomotion.AntiCaptcha.HCaptcha.TimedOut", "The captcha could not be solved in given timeout")
			return err
		}
		status, result, err = ControlHCaptchaTask(inToken, taskId)
		if err != nil {
			err = runtime.NewError("Robomotion.AntiCaptcha.HCaptcha.EncodeBase64", err.Error())
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
func (n *HCaptcha) OnClose() error {
	return nil
}

func createHCaptchaTask(token, websiteUrl, websiteKey string) (int, error) {
	const url = "https://api.anti-captcha.com/createTask"
	const method = "POST"

	task := RCTask{
		Type:       "HCaptchaTaskProxyless",
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

func ControlHCaptchaTask(token string, taskId int) (string, string, error) {
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
	var body HCTaskResultResponse
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
	gHcaptchaResponse := body.Solution.GRecaptchaResponse

	return status, gHcaptchaResponse, nil
}
