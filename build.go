package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/go-version"
)

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func main() {
	//reading index.json file
	data, err := ioutil.ReadFile("./index.json")
	if err != nil {
		fmt.Print(err)
	}
	var intr interface{}

	err = json.Unmarshal(data, &intr)
	if err != nil {
		fmt.Println("error:", err)
	}
	str := fmt.Sprintf("%v", intr)
	//index mapped
	indexmap := ReadIndex(str)

	var repoDirectory string //for mv repo command

	successError := make(map[string]string)

	// tree walk
	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			//checking path
			configCheck := strings.Contains(path, "config.json")
			distCheck := strings.Contains(path, "dist")
			if configCheck == true && distCheck == false {

				dataConfig, err := ioutil.ReadFile(path)
				if err != nil {
					fmt.Print(err)
				}

				var configIntr interface{}
				err = json.Unmarshal(dataConfig, &configIntr)
				if err != nil {
					fmt.Println("error:", err)
				}
				stringConfig := fmt.Sprintf("%v", configIntr)

				configMap, configLanguage, configPlatform, configIcon := ReadConfig(stringConfig)
				//path without config.json
				newPath := path[:len(path)-12]
				newPath = filepath.ToSlash(newPath)
				var out bytes.Buffer
				for keyConfig, valConfig := range configMap {
					folderName := strings.Split(keyConfig, ".")
					folderNameRepo := folderName[1]
					pythonLanguage, javaLanguage, dotnetLanguage := false, false, false
					if configLanguage[keyConfig] == "Python" {
						pythonLanguage = true
					} else if configLanguage[keyConfig] == "Java" {
						javaLanguage = true
						if runtime.GOOS == "linux" && runtime.GOARCH == "arm" {
							sedCommand := exec.Command("sh", "-c", " sed -i 's+x64_linux+arm_linux+' "+newPath+"/config.json ")
							sedCommand.Stdout = &out
							err = sedCommand.Run()
							if err != nil {
								fmt.Println("error when changing build script", err)
							}
						}
					} else if configLanguage[keyConfig] == "C#" {
						dotnetLanguage = true
					}
					if runtime.GOARCH != "arm" || !dotnetLanguage {
						runneros := os.Getenv("RUNNEROS")
						var commandStr *exec.Cmd
						if configIcon[keyConfig] == "icon.png" {
							if runneros == "Linux" {
								commandStr = exec.Command("sh", "-c", " cd "+newPath+" && rclone sync -P --s3-acl=public-read ./icon.png DO:robomotion-packages/contrib/"+folderNameRepo+"/ && sed -i 's+icon.png+https://packages.robomotion.io/contrib/"+folderNameRepo+"/icon.png+' config.json ")
							} else if runneros == "Windows" {
								commandStr = exec.Command("sh", "-c", " cd "+newPath+" && D:/a/packages-contrib/packages-contrib/rclone.exe sync -P --s3-acl=public-read ./icon.png DO:robomotion-packages/contrib/"+folderNameRepo+"/ && sed -i 's+icon.png+https://packages.robomotion.io/contrib/"+folderNameRepo+"/icon.png+' config.json ")
							} else if runneros == "macOS" {
								commandStr = exec.Command("sh", "-c", " cd "+newPath+" && rclone sync -P --s3-acl=public-read ./icon.png DO:robomotion-packages/contrib/"+folderNameRepo+"/ && perl -i -pe's+icon.png+https://packages.robomotion.io/contrib/"+folderNameRepo+"/icon.png+' config.json ")
							}
							commandStr.Stdout = &out
							err = commandStr.Run()
							if err != nil {
								fmt.Println("rclone commands error", err)
							}

						}
						platform := configPlatform[keyConfig]
						var platformCheck bool
						roboctl := "roboctl"

						if runneros == "Linux" {
							runneros = "linux"
							platformCheck = strings.Contains(platform, runneros)
							repoDirectory = "/home/runner/work/packages-contrib/packages-contrib/repo"
							if runtime.GOARCH == "arm" {
								runneros = "rpi"
								repoDirectory = "/home/pi/external/4c2f8c8e-6a71-4fb2-904a-2f76ffd42d8d/packages-runner/_work/packages-contrib/packages-contrib/repo"
							}
						} else if runneros == "Windows" {
							runneros = "windows"
							platformCheck = strings.Contains(platform, runneros)
							roboctl = "D:/a/packages-contrib/packages-contrib/roboctl"
							repoDirectory = "D:/a/packages-contrib/packages-contrib/repo"
						} else if runneros == "macOS" {
							repoDirectory = "/Users/runner/work/packages-contrib/packages-contrib/repo"
							runneros = "darwin"
							platformCheck = strings.Contains(platform, runneros)
						}

						//change config

						if strings.Contains(platform, "any") == true {
							platformCheck = true
						}

						exist := true
						_, err = os.Stat("repo/" + folderNameRepo + "")
						if os.IsNotExist(err) {
							exist = false
						}
						changelogCommand := exec.Command("sh", "-c", " cd "+newPath+" && rclone sync -P --s3-acl=public-read ./CHANGELOG.md DO:robomotion-packages/contrib/"+folderNameRepo+"/ ")
						if platformCheck {
							if value, ok := indexmap[keyConfig]; ok {
								//version compare
								v1, err := version.NewVersion(valConfig)
								if err != nil {
									fmt.Println(err)
								}
								v2, err := version.NewVersion(value)
								if err != nil {
									fmt.Println(err)
								}

								var comman *exec.Cmd
								if v1.GreaterThan(v2) {
									//config version > index
									if exist == false {
										comman = exec.Command("sh", "-c", " cd "+newPath+" && "+roboctl+" package -b "+runneros+" && mkdir "+folderNameRepo+" && mv robomotion-* "+folderNameRepo+"  && mv "+folderNameRepo+" "+repoDirectory+"/"+folderNameRepo+"")
										if pythonLanguage {
											comman = exec.Command("sh", "-c", " cd "+newPath+" && pipreqs . && pip install -r requirements.txt && "+roboctl+" package -b "+runneros+" && mkdir "+folderNameRepo+" && mv robomotion-* "+folderNameRepo+"  && mv "+folderNameRepo+" "+repoDirectory+"/"+folderNameRepo+"")
										} else if javaLanguage {
											comman = exec.Command("sh", "-c", " cd "+newPath+" && mv ../../OpenJDK11U-* . && "+roboctl+" package -b "+runneros+" && mkdir "+folderNameRepo+" && mv robomotion-* "+folderNameRepo+"  && mv "+folderNameRepo+" "+repoDirectory+"/"+folderNameRepo+"")
										}

									} else {
										comman = exec.Command("sh", "-c", " cd "+newPath+" && "+roboctl+" package -b "+runneros+" && mkdir "+folderNameRepo+" && mv robomotion-* "+folderNameRepo+"  && mv "+folderNameRepo+"/* "+repoDirectory+"/"+folderNameRepo+"")
										if pythonLanguage {
											comman = exec.Command("sh", "-c", " cd "+newPath+" && pipreqs . && pip install -r requirements.txt && "+roboctl+" package -b "+runneros+" && mkdir "+folderNameRepo+" && mv robomotion-* "+folderNameRepo+"  && mv "+folderNameRepo+"/* "+repoDirectory+"/"+folderNameRepo+"")
										} else if javaLanguage {
											comman = exec.Command("sh", "-c", " cd "+newPath+" && mv ../../OpenJDK11U-* . && "+roboctl+" package -b "+runneros+" && mkdir "+folderNameRepo+" && mv robomotion-* "+folderNameRepo+"  && mv "+folderNameRepo+"/* "+repoDirectory+"/"+folderNameRepo+"")
										}
									}

									comman.Stdout = &out
									err = comman.Run()
									if err != nil {
										successError[keyConfig] = err.Error()
										log.Fatal(keyConfig+" ", err.Error())
									} else {
										successError[keyConfig] = out.String()
									}
									changelogCommand.Stdout = &out
									err = changelogCommand.Run()
									if err != nil {
										fmt.Println("changelog command error", err)
									}
								}
							} else {

								var cmd *exec.Cmd
								if exist == false {
									cmd = exec.Command("sh", "-c", " cd "+newPath+" && "+roboctl+" package -b "+runneros+" && mkdir "+folderNameRepo+" && mv robomotion-* "+folderNameRepo+" && mv "+folderNameRepo+" "+repoDirectory+"/"+folderNameRepo+"")
									if pythonLanguage {
										cmd = exec.Command("sh", "-c", " cd "+newPath+" && pipreqs . && pip install -r requirements.txt && "+roboctl+" package -b "+runneros+" && mkdir "+folderNameRepo+" && mv robomotion-* "+folderNameRepo+" && mv "+folderNameRepo+" "+repoDirectory+"/"+folderNameRepo+"")
									} else if javaLanguage {
										cmd = exec.Command("sh", "-c", " cd "+newPath+" && mv ../../OpenJDK11U-* . && "+roboctl+" package -b "+runneros+" && mkdir "+folderNameRepo+" && mv robomotion-* "+folderNameRepo+" && mv "+folderNameRepo+" "+repoDirectory+"/"+folderNameRepo+"")
									}
								} else {

									cmd = exec.Command("sh", "-c", " cd "+newPath+" && "+roboctl+" package -b "+runneros+" && mkdir "+folderNameRepo+" && mv robomotion-* "+folderNameRepo+" && mv "+folderNameRepo+"/* "+repoDirectory+"/"+folderNameRepo+"")
									if pythonLanguage {
										cmd = exec.Command("sh", "-c", " cd "+newPath+" && pipreqs . && pip install -r requirements.txt && "+roboctl+" package -b "+runneros+" && mkdir "+folderNameRepo+" && mv robomotion-* "+folderNameRepo+" && mv "+folderNameRepo+"/* "+repoDirectory+"/"+folderNameRepo+"")
									} else if javaLanguage {
										cmd = exec.Command("sh", "-c", " cd "+newPath+" && mv ../../OpenJDK11U-* . && "+roboctl+" package -b "+runneros+" && mkdir "+folderNameRepo+" && mv robomotion-* "+folderNameRepo+" && mv "+folderNameRepo+"/* "+repoDirectory+"/"+folderNameRepo+"")
									}
								}

								cmd.Stdout = &out
								err = cmd.Run()
								if err != nil {
									successError[keyConfig] = err.Error()
									log.Fatal(keyConfig+" ", err.Error())
								} else {
									successError[keyConfig] = out.String()
								}

								changelogCommand.Stdout = &out
								err = changelogCommand.Run()
								if err != nil {
									fmt.Println("changelog command error", err)
								}

							}
						}
					}
				}
			}

			return nil
		})
	if err != nil {
		fmt.Println(err)
	}
	messageStr := strings.Builder{}
	subjectStr := strings.Builder{}
	subjectStr.WriteString("Main Repo's Deployment   ")
	messageStr.WriteString("RUNNER " + os.Getenv("RUNNEROS") + " " + runtime.GOARCH + "  ")
	for key, value := range successError {
		messageStr.WriteString(key + value)
		if strings.Contains(value, "Done.") {
			subjectStr.WriteString("" + key + "-> Successful ")
		} else {
			subjectStr.WriteString("" + key + "-> Failed ")
		}

	}

	password := os.Getenv("PASSWORD")
	from := os.Getenv("MAIL")

	to := []string{
		"rohat@robomotion.io",
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	request := Mail{
		Sender:  from,
		To:      to,
		Subject: subjectStr.String(),
		Body:    messageStr.String(),
	}
	msg := BuildMessage(request)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println("mail sent")

}
func ReadIndex(m string) map[string]string {
	nsvs := make(map[string]string)

	stringFirst := strings.Split(m, "namespace:")
	stringSecond := strings.Split(m, "version:")
	for i := 0; i < len(stringFirst)-1; i++ {
		resultNamesp := strings.Split(stringFirst[i+1], " ")
		resultVersion := strings.Split(stringSecond[i+1], " ")

		nsvs[resultNamesp[0]] = resultVersion[0]

	}
	//nsvs contains namespace and key

	return nsvs
}
func ReadConfig(m string) (map[string]string, map[string]string, map[string]string, map[string]string) {
	nsvs := make(map[string]string)
	nsln := make(map[string]string)
	nspl := make(map[string]string)
	nsic := make(map[string]string)
	stringFirst := strings.Split(m, "namespace:")
	stringSecond := strings.Split(m, "version:")
	for i := 0; i < len(stringFirst)-1; i++ {

		resultNamesp := strings.Split(stringFirst[i+1], " ")
		resultVersion := strings.Split(stringSecond[i+1], " ")
		versionString := resultVersion[0]
		nsvs[resultNamesp[0]] = versionString[:len(versionString)-1]

	}

	stringThird := strings.Split(m, "language:")

	for i := 0; i < len(stringFirst)-1; i++ {
		resultNamesp := strings.Split(stringFirst[i+1], " ")
		resultLanguage := strings.Split(stringThird[i+1], " ")
		languageString := resultLanguage[0]
		nsln[resultNamesp[0]] = languageString
	}

	stringFourth := strings.Split(m, "platforms:")

	for i := 0; i < len(stringFirst)-1; i++ {
		resultNamesp := strings.Split(stringFirst[i+1], " ")
		resultPlatform := strings.Split(stringFourth[i+1], "]")
		platformString := resultPlatform[0]
		nspl[resultNamesp[0]] = platformString
	}

	stringFifth := strings.Split(m, "icon:")

	for i := 0; i < len(stringFirst)-1; i++ {
		resultNamesp := strings.Split(stringFirst[i+1], " ")
		resultIcon := strings.Split(stringFifth[i+1], " ")
		iconString := resultIcon[0]
		nsic[resultNamesp[0]] = iconString
	}

	//nsln contains namespace and language
	return nsvs, nsln, nspl, nsic
}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}
