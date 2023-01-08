package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

var (
	ramFlag  = flag.Bool("r", false, "Clean RAM Cache")
	swapFlag = flag.Bool("s", false, "Reload Swap")
	tempFlag = flag.Bool("t", false, "Clean temp directory")
)

func check(err error) {
	if err != nil {
		color.Red("Ошибка: %s", err)
		os.Exit(67)
	}
}

func main() {
	if rootCheck() {
		flag.Parse()
		doClean()
	} else {
		color.Red("Недостаточно привилегий. Запустите утилиту от имени суперпользователя.")
		os.Exit(67)
	}
}

func rootCheck() bool {
	currentUser, err := user.Current()
	check(err)
	return currentUser.Username == "root"
}

func doClean() {
	displayStatus("Очистка началась")
	if *ramFlag {
		cleanRamCache()
	}
	if *swapFlag {
		restartSwap()
	}
	if *tempFlag {
		cleanTemp()
	}
	displayStatus("Очистка завершена")
}

func displayStatus(info string) {
	now := time.Now().Local()
	nowFormatted := fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())
	color.Green("=============== %s: %s ===============", info, nowFormatted)
}

func restartSwap() {
	fmt.Print("Перезагружаю Swap... ")
	cmd := "swapoff -a && swapon -a"
	err := exec.Command("bash", "-c", cmd).Run()
	check(err)
	fmt.Println("Успешно!")
}

func cleanRamCache() {
	fmt.Print("Очищаю кэш оперативной памяти... ")
	err := os.WriteFile("/proc/sys/vm/drop_caches", []byte("3"), 0)
	check(err)
	fmt.Println("Успешно!")
}

func cleanTemp() {
	fmt.Print("Удаляю временные файлы...")
	dir, err := os.Open("/tmp")
	check(err)
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	check(err)
	for _, name := range names {
		err = os.RemoveAll(filepath.Join("/tmp", name))
		check(err)
	}
	fmt.Println("Успешно!")
}
