package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Postgres struct {
			Host     string `yaml:"host"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Post     string `yaml:"post"`
			Name     string `yaml:"name"`
		} `yaml:"postgres"`
	} `yaml:"database"`
	Filter struct {
		Artist string  `yaml:"artist"`
		Title  string  `yaml:"title"`
		Price  float64 `yaml:"price"`
	} `yaml:"filter"`
}

type Album struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []Album{
	{Title: "Kind of Blue", Artist: "Miles Davis", Price: 39.99},
	{Title: "The Wall", Artist: "Pink Floyd", Price: 49.99},
	{Title: "The Beatles", Artist: "The Beatles", Price: 59.99},
	{Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{Title: "Time Out", Artist: "Dave Brubeck", Price: 19.99},
}

func loadConfig() (Config, error) {
	var config Config
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	return config, err
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	router := gin.Default()

	router.GET("/albmus", func(c *gin.Context) {
		title := c.Query("title")
		artist := c.Query("artisr")
		price := c.Query("price")

		var fillteradAlbums []Album
		for _, album := range albums {
			if title != "" && !strings.Contains(strings.ToLower(album.Title), strings.ToLower(title)) {
				continue
			}
			if artist != "" && !strings.Contains(strings.ToLower(album.Artist), strings.ToLower(artist)) {
				continue
			}
			if price != "" {
				if n, err := strconv.ParseFloat(price, 64); err == nil {
					if album.Price != n {
						continue
					}
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price parameter"})
					return
				}
			}
			fillteradAlbums = append(fillteradAlbums, album)
		}
		c.JSON(http.StatusOK, fillteradAlbums)

	})
	address := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	router.Run(address)
}
