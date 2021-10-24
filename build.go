package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"
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

	var repodirectory string //for mv repo command

	successerror := make(map[string]string)

	// tree walk
	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			//checking path
			configcheck := strings.Contains(path, "config.json")
			distcheck := strings.Contains(path, "dist")
			if configcheck == true && distcheck == false {

				dataconfig, err := ioutil.ReadFile(path)
				if err != nil {
					fmt.Print(err)
				}

				var configintr interface{}
				err = json.Unmarshal(dataconfig, &configintr)
				if err != nil {
					fmt.Println("error:", err)
				}
				stringconfig := fmt.Sprintf("%v", configintr)

				configmap, configlanguage, configplatform, configicon := ReadConfig(stringconfig)
				//path without config.json
				newpath := path[:len(path)-12]
				newpath = filepath.ToSlash(newpath)
				var out bytes.Buffer
				for key1, val1 := range configmap {
					foldername := strings.Split(key1, ".")
					foldernamerepo := foldername[1]
					pythonlanguage := false
					javalanguage := false
					if configlanguage[key1] == "Python" {
						pythonlanguage = true
					} else if configlanguage[key1] == "Java" {
						javalanguage = true
					}
					runneros := os.Getenv("RUNNEROS")
					iconcheck := true
					var commandstr *exec.Cmd
					if configicon[key1] == "icon.png" {
						iconcheck = false
						if runneros == "Linux" {
							commandstr = exec.Command("sh", "-c", " cd "+newpath+" && rclone sync -P --s3-acl=public-read ./icon.png DO:robomotion-packages/contrib/"+foldernamerepo+"/ && sed -i 's+icon.png+https://packages.robomotion.io/contrib/"+foldernamerepo+"/icon.png+' config.json ")
						} else if runneros == "Windows" {
							commandstr = exec.Command("sh", "-c", " cd "+newpath+" && D:/a/packages-contrib/packages-contrib/rclone.exe sync -P --s3-acl=public-read ./icon.png DO:robomotion-packages/contrib/"+foldernamerepo+"/ && sed -i 's+icon.png+https://packages.robomotion.io/contrib/"+foldernamerepo+"/icon.png+' config.json ")
						} else if runneros == "macOS" {
							commandstr = exec.Command("sh", "-c", " cd "+newpath+" && rclone sync -P --s3-acl=public-read ./icon.png DO:robomotion-packages/contrib/"+foldernamerepo+"/ && perl -i -pe's+icon.png+https://packages.robomotion.io/contrib/"+foldernamerepo+"/icon.png+' config.json ")
						}
						commandstr.Stdout = &out
						err = commandstr.Run()
						if err != nil {
							fmt.Println("rclone commands error", err)
						}

					}
					fmt.Println("Icon check", iconcheck)
					platform := configplatform[key1]
					var platformcheck bool
					roboctl := "roboctl"

					fmt.Println(runneros)
					if runneros == "Linux" {
						runneros = "linux"
						platformcheck = strings.Contains(platform, runneros)
						repodirectory = "/home/runner/work/packages-contrib/packages-contrib/repo"
					} else if runneros == "Windows" {
						runneros = "windows"
						platformcheck = strings.Contains(platform, runneros)
						roboctl = "D:/a/packages-contrib/packages-contrib/roboctl"
						repodirectory = "D:/a/packages-contrib/packages-contrib/repo"
					} else if runneros == "macOS" {
						repodirectory = "/Users/runner/work/packages-contrib/packages-contrib/repo"
						runneros = "darwin"
						platformcheck = strings.Contains(platform, runneros)
					}

					//change config

					if strings.Contains(platform, "any") == true {
						platformcheck = true
					}

					fmt.Println(platformcheck)
					fmt.Println(platform)
					exist := true
					_, err = os.Stat("repo/" + foldernamerepo + "")
					if os.IsNotExist(err) {
						exist = false
					}
					fmt.Println(exist)
					changelogcommand := exec.Command("sh", "-c", " cd "+newpath+" && rclone sync -P --s3-acl=public-read ./CHANGELOG.md DO:robomotion-packages/contrib/"+foldernamerepo+"/ ")
					if platformcheck {
						if value, ok := indexmap[key1]; ok {
							//version compare
							v1, err := version.NewVersion(val1)
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
									comman = exec.Command("sh", "-c", " cd "+newpath+" && "+roboctl+" package -b "+runneros+" && mkdir "+foldernamerepo+" && mv robomotion-* "+foldernamerepo+"  && mv "+foldernamerepo+" "+repodirectory+"/"+foldernamerepo+"")
									if pythonlanguage {
										comman = exec.Command("sh", "-c", " cd "+newpath+" && pipreqs . && pip install -r requirements.txt && "+roboctl+" package -b "+runneros+" && mkdir "+foldernamerepo+" && mv robomotion-* "+foldernamerepo+"  && mv "+foldernamerepo+" "+repodirectory+"/"+foldernamerepo+"")
									} else if javalanguage {
										comman = exec.Command("sh", "-c", " cd "+newpath+" && mv ../../OpenJDK11U-* . && "+roboctl+" package -b "+runneros+" && mkdir "+foldernamerepo+" && mv robomotion-* "+foldernamerepo+"  && mv "+foldernamerepo+" "+repodirectory+"/"+foldernamerepo+"")
									}

								} else {
									comman = exec.Command("sh", "-c", " cd "+newpath+" && "+roboctl+" package -b "+runneros+" && mkdir "+foldernamerepo+" && mv robomotion-* "+foldernamerepo+"  && mv "+foldernamerepo+"/* "+repodirectory+"/"+foldernamerepo+"")
									if pythonlanguage {
										comman = exec.Command("sh", "-c", " cd "+newpath+" && pipreqs . && pip install -r requirements.txt && "+roboctl+" package -b "+runneros+" && mkdir "+foldernamerepo+" && mv robomotion-* "+foldernamerepo+"  && mv "+foldernamerepo+"/* "+repodirectory+"/"+foldernamerepo+"")
									} else if javalanguage {
										comman = exec.Command("sh", "-c", " cd "+newpath+" && mv ../../OpenJDK11U-* . && "+roboctl+" package -b "+runneros+" && mkdir "+foldernamerepo+" && mv robomotion-* "+foldernamerepo+"  && mv "+foldernamerepo+"/* "+repodirectory+"/"+foldernamerepo+"")
									}
								}

								comman.Stdout = &out
								err = comman.Run()
								if err != nil {
									fmt.Println("error", err)
									successerror[key1] = err.Error()
								} else {
									successerror[key1] = out.String()
								}
								changelogcommand.Stdout = &out
								err = changelogcommand.Run()
								if err != nil {
									fmt.Println("changelog command error", err)
								}
							}
						} else {

							var cmd *exec.Cmd
							if exist == false {
								cmd = exec.Command("sh", "-c", " cd "+newpath+" && "+roboctl+" package -b "+runneros+" && mkdir "+foldernamerepo+" && mv robomotion-* "+foldernamerepo+" && mv "+foldernamerepo+" "+repodirectory+"/"+foldernamerepo+"")
								if pythonlanguage {
									cmd = exec.Command("sh", "-c", " cd "+newpath+" && pipreqs . && pip install -r requirements.txt && "+roboctl+" package -b "+runneros+" && mkdir "+foldernamerepo+" && mv robomotion-* "+foldernamerepo+" && mv "+foldernamerepo+" "+repodirectory+"/"+foldernamerepo+"")
								} else if javalanguage {
									cmd = exec.Command("sh", "-c", " cd "+newpath+" && mv ../../OpenJDK11U-* . && "+roboctl+" package -b "+runneros+" && mkdir "+foldernamerepo+" && mv robomotion-* "+foldernamerepo+" && mv "+foldernamerepo+" "+repodirectory+"/"+foldernamerepo+"")
								}
							} else {

								cmd = exec.Command("sh", "-c", " cd "+newpath+" && "+roboctl+" package -b "+runneros+" && mkdir "+foldernamerepo+" && mv robomotion-* "+foldernamerepo+" && mv "+foldernamerepo+"/* "+repodirectory+"/"+foldernamerepo+"")
								if pythonlanguage {
									cmd = exec.Command("sh", "-c", " cd "+newpath+" && pipreqs . && pip install -r requirements.txt && "+roboctl+" package -b "+runneros+" && mkdir "+foldernamerepo+" && mv robomotion-* "+foldernamerepo+" && mv "+foldernamerepo+"/* "+repodirectory+"/"+foldernamerepo+"")
								} else if javalanguage {
									cmd = exec.Command("sh", "-c", " cd "+newpath+" && mv ../../OpenJDK11U-* . && "+roboctl+" package -b "+runneros+" && mkdir "+foldernamerepo+" && mv robomotion-* "+foldernamerepo+" && mv "+foldernamerepo+"/* "+repodirectory+"/"+foldernamerepo+"")
								}
							}

							cmd.Stdout = &out
							err = cmd.Run()
							if err != nil {
								fmt.Println("error", err)
								successerror[key1] = err.Error()
							} else {
								successerror[key1] = out.String()
							}

							changelogcommand.Stdout = &out
							err = changelogcommand.Run()
							if err != nil {
								fmt.Println("changelog command error", err)
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
	messagestr := strings.Builder{}
	subjectstr := strings.Builder{}
	subjectstr.WriteString("Contrib Repo's Deployment   ")
	messagestr.WriteString("RUNNER OS= " + os.Getenv("RUNNEROS") + " ")
	for key, value := range successerror {
		messagestr.WriteString(key + value)
		if strings.Contains(value, "Done.") {
			subjectstr.WriteString("" + key + "-> Successful ")
		} else {
			subjectstr.WriteString("" + key + "-> Failed ")
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
		Subject: subjectstr.String(),
		Body:    messagestr.String(),
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

	stringfirst := strings.Split(m, "namespace:")
	stringsecond := strings.Split(m, "version:")
	for i := 0; i < len(stringfirst)-1; i++ {
		resultnamesp := strings.Split(stringfirst[i+1], " ")
		resultversion := strings.Split(stringsecond[i+1], " ")

		nsvs[resultnamesp[0]] = resultversion[0]

	}
	//nsvs contains namespace and key

	return nsvs
}
func ReadConfig(m string) (map[string]string, map[string]string, map[string]string, map[string]string) {
	nsvs := make(map[string]string)
	nsln := make(map[string]string)
	nspl := make(map[string]string)
	nsic := make(map[string]string)
	stringfirst := strings.Split(m, "namespace:")
	stringsecond := strings.Split(m, "version:")
	for i := 0; i < len(stringfirst)-1; i++ {

		resultnamesp := strings.Split(stringfirst[i+1], " ")
		resultversion := strings.Split(stringsecond[i+1], " ")
		versionstring := resultversion[0]
		nsvs[resultnamesp[0]] = versionstring[:len(versionstring)-1]

	}

	stringthird := strings.Split(m, "language:")

	for i := 0; i < len(stringfirst)-1; i++ {
		resultnamesp := strings.Split(stringfirst[i+1], " ")
		resultlanguage := strings.Split(stringthird[i+1], " ")
		languagestring := resultlanguage[0]
		nsln[resultnamesp[0]] = languagestring
	}

	stringfourth := strings.Split(m, "platforms:")

	for i := 0; i < len(stringfirst)-1; i++ {
		resultnamesp := strings.Split(stringfirst[i+1], " ")
		resultplatform := strings.Split(stringfourth[i+1], "]")
		platformstring := resultplatform[0]
		nspl[resultnamesp[0]] = platformstring
	}

	stringfifth := strings.Split(m, "icon:")

	for i := 0; i < len(stringfirst)-1; i++ {
		resultnamesp := strings.Split(stringfirst[i+1], " ")
		resulticon := strings.Split(stringfifth[i+1], " ")
		iconstring := resulticon[0]
		nsic[resultnamesp[0]] = iconstring
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
