package api

import (
	"time"
	"trade_wizzard/ollama"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func StartApiServer() {

	app := echo.New()

	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	client, _ := XtbClient()
	llvm := ollama.NewOllama()

	app.GET("/chat", func(c echo.Context) error {
		message := c.QueryParam("message")
		llvm.SendMessage(ollama.OllamaMessage{Role: "user", Content: message})

		msg := llvm.Chat()
		return c.JSON(200, msg.Content)
	})

	app.GET("/news", func(c echo.Context) error {
		end := time.Now()
		start := end.Add(-24 * time.Hour)
		data, err := client.GetNews(end, start)

		if len(data.ReturnData) == 0 {
			return c.String(200, "No news found for today")
		}

		if err != nil {
			return c.String(500, "Error getting news")
		}

		for _, news := range data.ReturnData {
			llvm.SendMessage(ollama.OllamaMessage{Role: "user", Content: news.Title + " " + news.Body})
		}

		analysis := llvm.Chat()

		return c.JSON(200, analysis.Content)
	})

	app.GET("/stream/news", func(c echo.Context) error {
		end := time.Now()
		start := end.Add(-24 * time.Hour)
		data, err := client.GetNews(end, start)

		if len(data.ReturnData) == 0 {
			return c.String(200, "No news found for today")
		}

		if err != nil {
			return c.String(500, "Error getting news")
		}

		for _, news := range data.ReturnData {
			llvm.SendMessage(ollama.OllamaMessage{Role: "user", Content: news.Title + "\n" + news.Body})
		}

		responseChan := make(chan ollama.OllamaMessage)
		go llvm.StreamChat(responseChan)

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
		c.Response().WriteHeader(200)

		for msg := range responseChan {
			if _, err := c.Response().Write([]byte(msg.Content)); err != nil {
				return err
			}

			c.Response().Flush()
		}

		return nil
	})

	defer client.Close()

	app.Start(":1420")
}
