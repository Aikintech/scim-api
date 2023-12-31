package utils

import (
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateRandomNumbers(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(crand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func GenerateRandomString(length int) string {
	pool := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := ""
	for i := 0; i < length; i++ {
		result += string(pool[rand.Intn(len(pool))])
	}
	return result
}

func DumpRoutesToFile(app *fiber.App) error {
	routes := app.GetRoutes(true)

	// Open the file for writing
	files := []string{"routes.txt", "routes.json"}
	error := error(nil)

	for _, filename := range files {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		if filename == "routes.txt" {
			for _, route := range routes {
				_, err := file.WriteString(fmt.Sprintf("%s %s\n", route.Method, route.Path))
				if err != nil {
					return err
				}
			}
		}

		if filename == "routes.json" {
			data, _ := json.MarshalIndent(routes, "", "  ")

			_, err = file.Write(data)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Routes dumped to %s\n", filename)
	}

	return error
}

func GetMimeExtension(mime string) string {
	split := strings.Split(mime, "/")
	if len(split) > 1 {
		return split[len(split)-1]
	}

	return ""
}

func StructToPtr(i interface{}) interface{} {
	return &i
}
