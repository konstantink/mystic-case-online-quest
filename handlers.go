package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mystic-case-online-quest/config"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	danger "github.com/iromli/go-itsdangerous"
)

func (oq *OnlineQuest) initRoutes() {
	oq.app.Route("/", func(router fiber.Router) {
		router.Add(fiber.MethodGet, "", AuthMiddleware, oq.homeIndexHandler)
		router.Add(fiber.MethodPost, "", oq.homeAnswerHandler)
		router.Route("/key", func(router fiber.Router) {
			router.Add(fiber.MethodGet, "", AuthMiddleware, oq.keyIndexHandler)
			router.Add(fiber.MethodPost, "", oq.keyAnswerHandler)
			router.Route("/face", func(router fiber.Router) {
				router.Add(fiber.MethodGet, "", AuthMiddleware, oq.faceIndexHandler)
				router.Add(fiber.MethodPost, "", oq.faceAnswerHandler)
				router.Route("/ghost", func(router fiber.Router) {
					router.Add(fiber.MethodGet, "", AuthMiddleware, oq.ghostIndexHandler)
					router.Add(fiber.MethodPost, "", oq.ghostAnswerHandler)
					router.Route("/034", func(router fiber.Router) {
						router.Add(fiber.MethodGet, "", AuthMiddleware, oq.zeroThreeFourIndexHandler)
						router.Add(fiber.MethodPost, "", oq.zeroThreeFourAnswerHandler)
						router.Route("/attention", func(router fiber.Router) {
							router.Add(fiber.MethodGet, "", AuthMiddleware, oq.attentionIndexHandler)
							router.Add(fiber.MethodPost, "", oq.attentionAnswerHandler)
							router.Route("/umbrella", func(router fiber.Router) {
								router.Add(fiber.MethodGet, "", AuthMiddleware, oq.umbrellaIndexHandler)
								router.Add(fiber.MethodPost, "", oq.umbrellaAnswerHandler)
								router.Route("/bishop", func(router fiber.Router) {
									router.Add(fiber.MethodGet, "", AuthMiddleware, oq.bishopIndexHandler)
									router.Add(fiber.MethodPost, "", oq.bishopAnswerHandler)
									router.Route("/cheshire", func(router fiber.Router) {
										router.Add(fiber.MethodGet, "", AuthMiddleware, oq.cheshireIndexHandler)
									}, "level_eight")
								}, "level_seven")
							}, "level_six")
						}, "level_five")
					}, "level_four")
				}, "level_tree")
			}, "level_two")
		}, "level_one")
	}, "home_page")
}

type SessionData struct {
	Role   string
	UserID uuid.UUID
	IsNew  bool
}

// def sign(self, value: str | bytes) -> bytes:
//     """Signs the given string and also attaches time information."""
//     value = want_bytes(value)
//     timestamp = base64_encode(int_to_bytes(self.get_timestamp()))
//     sep = want_bytes(self.sep)
//     value = value + sep + timestamp
//     return value + sep + self.get_signature(value)

type SessionKey string

