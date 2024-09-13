package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mystic-case-online-quest/config"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (oq *OnlineQuest) initRoutes() {
	oq.app.Route("/", func(router fiber.Router) {
		router.Add(fiber.MethodGet, "", homeIndexHandler)
		router.Add(fiber.MethodPost, "", homeAnswerHandler)
	}, "home_page")

	oq.app.Route("/key", func(router fiber.Router) {
		router.Add(fiber.MethodGet, "", keyIndexHandler)
		router.Add(fiber.MethodPost, "", keyAnswerHandler)
		router.Route("/face", func(router fiber.Router) {
			router.Add(fiber.MethodGet, "", faceIndexHandler)
			router.Add(fiber.MethodPost, "", faceAnswerHandler)
			router.Route("/ghost", func(router fiber.Router) {
				router.Add(fiber.MethodGet, "", ghostIndexHandler)
				router.Add(fiber.MethodPost, "", ghostAnswerHandler)
				router.Route("/034", func(router fiber.Router) {
					router.Add(fiber.MethodGet, "", zeroThreeFourIndexHandler)
					router.Add(fiber.MethodPost, "", zeroThreeFourAnswerHandler)
					router.Route("/attention", func(router fiber.Router) {
						router.Add(fiber.MethodGet, "", attentionIndexHandler)
						router.Add(fiber.MethodPost, "", attentionAnswerHandler)
						router.Route("/umbrella", func(router fiber.Router) {
							router.Add(fiber.MethodGet, "", umbrellaIndexHandler)
							router.Add(fiber.MethodPost, "", umbrellaAnswerHandler)
							router.Route("/bishop", func(router fiber.Router) {
								router.Add(fiber.MethodGet, "", bishopIndexHandler)
								router.Add(fiber.MethodPost, "", bishopAnswerHandler)
								router.Route("/cheshire", func(router fiber.Router) {
									router.Add(fiber.MethodGet, "", cheshireIndexHandler)
								})
							})
						})
					})
				})
			}, "level_tree")
		}, "level_two")
	}, "level_one")
}

func (oq *OnlineQuest) initSystemHandlers() {
	oq.app.Static("/favicon.ico", "./images/favicon.ico")
	oq.app.Static("/sitemap.xml", "./sitemap.xml")
	oq.app.Static("/robots.txt", "./robots.txt")
	oq.app.Static("/static/", "./static")
	oq.app.Static("/images/", "./assets/images")
}

type (
	Passcode struct {
		Passcode string `json:"passcode" form:"passcode"`
	}

	Answer struct {
		Correct  bool
		NextPath string
		Passcode string
	}

	Hint struct {
		Name string `json:"name"`
		Text string `json:"text"`
	}

	Context struct {
		GtmID string
		Title string
		Path  string
		Hints *[]Hint
	}
)

func checkPasscode(c *fiber.Ctx, nextPath string, answers ...string) error {
	var p = Passcode{}
	if err := c.BodyParser(&p); err != nil {
		log.Print("missing passcode")
		c.Status(http.StatusBadRequest)
		return errors.New("missing passcode")
	}

	for i := range answers {
		if strings.EqualFold(p.Passcode, answers[i]) {
			return c.Render("answer", &Answer{true, nextPath, p.Passcode})
		}
	}

	return c.Render("answer", &Answer{false, "", p.Passcode})
}

func getHints(page string) *[]Hint {
	var hints = new([]Hint)
	hintsStream, _ := os.ReadFile(fmt.Sprintf("./hints/%s.json", page))
	err := json.Unmarshal(hintsStream, hints)
	if err != nil {
		log.Print(err.Error())
	}

	return hints
}

func homeIndexHandler(c *fiber.Ctx) error {
	var context = Context{
		config.Config("MYSTIC_CASE_GTM_ID"),
		"Introduction",
		"/",
		nil,
	}
	return c.Render("home", context, "layouts/base")
}

func homeAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, "/key", "key")
}

func keyIndexHandler(c *fiber.Ctx) error {
	var (
		hints   = getHints("key")
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Level #1",
			"/key",
			hints,
		}
	)
	return c.Render("key", context, "layouts/base")
}

func keyAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, fmt.Sprintf("%s/face", c.Path()), "face")
}

func faceIndexHandler(c *fiber.Ctx) error {
	var (
		hints   = getHints("face")
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Level #2",
			"/key/face",
			hints,
		}
	)
	return c.Render("face", context, "layouts/base")
}

func faceAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, fmt.Sprintf("%s/ghost", c.Path()), "ghost")
}

func ghostIndexHandler(c *fiber.Ctx) error {
	var (
		hints   = getHints("ghost")
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Level #3",
			"/key/face/ghost",
			hints,
		}
	)
	return c.Render("ghost", context, "layouts/base")
}

func ghostAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, fmt.Sprintf("%s/034", c.Path()), "034", "043", "304", "340", "430", "403")
}

func zeroThreeFourIndexHandler(c *fiber.Ctx) error {
	var (
		hints   = getHints("034")
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Level #4",
			"/key/face/ghost/034",
			hints,
		}
	)
	return c.Render("034", context, "layouts/base")
}

func zeroThreeFourAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, fmt.Sprintf("%s/attention", c.Path()), "attention")
}

func attentionIndexHandler(c *fiber.Ctx) error {
	var (
		hints   = getHints("attention")
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Level #5",
			"/key/face/ghost/034/attention",
			hints,
		}
	)
	return c.Render("attention", context, "layouts/base")
}

func attentionAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, fmt.Sprintf("%s/umbrella", c.Path()), "umbrella")
}

func umbrellaIndexHandler(c *fiber.Ctx) error {
	var (
		hints   = getHints("umbrella")
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Level #6",
			"/key/face/ghost/034/attention/umbrella",
			hints,
		}
	)
	return c.Render("umbrella", context, "layouts/base")
}

func umbrellaAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, fmt.Sprintf("%s/bishop", c.Path()), "bishop")
}

func bishopIndexHandler(c *fiber.Ctx) error {
	var (
		hints   = getHints("bishop")
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Level #7",
			"/key/face/ghost/034/attention/umbrella/bishop",
			hints,
		}
	)
	return c.Render("bishop", context, "layouts/base")
}

func bishopAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, fmt.Sprintf("%s/cheshire", c.Path()), "cheshire", "cheshire cat", "cheshirecat", "dinah")
}

func cheshireIndexHandler(c *fiber.Ctx) error {
	var (
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Congratulations",
			"/key/face/ghost/034/attention/umbrella/bishop/cheshire",
			nil,
		}
	)
	return c.Render("cheshire", context, "layouts/base")
}
