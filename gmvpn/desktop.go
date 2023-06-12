/*
 * @Author: liuwei lyy9645@163.com
 * @Date: 2023-05-20 21:59:49
 * @LastEditors: liuwei lyy9645@163.com
 * @LastEditTime: 2023-05-21 15:57:17
 * @FilePath: /gmvpn/ui/desktop.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"gmvpn/common"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

var configDir string

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	var cfg_file *os.File
	var err error
	if !common.FileIsExisted("gmvpn.cfg") {
		cfg_file, err = os.Create("gmvpn.cfg")
	} else {
		cfg_file, err = os.OpenFile("gmvpn.cfg", os.O_RDWR, 6)
	}
	if err != nil {
		logrus.Warnln("打开或创建配置文件失败", err)
		return
	}
	defer cfg_file.Close()

	t := time.Now().Format("2006-01-02 15:04:05")
	_, err = cfg_file.WriteString(t)
	if err != nil {
		logrus.Warnln("配置文件写入数据失败")
	}
	if !common.FileIsExisted("cert") {
		err = common.MakeDir("cert")
		if err != nil {
			logrus.Warnln("创建cert目录失败", err)
			return
		}

		caCert, _ := os.OpenFile("cert/ca.crt", os.O_RDWR|os.O_CREATE, 6)
		caCert.WriteString(string(resourceEnccertCrt.StaticContent))
		encCert, _ := os.OpenFile("cert/encert.crt", os.O_RDWR|os.O_CREATE, 6)
		encCert.WriteString(string(resourceEnccertCrt.StaticContent))
		encKey, _ := os.OpenFile("cert/enckey.key", os.O_RDWR|os.O_CREATE, 6)
		encKey.WriteString(string(resourceEnccertCrt.StaticContent))
		signCert, _ := os.OpenFile("cert/signcert.crt", os.O_RDWR|os.O_CREATE, 6)
		signCert.WriteString(string(resourceEnccertCrt.StaticContent))
		signKey, _ := os.OpenFile("cert/signkey.key", os.O_RDWR|os.O_CREATE, 6)
		signKey.WriteString(string(resourceEnccertCrt.StaticContent))
	}

}

func main() {
	a := app.New()
	// 设置主题
	a.Settings().SetTheme(&myTheme{})

	w := a.NewWindow("Gmvpn")
	w.Resize(fyne.NewSize(340, 150))
	w.SetMaster()
	w.SetFixedSize(true)

	var indexCon *fyne.Container

	serverEntry := widget.NewEntry()
	serverEntry.SetPlaceHolder("例如: 122.231.117.70:3001")
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("例如: zhangsan")
	passwordEntry := widget.NewPasswordEntry()

	infinite := widget.NewProgressBarInfinite()

	form := &widget.Form{}

	form.Append("IP地址:", serverEntry)
	form.Append("账 户:", usernameEntry)
	form.Append("密 码:", passwordEntry)

	loginBtn := widget.NewButton("登录", func() {
		c := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), widget.NewLabel("登录中..."), layout.NewSpacer())
		w.SetContent(fyne.NewContainerWithLayout(
			layout.NewVBoxLayout(),
			layout.NewSpacer(),
			infinite,
			c,
			layout.NewSpacer(),
		))
		time.Sleep(3 * time.Second)
		w.SetContent(
			indexCon,
		)
		log.Println("登录了。。。")
	})

	masterLaout := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		form,
		loginBtn,
	)

	w.SetContent(
		masterLaout,
	)
	logoutBtn := widget.NewButton("退出", func() {
		log.Println("退出了。。。")
		w.SetContent(masterLaout)
	})
	indexCon = fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		widget.NewLabel("登录成功"),
		logoutBtn,
	)
	w.ShowAndRun()
}
