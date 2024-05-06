package learn_fiber

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var app *fiber.App = fiber.New()

func TestHelloWorld(t *testing.T) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	result, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, World!", string(result))
}

func TestHello(t *testing.T) {
	app.Get("/hello", func(ctx *fiber.Ctx) error {
		name := ctx.Query("name", "Guest")
		return ctx.SendString("Hello " + name)
	})

	request, _ := http.NewRequest(http.MethodGet, "/hello?name=Raihan", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)

	result, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Raihan", string(result))
}

func TestHttpRequest(t *testing.T) {
	app.Get("/httptest", func(ctx *fiber.Ctx) error {
		headerName := ctx.Get("name", "Guest")    //header
		cookieFullName := ctx.Cookies("fullname") //cookie

		return ctx.SendString("Hello " + headerName + " or " + cookieFullName)
	})

	request, _ := http.NewRequest(http.MethodGet, "/httptest", nil)
	request.Header.Set("name", "Gaou")
	request.AddCookie(&http.Cookie{Name: "fullname", Value: "Raihanhori"})

	response, err := app.Test(request)
	assert.Nil(t, err)

	result, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Gaou or Raihanhori", string(result))
}

func TestRouteParams(t *testing.T) {
	app.Get("/users/:username/orders/:orderId", func(ctx *fiber.Ctx) error {
		username := ctx.Params("username")
		orderId := ctx.Params("orderId")

		return ctx.SendString("User " + username + " with order id " + orderId)
	})

	request, _ := http.NewRequest(http.MethodGet, "/users/raihanhori/orders/123", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)

	result, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "User raihanhori with order id 123", string(result))
}

func TestFormRequest(t *testing.T) {
	app.Post("/form", func(ctx *fiber.Ctx) error {
		name := ctx.FormValue("name")
		return ctx.SendString("Hello " + name)
	})

	body := strings.NewReader("name=Raihan")
	request, _ := http.NewRequest(http.MethodPost, "/form", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := app.Test(request)
	assert.Nil(t, err)

	result, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Raihan", string(result))
}

//go:embed source/contoh.txt
var contohFile []byte

func TestFileUpload(t *testing.T) {
	app.Post("/upload", func(ctx *fiber.Ctx) error {
		file, err := ctx.FormFile("file")
		if err != nil {
			return err
		}

		// save file
		errSaveFile := ctx.SaveFile(file, "./target/"+file.Filename)
		if errSaveFile != nil {
			return errSaveFile
		}

		return ctx.SendString("Upload file successfully")
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, _ := writer.CreateFormFile("file", "contoh.txt")
	file.Write(contohFile)
	writer.Close()

	request, _ := http.NewRequest(http.MethodPost, "/upload", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := app.Test(request)
	assert.Nil(t, err)
	result, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Upload file successfully", string(result))
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func TestRequestBody(t *testing.T) {
	app.Post("/login", func(ctx *fiber.Ctx) error {
		body := ctx.Body()
		loginRequest := LoginRequest{}

		err := json.Unmarshal(body, &loginRequest)
		if err != nil {
			return err
		}

		return ctx.SendString("Hello " + loginRequest.Email + " - " + loginRequest.Password)
	})

	body := strings.NewReader(`{"email":"raihan@test.com","password":"password"}`)
	request, _ := http.NewRequest(http.MethodPost, "/login", body)
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)
	result, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello raihan@test.com - password", string(result))
}

type RegisterRequest struct {
	Name     string `json:"name" xml:"name" form:"name"`
	Email    string `json:"email" xml:"email" form:"email"`
	Password string `json:"password" xml:"password" form:"password"`
}

func TestBodyParser(t *testing.T) {
	app.Post("/register", func(ctx *fiber.Ctx) error {
		registerRequest := new(RegisterRequest)
		err := ctx.BodyParser(registerRequest)
		if err != nil {
			return err
		}

		return ctx.SendString("Register success " + registerRequest.Name)
	})

	body := strings.NewReader(`{"name":"raihanhori","email":"raihan@test.com","password":"password"}`)
	request, _ := http.NewRequest(http.MethodPost, "/register", body)
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)
	result, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Register success raihanhori", string(result))
}

func TestResponseJson(t *testing.T) {
	app.Get("/user", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"name":  "Raihanhori",
			"email": "raihanki02@gmail.com",
		})
	})

	request, _ := http.NewRequest(http.MethodGet, "/user", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	result, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"email":"raihanki02@gmail.com","name":"Raihanhori"}`, string(result))
}

func TestRouteGroup(t *testing.T) {
	helloWorld := func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello World!")
	}

	api := app.Group("/api")
	api.Get("/hello", helloWorld)
	api.Get("/world", helloWorld)

	request, _ := http.NewRequest(http.MethodGet, "/api/hello", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	result, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello World!", string(result))
}
