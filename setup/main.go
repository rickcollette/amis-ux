package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"amis-x/bbscommon"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

var db *sql.DB
var config bbscommon.Config
var app *tview.Application
var mainMenu *tview.List

func main() {
	var err error
	db, err = bbscommon.OpenDatabase("./bbs.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	bbscommon.CreateTables(db)
	config, err = bbscommon.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app = tview.NewApplication()

	mainMenu = tview.NewList().
		AddItem("Configure BBS Settings", "Set up the BBS configuration", '1', func() {
			configureBBS()
		}).
		AddItem("Manage Message Bases", "Create, update, delete message bases", '2', func() {
			manageMessageBases()
		}).
		AddItem("Save and Exit", "Save the configuration and exit", '3', func() {
			bbscommon.SaveConfig("config.json", config)
			app.Stop()
		})

	mainMenu.SetBorder(true).SetTitle("Main Menu").SetTitleAlign(tview.AlignLeft)

	if err := app.SetRoot(mainMenu, true).Run(); err != nil {
		panic(err)
	}
}

func configureBBS() {
	form := tview.NewForm().
		AddInputField("BBS Name", config.BBSName, 20, nil, func(text string) {
			config.BBSName = text
		}).
		AddInputField("Sysop Name", config.SysopName, 20, nil, func(text string) {
			config.SysopName = text
		}).
		AddDropDown("Allow New Users", []string{"true", "false"}, boolToIndex(config.AllowNewUsers), func(option string, index int) {
			config.AllowNewUsers = indexToBool(index)
		}).
		AddInputField("ASCII Folder", config.AsciiFolder, 20, nil, func(text string) {
			config.AsciiFolder = text
		}).
		AddInputField("ATASCII Folder", config.AtasciiFolder, 20, nil, func(text string) {
			config.AtasciiFolder = text
		}).
		AddInputField("ANSI Folder", config.AnsiFolder, 20, nil, func(text string) {
			config.AnsiFolder = text
		}).
		AddInputField("Menus Folder", config.MenusFolder, 20, nil, func(text string) {
			config.MenusFolder = text
		}).
		AddInputField("Executables Folder", config.ExecutablesFolder, 20, nil, func(text string) {
			config.ExecutablesFolder = text
		}).
		AddPasswordField("Sysop Password", config.SysopPassword, 20, '*', func(text string) {
			config.SysopPassword = text
		}).
		AddInputField("Port Number", fmt.Sprintf("%d", config.PortNumber), 20, nil, func(text string) {
			fmt.Sscanf(text, "%d", &config.PortNumber)
		}).
		AddButton("Save", func() {
			bbscommon.SaveConfig("config.json", config)
			app.SetRoot(mainMenu, true)
		}).
		AddButton("Cancel", func() {
			app.SetRoot(mainMenu, true)
		})

	form.SetBorder(true).SetTitle("Configure BBS Settings").SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true)
}

func manageMessageBases() {
	list := tview.NewList()

	messageBases, err := bbscommon.ListMessageBases(db)
	if err != nil {
		log.Fatalf("Failed to list message bases: %v", err)
	}

	for _, mb := range messageBases {
		list.AddItem(fmt.Sprintf("%d. %s (Read: %d, Post: %d)", mb.ID, mb.Name, mb.AccessRead, mb.AccessPost), "", 0, nil)
	}

	list.AddItem("Create New Message Base", "", 'c', func() {
		createMessageBaseForm()
	}).
		AddItem("Update Message Base", "", 'u', func() {
			updateMessageBaseForm()
		}).
		AddItem("Delete Message Base", "", 'd', func() {
			deleteMessageBaseForm()
		}).
		AddItem("Back to Main Menu", "", 'b', func() {
			app.SetRoot(mainMenu, true)
		})

	list.SetBorder(true).SetTitle("Manage Message Bases").SetTitleAlign(tview.AlignLeft)
	app.SetRoot(list, true)
}

func createMessageBaseForm() {
	var form = tview.NewForm()
	form = tview.NewForm().
		AddInputField("Name", "", 20, nil, nil).
		AddInputField("Access Read", "0", 20, nil, nil).
		AddInputField("Access Post", "0", 20, nil, nil).
		AddButton("Save", func() {
			name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
			accessRead, _ := strconv.Atoi(form.GetFormItemByLabel("Access Read").(*tview.InputField).GetText())
			accessPost, _ := strconv.Atoi(form.GetFormItemByLabel("Access Post").(*tview.InputField).GetText())
			bbscommon.CreateMessageBase(db, name, accessRead, accessPost)
			app.SetRoot(mainMenu, true)
		}).
		AddButton("Cancel", func() {
			app.SetRoot(mainMenu, true)
		})

	form.SetBorder(true).SetTitle("Create Message Base").SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true)
}

func updateMessageBaseForm() {
	var form = tview.NewForm()
	form = tview.NewForm().
		AddInputField("ID", "", 20, nil, nil).
		AddInputField("Name", "", 20, nil, nil).
		AddInputField("Access Read", "0", 20, nil, nil).
		AddInputField("Access Post", "0", 20, nil, nil).
		AddButton("Save", func() {
			id, _ := strconv.Atoi(form.GetFormItemByLabel("ID").(*tview.InputField).GetText())
			name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
			accessRead, _ := strconv.Atoi(form.GetFormItemByLabel("Access Read").(*tview.InputField).GetText())
			accessPost, _ := strconv.Atoi(form.GetFormItemByLabel("Access Post").(*tview.InputField).GetText())
			bbscommon.UpdateMessageBase(db, id, name, accessRead, accessPost)
			app.SetRoot(mainMenu, true)
		}).
		AddButton("Cancel", func() {
			app.SetRoot(mainMenu, true)
		})

	form.SetBorder(true).SetTitle("Update Message Base").SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true)
}

func deleteMessageBaseForm() {
	var form = tview.NewForm()
	form = tview.NewForm().
		AddInputField("ID", "", 20, nil, nil).
		AddButton("Delete", func() {
			id, _ := strconv.Atoi(form.GetFormItemByLabel("ID").(*tview.InputField).GetText())
			bbscommon.DeleteMessageBase(db, id)
			app.SetRoot(mainMenu, true)
		}).
		AddButton("Cancel", func() {
			app.SetRoot(mainMenu, true)
		})

	form.SetBorder(true).SetTitle("Delete Message Base").SetTitleAlign(tview.AlignLeft)
	app.SetRoot(form, true)
}

func boolToIndex(b bool) int {
	if b {
		return 0
	}
	return 1
}

func indexToBool(index int) bool {
	return index == 0
}
