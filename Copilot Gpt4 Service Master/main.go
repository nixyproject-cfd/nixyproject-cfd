package main

import (
	"github.com/gin-gonic/gin"

	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"copilot-gpt4-service/config"
	"copilot-gpt4-service/log"
	"copilot-gpt4-service/utils"
)

// Handle the Cross-Origin Resource Sharing (CORS) for requests.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

// Represent the JSON data structure for the request body.
type JsonData struct {
	Messages    interface{} `json:"messages"`
	Model       string      `json:"model"`
	Temperature float64     `json:"temperature"`
	TopP        float64     `json:"top_p"`
	N           int64       `json:"n"`
	Stream      bool        `json:"stream"`
}

type Delta struct {
	Content string `json:"content"`
}

type Choice struct {
	Delta Delta `json:"delta"`
	Index int   `json:"index"`
}

type Data struct {
	Choices []Choice `json:"choices,omitempty"`
	Created int      `json:"created,omitempty"`
	ID      string   `json:"id,omitempty"`
}

// Generate a random hexadecimal string of the specified length.
func genHexStr(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

// Create request headers to mock Github Copilot Chat requests.
func createHeaders(authorization string) map[string]string {
	headers := make(map[string]string, 0)
	headers["Authorization"] = "Bearer " + authorization
	headers["X-Request-Id"] = genHexStr(8) + "-" + genHexStr(4) + "-" + genHexStr(4) + "-" + genHexStr(4) + "-" + genHexStr(12)
	headers["Vscode-Sessionid"] = genHexStr(8) + "-" + genHexStr(4) + "-" + genHexStr(4) + "-" + genHexStr(4) + "-" + genHexStr(25)
	headers["Vscode-Machineid"] = genHexStr(64)
	headers["Editor-Version"] = "vscode/1.83.1"
	headers["Editor-Plugin-Version"] = "copilot-chat/0.8.0"
	headers["Openai-Organization"] = "github-copilot"
	headers["Openai-Intent"] = "conversation-panel"
	headers["Content-Type"] = "text/event-stream; charset=utf-8"
	headers["User-Agent"] = "GitHubCopilotChat/0.8.0"
	headers["Accept"] = "*/*"
	headers["Accept-Encoding"] = "gzip,deflate,br"
	headers["Connection"] = "close"

	return headers
}

func chatCompletions(c *gin.Context) {
	content := c.Query("content")
	url := "https://api.githubcopilot.com/chat/completions"

	// Get app token from request header
	appToken, ok := utils.GetAuthorization(c)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "Unauthorized",
			},
		)
		return
	}

	copilotToken, statusCode, errorInfo := utils.GetAuthorizationFromToken(appToken)
	if len(errorInfo) != 0 {
		c.JSON(
			statusCode,
			gin.H{
				"error": errorInfo,
				"code":  statusCode,
			},
		)
		return
	}

	headers := createHeaders(copilotToken)
	jsonBody := &JsonData{
		Messages: []map[string]string{
			{"role": "system",
				"content": "\nYou are ChatGPT, a large language model trained by OpenAI.\nKnowledge cutoff: 2021-09\nCurrent model: gpt-4\n"},
			{"role": "user",
				"content": content},
		},
		Model:       "gpt-4",
		Temperature: 0.5,
		TopP:        1,
		N:           1,
		Stream:      false,
	}
	_ = c.BindJSON(&jsonBody)

	jsonData, err := json.Marshal(jsonBody)
	if err != nil {
		return
	}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Encountering an error when sending the request.")
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return
		} else {
			// Set the headers for the response
			c.Writer.Header().Set("Transfer-Encoding", "chunked")
			c.Writer.Header().Set("X-Accel-Buffering", "no")
			c.Header("Content-Type", "text/event-stream; charset=utf-8")
			c.Header("Cache-Control", "no-cache")
			c.Header("Connection", "keep-alive")
			// Scan the response body line by line
			scanner := bufio.NewScanner(resp.Body)
			for scanner.Scan() {
				line := scanner.Bytes()

				var object string
				if jsonBody.Stream {
					object = "chat.completion.chunk"
				} else {
					object = "chat.completion"
				}

				if bytes.Contains(line, []byte(`"choices":`)) {
					if !bytes.Contains(line, []byte(`"object":`)) {
						line = bytes.ReplaceAll(line, []byte(`"choices":`), []byte(fmt.Sprintf(`"object": "%s", "choices":`, object)))
					}
					if !bytes.Contains(line, []byte(`"model":`)) {
						line = bytes.ReplaceAll(line, []byte(`"choices":`), []byte(fmt.Sprintf(`"model": "%s", "choices":`, jsonBody.Model)))
					}
				}

				c.Writer.Write(line)
				c.Writer.Write([]byte("\n")) // Add newline to the end of each line
				c.Writer.Flush()
			}
			if err := scanner.Err(); err != nil {
				c.AbortWithError(http.StatusBadGateway, err)
				return
			}
		}
	}
}

