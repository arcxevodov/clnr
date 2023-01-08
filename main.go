package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	ramFlag  = flag.Bool("r", false, localString("flagRam"))
	swapFlag = flag.Bool("s", false, localString("flagSwap"))
	tempFlag = flag.Bool("t", false, localString("flagTemp"))

	localeArg = flag.String("lang", "en", "Set locale")
)

func initLocalizer() *i18n.Localizer {
	var bundle *i18n.Bundle
	var localizer *i18n.Localizer
	if *localeArg == "ru" {
		bundle = i18n.NewBundle(language.Russian)
	} else {
		bundle = i18n.NewBundle(language.English)
	}

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.LoadMessageFile("locales/en.json")
	bundle.LoadMessageFile("locales/ru.json")

	if *localeArg == "ru" {
		localizer = i18n.NewLocalizer(bundle, language.Russian.String())
	} else {
		localizer = i18n.NewLocalizer(bundle, language.English.String())
	}
	return localizer
}

func localString(id string) string {
	localizer := initLocalizer()
	localzeConfig := i18n.LocalizeConfig{
		MessageID: id,
	}
	result, err := localizer.Localize(&localzeConfig)
	check(err)
	return result
}

func check(err error) {
	if err != nil {
		color.Red("%s: %s", localString("error"), err)
		os.Exit(67)
	}
}

func main() {
	if rootCheck() {
		flag.Parse()
		if !*ramFlag && !*swapFlag && !*tempFlag {
			color.Red(localString("flagError"))
			os.Exit(67)
		}
		doClean()
	} else {
		color.Red(localString("noRoot"))
		os.Exit(67)
	}
}

func rootCheck() bool {
	currentUser, err := user.Current()
	check(err)
	return currentUser.Username == "root"
}

func doClean() {
	displayStatus(localString("clearStart"))
	if *ramFlag {
		cleanRamCache()
	}
	if *swapFlag {
		restartSwap()
	}
	if *tempFlag {
		cleanTemp()
	}
	displayStatus(localString("clearEnd"))
}

func displayStatus(info string) {
	now := time.Now().Local()
	nowFormatted := fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())
	color.Green("=============== %s: %s ===============", info, nowFormatted)
}

func restartSwap() {
	fmt.Print(localString("flagSwap"), "... ")
	cmd := "swapoff -a && swapon -a"
	err := exec.Command("bash", "-c", cmd).Run()
	check(err)
	fmt.Println(localString("success"))
}

func cleanRamCache() {
	fmt.Print(localString("flagRam"), "... ")
	err := os.WriteFile("/proc/sys/vm/drop_caches", []byte("3"), 0)
	check(err)
	fmt.Println(localString("success"))
}

func cleanTemp() {
	fmt.Print(localString("flagTemp"), "... ")
	dir, err := os.Open("/tmp")
	check(err)
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	check(err)
	for _, name := range names {
		err = os.RemoveAll(filepath.Join("/tmp", name))
		check(err)
	}
	fmt.Println(localString("success"))
}