func AuthMiddleware(c *fiber.Ctx) error {
	var (
		signedSession string
		err           error
	)
	sessionData := new(SessionData)
	// salt := make([]byte, 32)
	// _, err := io.ReadFull(rand.Reader, salt)
	// if err != nil {
	// 	log.Printf("[WARN] failed to generate salt %s", err.Error())
	// }
	signature := danger.NewSignature(config.Config("SECRET_KEY"), "", ".", "", nil, nil)
	session_id := c.Cookies("session", "")
	if session_id != "" {
		log.Printf("[INFO] found session %s", session_id)
		unsigned, err := signature.Unsign(session_id)
		if err != nil {
			log.Printf("[WARN] failed to unsign the cookie %s", err.Error())
		}
		decodedUnsigned, err := base64.StdEncoding.DecodeString(unsigned)
		if err != nil {
			log.Printf("[WARN] failed to decode unsigned cookie %s", err.Error())
		}
		err = json.Unmarshal([]byte(decodedUnsigned), sessionData)
		if err != nil {
			log.Printf("[WARN] failed to unmarshal cookie %s", err.Error())
		}
		if sessionData.IsNew {
			sessionData.IsNew = false
		}
	} else {
		userID, _ := uuid.NewUUID()
		sessionData = &SessionData{
			Role:   "user",
			UserID: userID,
			IsNew:  true,
		}
		serializedSession, _ := json.Marshal(sessionData)
		encodedSession := base64.StdEncoding.EncodeToString(serializedSession)
		signedSession, err = signature.Sign(encodedSession)
		if err != nil {
			log.Printf("[WARN] failed to sign cookie value %s", err.Error())
		}
		log.Printf("[INFO] Signed cookie %s", signedSession)
		session_id = signedSession
	}

	c.SetUserContext(context.WithValue(c.UserContext(), SessionKey("userSession"), sessionData))
	err = c.Next()

	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    session_id,
		Path:     "/",
		Domain:   "mcoq.slipkid.xyz", // config.Config("MYSTIC_CASE_DOMAIN"),
		MaxAge:   int(time.Hour) * 24 * 30,
		Expires:  time.Now().Add(time.Hour * 24),
		Secure:   false,
		HTTPOnly: true,
		SameSite: "lax",
	})

	return err
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

func (oq *OnlineQuest) homeIndexHandler(c *fiber.Ctx) error {
	var context = Context{
		config.Config("MYSTIC_CASE_GTM_ID"),
		"Introduction",
		"/",
		nil,
	}
	sessionData := c.UserContext().Value(SessionKey("userSession")).(*SessionData)
	log.Printf("[INFO] Visiting home page %s", sessionData.UserID.String())
	oq.botChan <- BotMessage{
		UserID: sessionData.UserID,
		IsNew:  sessionData.IsNew,
		Page:   "Introduction",
	}
	return c.Render("home", context, "layouts/base")
}

func (oq *OnlineQuest) homeAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, "/key", "key")
}

func (oq *OnlineQuest) keyIndexHandler(c *fiber.Ctx) error {
	var (
		hints   = getHints("key")
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Level #1",
			"/key",
			hints,
		}
	)
	sessionData := c.UserContext().Value(SessionKey("userSession")).(*SessionData)
	log.Printf("[INFO] Visiting level 1 page %s", sessionData.UserID.String())
	oq.botChan <- BotMessage{
		UserID: sessionData.UserID,
		IsNew:  sessionData.IsNew,
		Page:   "Level #1",
	}
	return c.Render("key", context, "layouts/base")
}

func (oq *OnlineQuest) keyAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, fmt.Sprintf("%s/face", c.Path()), "face")
}

func (oq *OnlineQuest) faceIndexHandler(c *fiber.Ctx) error {
	var (
		hints   = getHints("face")
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Level #2",
			"/key/face",
			hints,
		}
	)
	sessionData := c.UserContext().Value(SessionKey("userSession")).(*SessionData)
	log.Printf("[INFO] Visiting level 2 page %s", sessionData.UserID.String())
	oq.botChan <- BotMessage{
		UserID: sessionData.UserID,
		IsNew:  sessionData.IsNew,
		Page:   "Level #2",
	}
	return c.Render("face", context, "layouts/base")
}

func (oq *OnlineQuest) faceAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, fmt.Sprintf("%s/ghost", c.Path()), "ghost")
}

func (oq *OnlineQuest) ghostIndexHandler(c *fiber.Ctx) error {
	var (
		hints   = getHints("ghost")
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Level #3",
			"/key/face/ghost",
			hints,
		}
	)
	sessionData := c.UserContext().Value(SessionKey("userSession")).(*SessionData)
	log.Printf("[INFO] Visiting level 3 page %s", sessionData.UserID.String())
	oq.botChan <- BotMessage{
		UserID: sessionData.UserID,
		IsNew:  sessionData.IsNew,
		Page:   "Level #3",
	}
	return c.Render("ghost", context, "layouts/base")
}

