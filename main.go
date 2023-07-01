package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// setting views engine
	engine := html.New("./", ".html")

	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		var input struct {
			NamaGambar string
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		gambar, err := c.FormFile("gambar")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		// mengambil nama file
		fmt.Printf("Nama file: %s", gambar.Filename)

		// mengambil ext file
		splitDots := strings.Split(gambar.Filename, ".")
		ext := splitDots[len(splitDots)-1]
		fmt.Println(ext)

		nameFileBaru := fmt.Sprintf("%s.%s \n", time.Now().Format("2006-01-02-15-04-05"), ext)
		fmt.Println(nameFileBaru)

		// Mengambil ukuran gambar
		fileHeader, _ := gambar.Open()
		defer fileHeader.Close()
		imageConfig, _, err := image.DecodeConfig(fileHeader)
		if err != nil {
			log.Print(err)
		}
		width := imageConfig.Width
		height := imageConfig.Height
		fmt.Printf("width %d \n", width)
		fmt.Printf("height %d \n", height)

		// Membuat folder upload
		folderUpload := filepath.Join(".", "uploads")

		if err := os.MkdirAll(folderUpload, 0770); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		// Menyimpan gambar ke direktori uploads
		if err := c.SaveFile(gambar, "./uploads/"+nameFileBaru); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.Status(500).JSON(fiber.Map{
			"title":       input.NamaGambar,
			"nama_gambar": nameFileBaru,
			"message":     "Gambar berhasil diupload",
		})
	})

	app.Listen(":8087")
}
