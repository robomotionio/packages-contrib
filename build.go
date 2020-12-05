package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
)

var (
	pkgDirs []string
)

func filter(path string, info os.FileInfo, err error) error {
	if filepath.Base(path) == "config.json" {
		dir := filepath.Dir(path)
		if !funk.Contains(pkgDirs, filepath.Dir(dir)) { // do not include dist directories
			pkgDirs = append(pkgDirs, filepath.Dir(path))
		}
	}
	return nil
}

func main() {

	filepath.Walk("./src", filter)
	err := os.MkdirAll("repo", 0755)
	if err != nil {
		fmt.Printf("error: cannot make repo dir: %v\n", err)
		return
	}

	for _, pkg := range pkgDirs {

		cfg, modTime, err := readConfig(pkg)
		if err != nil {
			fmt.Printf("error: %s config.json not found: %+v\n", pkg, err)
			continue
		}

		pkgFile, err := getPackageName(pkg, cfg)
		if err != nil {
			fmt.Printf("error: %s validation error: %v\n", pkg, err)
			continue
		}

		packagePath := path.Join("./repo", pkgFile)
		if fileExists(packagePath) {
			info, err := os.Stat(packagePath)
			if err != nil {
				fmt.Printf("error: %s cannot get stat: %v\n", pkg, err)
			}

			if info != nil && info.ModTime().Unix() > modTime.Unix() {
				fmt.Printf("* Skipping %s\n", pkg)
				continue
			}
		}

		fmt.Printf("* Building %s\n", pkg)

		//roboctl, _ := exec.LookPath("roboctl.exe")
		cmd := exec.Command("roboctl", "package")
		cmd.Dir = pkg

		out, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf("get stdout error: %v\n", err)
			continue
		}

		err = cmd.Start()
		if err != nil {
			fmt.Printf("build error: %v\n", err)
			continue
		}

		c := make(chan bool)

		go func() {
			for {
				buf := make([]byte, 1024)
				n, err := out.Read(buf)
				if err == io.EOF {
					c <- true
					return
				} else if err != nil {
					fmt.Printf("pipe read error: %v\n", err)
					return
				}

				fmt.Printf(string(buf[:n]))
			}
		}()

		<-c
		fmt.Println("")

		filepath.Walk(pkg, func(p string, info os.FileInfo, err error) error {
			if strings.HasSuffix(p, ".tgz") {
				base := filepath.Base(p)
				srcPkgFile := path.Join(pkg, base)
				dstPkgFile := path.Join("./repo", base)

				err = os.Rename(srcPkgFile, dstPkgFile)
				if err != nil {
					fmt.Printf("mv error: %v\n", err)
				}
			}
			return nil
		})

	}

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getPackageName(pkg string, cfg string) (string, error) {

	// Get plugin version
	version := gjson.Get(cfg, "version")
	if !version.Exists() {
		return "", fmt.Errorf("Invalid %s version", pkg)
	}

	// Validate plugin version
	_, err := semver.Make(version.String())
	if err != nil {
		return "", err
	}

	// Get plugin name
	rName := gjson.Get(cfg, "namespace")
	if !rName.Exists() {
		return "", fmt.Errorf("Invalid %s namespace", pkg)
	}
	name := strings.ToLower(strings.Replace(rName.String(), ".", "-", -1))

	fileName := fmt.Sprintf("%s-%s", name, version.String())

	pl := gjson.Get(cfg, "platforms").Array()
	platforms := funk.Map(pl, func(p gjson.Result) string {
		return strings.ToLower(p.String())
	}).([]string)

	// Compress product folder
	output := fmt.Sprintf("%s-%s.tgz", fileName, getPlatform())
	if platforms[0] == "any" {
		output = fmt.Sprintf("%s.tgz", fileName)
	}

	return output, nil
}

func readConfig(pkg string) (string, time.Time, error) {
	pkg = path.Join(pkg, "config.json")
	dat, err := ioutil.ReadFile(pkg)
	if err != nil {
		return "", time.Unix(0, 0), err
	}

	info, err := os.Stat(pkg)
	if err != nil {
		return "", time.Unix(0, 0), err
	}

	return string(dat), info.ModTime(), nil
}

func getPlatform() string {
	return fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
}