func (oq *OnlineQuest) ghostAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, fmt.Sprintf("%s/034", c.Path()), "034", "043", "304", "340", "430", "403")
}

func (oq *OnlineQuest) zeroThreeFourIndexHandler(c *fiber.Ctx) error {
	var (
		hints   = getHints("034")
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Level #4",
			"/key/face/ghost/034",
			hints,
		}
	)
	sessionData := c.UserContext().Value(SessionKey("userSession")).(*SessionData)
	log.Printf("[INFO] Visiting level 4 page %s", sessionData.UserID.String())
	oq.botChan <- BotMessage{
		UserID: sessionData.UserID,
		IsNew:  sessionData.IsNew,
		Page:   "Level #4",
	}
	return c.Render("034", context, "layouts/base")
}

func (oq *OnlineQuest) zeroThreeFourAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, fmt.Sprintf("%s/attention", c.Path()), "attention")
}

func (oq *OnlineQuest) attentionIndexHandler(c *fiber.Ctx) error {
	var (
		hints   = getHints("attention")
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Level #5",
			"/key/face/ghost/034/attention",
			hints,
		}
	)
	sessionData := c.UserContext().Value(SessionKey("userSession")).(*SessionData)
	log.Printf("[INFO] Visiting level 5 page %s", sessionData.UserID.String())
	oq.botChan <- BotMessage{
		UserID: sessionData.UserID,
		IsNew:  sessionData.IsNew,
		Page:   "Level #5",
	}
	return c.Render("attention", context, "layouts/base")
}

func (oq *OnlineQuest) attentionAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, fmt.Sprintf("%s/umbrella", c.Path()), "umbrella")
}

func (oq *OnlineQuest) umbrellaIndexHandler(c *fiber.Ctx) error {
	var (
		hints   = getHints("umbrella")
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Level #6",
			"/key/face/ghost/034/attention/umbrella",
			hints,
		}
	)
	sessionData := c.UserContext().Value(SessionKey("userSession")).(*SessionData)
	log.Printf("[INFO] Visiting level 6 page %s", sessionData.UserID.String())
	oq.botChan <- BotMessage{
		UserID: sessionData.UserID,
		IsNew:  sessionData.IsNew,
		Page:   "Level #6",
	}
	return c.Render("umbrella", context, "layouts/base")
}

func (oq *OnlineQuest) umbrellaAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, fmt.Sprintf("%s/bishop", c.Path()), "bishop")
}

func (oq *OnlineQuest) bishopIndexHandler(c *fiber.Ctx) error {
	var (
		hints   = getHints("bishop")
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Level #7",
			"/key/face/ghost/034/attention/umbrella/bishop",
			hints,
		}
	)
	sessionData := c.UserContext().Value(SessionKey("userSession")).(*SessionData)
	log.Printf("[INFO] Visiting level 7 page %s", sessionData.UserID.String())
	oq.botChan <- BotMessage{
		UserID: sessionData.UserID,
		IsNew:  sessionData.IsNew,
		Page:   "Level #7",
	}
	return c.Render("bishop", context, "layouts/base")
}

func (oq *OnlineQuest) bishopAnswerHandler(c *fiber.Ctx) error {
	return checkPasscode(c, fmt.Sprintf("%s/cheshire", c.Path()), "cheshire", "cheshire cat", "cheshirecat", "dinah")
}

func (oq *OnlineQuest) cheshireIndexHandler(c *fiber.Ctx) error {
	var (
		context = Context{
			config.Config("MYSTIC_CASE_GTM_ID"),
			"Congratulations",
			"/key/face/ghost/034/attention/umbrella/bishop/cheshire",
			nil,
		}
	)
	sessionData := c.UserContext().Value(SessionKey("userSession")).(*SessionData)
	log.Printf("[INFO] Visiting feedback page %s", sessionData.UserID.String())
	oq.botChan <- BotMessage{
		UserID: sessionData.UserID,
		IsNew:  sessionData.IsNew,
		Page:   "Feedback",
	}
	return c.Render("cheshire", context, "layouts/base")
}