func createMockModel(modelId string) gin.H {
	return gin.H{
		"id":       modelId,
		"object":   "model",
		"created":  1677610602,
		"owned_by": "openai",
		"permission": []gin.H{
			{
				"id":                   "modelperm-" + genHexStr(12),
				"object":               "model_permission",
				"created":              1677610602,
				"allow_create_engine":  false,
				"allow_sampling":       true,
				"allow_logprobs":       true,
				"allow_search_indices": false,
				"allow_view":           true,
				"allow_fine_tuning":    false,
				"organization":         "*",
				"group":                nil,
				"is_blocking":          false,
			},
		},
		"root":   modelId,
		"parent": nil,
	}
}

func createMockModelsResponse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"object": "list",
		"data": []gin.H{
			createMockModel("gpt-3.5-turbo"),
			createMockModel("gpt-4"),
		},
	})
}

func LoggerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		log.ZLog.Log.Info().Msgf("Request Info:\nMethod: %s\nHost: %s\nURL: %s",
			c.Request.Method, c.Request.Host, c.Request.URL)
		log.ZLog.Log.Debug().Msgf("Request Header:\n%v", c.Request.Header)

		c.Next()

		latency := time.Since(t)
		log.ZLog.Log.Info().Msgf("Response Time: %s\nStatus: %d",
			latency.String(), c.Writer.Status())
		log.ZLog.Log.Debug().Msgf("Response Header:\n%v", c.Writer.Header())
	}
}

func main() {
	if config.ConfigInstance.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.Use(LoggerHandler())

	router.POST("/v1/chat/completions", chatCompletions)
	router.GET("/v1/models", createMockModelsResponse)
	router.GET("/healthz", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "ok",
		})
	})
	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, `非常重要：请不要将此服务公开，仅供个人使用，否则账户或 Copilot 将被封禁。Very important: please do not make this service public, for personal use only, otherwise the account or Copilot will be banned. 非常に重要：このサービスを公開しないでください、個人使用のみにしてください。そうしないと、アカウントまたは Copilot が禁止されます。`)
	})
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
	})

	fmt.Printf("Cache enabled: %t, Cache path: %s, Logging: %t, LOG_LEVEL: %s, Debug: %t\n", config.ConfigInstance.Cache, config.ConfigInstance.CachePath, config.ConfigInstance.Logging, config.ConfigInstance.LogLevel, config.ConfigInstance.Debug)
	fmt.Printf("Starting server on http://%s:%s\n\n", config.ConfigInstance.Host, config.ConfigInstance.Port)

	fmt.Println("\033[31m 非常重要：请不要将此服务公开，仅供个人使用，否则账户或 Copilot 将被封禁。\033[0m")
	fmt.Println("\033[31m Very important: please do not make this service public, for personal use only, otherwise the account or Copilot will be banned.\033[0m")
	fmt.Println("\033[31m 非常に重要：このサービスを公開しないでください、個人使用のみにしてください。そうしないと、アカウントまたは Copilot が禁止されます。\033[0m\n")

	// router.Run(":8080")
	router.Run(config.ConfigInstance.Host + ":" + config.ConfigInstance.Port)
}
