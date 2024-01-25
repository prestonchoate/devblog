package config

import (
	"html/template"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var configInstance *Config

type Link struct {
	Title string
	URL   string
}

type databaseConfig struct {
	Host         string
	Port         string
	DatabaseName string
	User         string
	Password     string
}

type Config struct {
	templatesDirs []string
	templates     *template.Template
	links         []Link
	dbConfig      databaseConfig
}

func GetInstance() *Config {
	if configInstance == nil {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}

		dbConfig := databaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			DatabaseName: getEnv("DB_NAME", "devblog"),
			User:         getEnv("DB_USER", "devblog"),
			Password:     getEnv("DB_PASSWORD", "devblog"),
		}

		configInstance = &Config{
			templatesDirs: []string{
				"templates/layout/*.html",
				"templates/components/*.html",
				"templates/*.html",
			},
			links: []Link{
				{Title: "Home", URL: "/"},
				{Title: "Blog", URL: "/blog"},
				{Title: "Projects", URL: "/projects"},
				{Title: "About", URL: "/about"},
			},
			templates: template.New(""),
			dbConfig:  dbConfig,
		}
		configInstance.loadTemplates()
	}
	return configInstance
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func (c *Config) loadTemplates() *template.Template {
	for _, dir := range c.templatesDirs {
		var err error
		_, err = c.templates.ParseGlob(dir)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println(c.templates.DefinedTemplates())
	return c.templates
}

func (c *Config) GetTemplates() *template.Template {
	return c.templates
}

func (c *Config) GetLinks() []Link {
	return c.links
}

func (c *Config) GetDBConfig() databaseConfig {
	return c.dbConfig
}
