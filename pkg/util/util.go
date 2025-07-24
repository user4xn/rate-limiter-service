package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func GetEnv(key string, fallback string) string {
	a, _ := godotenv.Read()
	var (
		val     string
		isExist bool
	)

	val, isExist = a[key]
	if !isExist {
		val = os.Getenv(key)
		if val == "" {
			val = fallback
		}
	}
	return val
}

func GracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]func(ctx context.Context) error) <-chan struct{} {
	wait := make(chan struct{})

	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
		<-s

		log.Println("shutting down service")

		defer close(s)

		timeoutFunc := time.AfterFunc(timeout, func() {
			fmt.Printf("timeout %vs has been elapsed, forced exit", timeout.Seconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		for key, op := range ops {
			log.Printf("cleaning up: %s", key)

			if err := op(ctx); err != nil {
				fmt.Printf("cleaning up %s failed: %s", key, err.Error())
				return
			}
		}

		close(wait)
	}()

	return wait
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse

}

func CreateErrorLog(errMessage error) {
	fmt.Println(errMessage.Error())
	fileName := fmt.Sprintf("./storage/error_logs/error-%s.log", time.Now().Format("2006-01-02"))

	// open log file
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}

	defer logFile.Close()

	// set log out put
	log.SetOutput(logFile)

	log.SetFlags(log.LstdFlags)

	_, fileName, line, _ := runtime.Caller(1)
	log.Printf("[Error] in [%s:%d] %v", fileName, line, errMessage.Error())
}

func LoadDefaultWindow() int {
	raw := GetEnv("FIX_WINDOW_DEFAULT_DURATION", "60")
	window, err := strconv.Atoi(raw)
	if err != nil {
		log.Printf("invalid FIX_WINDOW_DEFAULT_DURATION: %v, fallback to 60", err)
		CreateErrorLog(err)
		return 60
	}
	return window
}

func LoadDefaultMaxRequest() int {
	raw := GetEnv("FIX_WINDOW_DEFAULT_MAX_REQUEST", "100")
	maxRequest, err := strconv.Atoi(raw)
	if err != nil {
		log.Printf("invalid FIX_WINDOW_DEFAULT_MAX_REQUEST: %v, fallback to 100", err)
		CreateErrorLog(err)
		return 100
	}
	return maxRequest
}
