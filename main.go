package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"mystic-case-online-quest/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type OnlineQuest struct {
	app *fiber.App
}

func NewOnlineQuest() *OnlineQuest {
	engine := html.New("./templates/views", ".html")
	engine.AddFuncMap(template.FuncMap{
		"safeHTML": func(val any) template.HTML {
			return template.HTML(fmt.Sprint(val))
		},
		"split": func(val string) []string {
			return strings.Split(val, "|")
		},
		"length": func(val []string) int {
			return len(val)
		},
	})

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	return &OnlineQuest{
		app: app,
	}
}

func (oc *OnlineQuest) StartServer() {
	log.Fatal(oc.app.Listen(fmt.Sprintf("%s:%s", config.Config("MYSTIC_CASE_DOMAIN"), config.Config("MYSTIC_CASE_PORT"))))
}

func main() {
	quest := NewOnlineQuest()

	quest.initSystemHandlers()
	quest.initRoutes()

	go quest.StartServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	err := <-c

	log.Printf("Terminated: %s", err.String())
}
