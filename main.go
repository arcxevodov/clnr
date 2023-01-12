package main

// #include <unistd.h>
import "C"

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
	// Флаги для выполнения различных функций программы
	// Flags for performing various program functions
	ramFlag  = flag.Bool("r", false, localString("flagRam"))
	swapFlag = flag.Bool("s", false, localString("flagSwap"))
	tempFlag = flag.Bool("t", false, localString("flagTemp"))
	infoFlag = flag.Bool("i", false, localString("flagInfo"))

	// Аргумент для указания языка программы
	// Argument for specifying the language of the program
	localeArg = flag.String("lang", "en", "Set locale")
)

const (
	NoRootError  = 1
	NoFlagsError = 2
	UnknownError = 3
)

func main() {
	flag.Parse()
	if rootCheck() {
		noFlags := !*ramFlag && !*swapFlag && !*tempFlag

		switch {
		case !*infoFlag && noFlags:
			color.Red(localString("flagError"))
			os.Exit(NoFlagsError)
		case *infoFlag:
			total, free, used := getRam()
			fmt.Printf("%s | %s | %s\n", total, free, used)
			if noFlags {
				os.Exit(0)
			}
			fallthrough
		default:
			doClean()
		}
	} else {
		color.Red(localString("noRoot"))
		os.Exit(NoRootError)
	}
}

// Получение полной и доступной оперативной памяти используя sysconf
// Getting total and available RAM using sysconf
func getRam() (string, string, string) {
	bTotal := C.sysconf(C._SC_PHYS_PAGES) * C.sysconf(C._SC_PAGE_SIZE)
	gbTotal := float64(bTotal) / 1024 / 1024 / 1024
	fmtTotal := fmt.Sprintf(localString("totalRam")+"%.1f GB", gbTotal)

	bFree := C.sysconf(C._SC_AVPHYS_PAGES) * C.sysconf(C._SC_PAGE_SIZE)
	gbFree := float64(bFree) / 1024 / 1024 / 1024
	fmtFree := fmt.Sprintf(localString("freeRam")+"%.1f GB", gbFree)

	fmtUsed := fmt.Sprintf(localString("usedRam")+"%.1f GB", gbTotal-gbFree)

	return fmtTotal, fmtFree, fmtUsed
}

// Инициализация локализатора
// Localizer initialization
func initLocalizer() *i18n.Localizer {
	var bundle *i18n.Bundle
	var localizer *i18n.Localizer

	path, err := os.Executable()
	check(err)

	if *localeArg == "ru" {
		bundle = i18n.NewBundle(language.Russian)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
		_, err := bundle.LoadMessageFile(path[:len(path)-5] + "/locales/ru.json")
		check(err)
		localizer = i18n.NewLocalizer(bundle, language.Russian.String())
	} else {
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
		_, err := bundle.LoadMessageFile(path[:len(path)-5] + "/locales/en.json")
		check(err)
		localizer = i18n.NewLocalizer(bundle, language.English.String())
	}
	return localizer
}

// Функция, возвращающая строку в нужном языке по его ID
// A function that returns a string in the desired language by its ID
func localString(id string) string {
	localizer := initLocalizer()
	localzeConfig := i18n.LocalizeConfig{
		MessageID: id,
	}
	result, err := localizer.Localize(&localzeConfig)
	check(err)
	return result
}

// Автоматизация обработки типичных ошибок
// Automation of handling common errors
func check(err error) {
	if err != nil {
		color.Red("%s: %s", localString("error"), err)
		os.Exit(UnknownError)
	}
}

// Проверка пользователя на наличие прав суперпользователя
// Checking if a user has superuser rights
func rootCheck() bool {
	currentUser, err := user.Current()
	check(err)
	return currentUser.Username == "root"
}

// Начать очистку
// Start cleaning
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

// Вывод статуса и даты
// Display status and date
func displayStatus(info string) {
	now := time.Now().Local()
	nowFormatted := fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())
	color.Green("=============== %s: %s ===============", info, nowFormatted)
}

// Функция перезагрузки Swap файла путем вызова Linux команд
// Function to reload the Swap file by calling Linux commands
func restartSwap() {
	fmt.Print(localString("flagSwap"), "... ")
	cmd := "swapoff -a && swapon -a"
	err := exec.Command("bash", "-c", cmd).Run()
	check(err)
	fmt.Println(localString("success"))
}

// Очистка оперативной памяти путем записи в файл drop_caches
// Clear RAM by writing to the drop_caches file
func cleanRamCache() {
	fmt.Print(localString("flagRam"), "... ")
	err := exec.Command("sync").Run()
	check(err)
	err = os.WriteFile("/proc/sys/vm/drop_caches", []byte("3"), 0)
	check(err)
	fmt.Println(localString("success"))
}

// Очистка папки /tmp путем перебора и рекурсивного удаления всех файлов
// Clean up the /tmp folder by looping through and recursively deleting all files
func cleanTemp() {
	fmt.Print(localString("flagTemp"), "... ")
	dir, err := os.Open("/tmp")
	check(err)
	defer func(dir *os.File) {
		err := dir.Close()
		check(err)
	}(dir)
	names, err := dir.Readdirnames(-1)
	check(err)
	for _, name := range names {
		err = os.RemoveAll(filepath.Join("/tmp", name))
		check(err)
	}
	fmt.Println(localString("success"))
}