/////////////////////////////////
/*
package repository

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"robomotion/robomotion-cli/util"
	"runtime"
	"strings"

	"github.com/blang/semver"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
)

const (
	BuildDir = "dist"
)

func Package(dir, cfg string, noBuild bool) {

	fmt.Printf("Reading config file...")

	// Read config.json
	dat, err := ioutil.ReadFile(path.Join(dir, cfg))
	if err != nil {
		log.Fatal(err)
	}

	data := string(dat)

	pl := gjson.Get(data, "platforms").Array()
	platforms := funk.Map(pl, func(p gjson.Result) string {
		return strings.ToLower(p.String())
	}).([]string)

	if (funk.Contains(platforms, "any") && len(platforms) > 1) || len(platforms) == 0 {
		log.Fatalln("invalid platforms")
	}

	fmt.Printf("\r%s", strings.Repeat(" ", 32))
	fmt.Printf("\rGetting plugin version...")

	// Get plugin version
	version := gjson.Get(data, "version")
	if !version.Exists() {
		log.Fatalln("invalid version: empty")
	}

	fmt.Printf("\r%s", strings.Repeat(" ", 32))
	fmt.Printf("\rValidating plugin version...")

	// Validate plugin version
	_, err = semver.Make(version.String())
	if err != nil {
		log.Fatalln("invalid version:", err)
	}

	fmt.Printf("\r%s", strings.Repeat(" ", 32))
	fmt.Printf("\rGetting plugin name...")

	// Get plugin name
	rName := gjson.Get(data, "namespace")
	if !rName.Exists() {
		log.Fatalln("plugin name does not exist")
	}
	name := strings.ToLower(strings.Replace(rName.String(), ".", "-", -1))

	fileName := fmt.Sprintf("%s-%s", name, version.String())

	fmt.Printf("\r%s", strings.Repeat(" ", 32))
	fmt.Printf("\rCleaning...")

	tmp := fmt.Sprintf("scripts.%s.clean", runtime.GOOS)
	cleanScript := gjson.Get(data, tmp).Array()
	runScript(cleanScript)

	fmt.Printf("\r%s", strings.Repeat(" ", 32))
	fmt.Printf("\rBuilding plugin...")

	// Build plugin
	if !noBuild {

		os.RemoveAll(BuildDir)
		os.Mkdir(BuildDir, os.ModeDir|0755)

		build := fmt.Sprintf("scripts.%s.build", runtime.GOOS)
		if platforms[0] == "any" {
			build = "scripts.build"
		}

		buildScript := gjson.Get(data, build).Array()
		runScript(buildScript)
	}

	fmt.Printf("\r%s", strings.Repeat(" ", 32))

	fmt.Printf("\rGenerating spec file...")

	// Generate spec file
	run := fmt.Sprintf("scripts.%s.run", runtime.GOOS)
	if platforms[0] == "any" {
		run = "scripts.run"
	}

	runScript := gjson.Get(data, run).String()
	generateSpecFile(runScript, fileName, gjson.Get(data, "name").String(), version.String(), path.Join(dir, BuildDir))

	// Copy config.json
	err = ioutil.WriteFile(path.Join(BuildDir, "config.json"), dat, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\r%s", strings.Repeat(" ", 100))
	fmt.Printf("\rCompressing files...")

	// Compress product folder
	output := fmt.Sprintf("%s-%s.tgz", fileName, getPlatform())
	if platforms[0] == "any" {
		output = fmt.Sprintf("%s.tgz", fileName)
	}

	excludes := []string{cfg , output}
	if util.CompressAll(BuildDir, output, excludes...); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\r%s", strings.Repeat(" ", 32))
	fmt.Printf("\rDone.")
}

func runScript(script []gjson.Result) {

	for _, cmd := range script {
		parts := strings.Split(cmd.String(), " ")
		args := []string{}
		if len(parts) > 1 {
			args = parts[1:]
		}

		command := exec.Command(parts[0], args...)
		err := command.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func generateSpecFile(script, file, pluginName, version, pkg string) string {

	runs := strings.Split(script, "&&")

	args := []string{}
	for _, r := range runs {
		r = strings.Trim(r, " ")
		args = append(args, strings.Split(r, " ")...)
		args = append(args, "&&")
	}

	args[len(args)-1] = "-s"
	args = append(args, pluginName, version)

	var buff bytes.Buffer
	generateCmd := exec.Command("./"+args[0], args[1:]...)
	generateCmd.Stdout = &buff
	generateCmd.Stderr = &buff
	generateCmd.Dir = pkg

	if err := generateCmd.Run(); err != nil {
		log.Fatal(err)
	}

	name := path.Join(fmt.Sprintf("%s.pspec", file))
	name = path.Join(pkg, name)
	ioutil.WriteFile(name, buff.Bytes(), 0644)

	return name
}

func getPlatform() string {
	return fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
}
*/
/////////////////////////////////
