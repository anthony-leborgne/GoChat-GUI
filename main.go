package main

import (
	"database/sql"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("GoChat GUI")

	db, err := sql.Open("sqlite3", "./gochat.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			content TEXT NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	messageList := widget.NewMultiLineEntry()
	messageList.SetPlaceHolder("Aucun message pour l'instant...")
	messageList.Disable()

	updateMessages(db, messageList)

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Votre nom")
	messageEntry := widget.NewEntry()
	messageEntry.SetPlaceHolder("Écrivez un message...")

	sendButton := widget.NewButton("Envoyer", func() {
		username := usernameEntry.Text
		content := messageEntry.Text
		if username == "" || content == "" {
			return
		}
		saveMessage(db, username, content)
		updateMessages(db, messageList)
		messageEntry.SetText("")
	})

	content := container.NewVBox(
		widget.NewLabel("GoChat GUI"),
		messageList,
		usernameEntry,
		messageEntry,
		sendButton,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.ShowAndRun()
}

func saveMessage(db *sql.DB, username, content string) {
	_, err := db.Exec("INSERT INTO messages (username, content) VALUES (?, ?)", username, content)
	if err != nil {
		log.Println("Erreur lors de l'insertion du message :", err)
	}
}

func updateMessages(db *sql.DB, messageList *widget.Entry) {
	rows, err := db.Query("SELECT username, content, timestamp FROM messages ORDER BY timestamp ASC")
	if err != nil {
		log.Println("Erreur lors de la récupération des messages :", err)
		return
	}
	defer rows.Close()

	var messages string
	for rows.Next() {
		var username, content, timestamp string
		err := rows.Scan(&username, &content, &timestamp)
		if err != nil {
			log.Println("Erreur lors de la lecture d'un message :", err)
			continue
		}
		messages += fmt.Sprintf("[%s] %s: %s\n", timestamp, username, content)
	}
	messageList.SetText(messages)
}
