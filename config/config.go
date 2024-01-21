package config

import (
	"html/template"
	"log"
)

var configInstance *Config

type Link struct {
	Title string
	URL   string
}

type Config struct {
	templatesDirs []string
	templates     *template.Template
	links         []Link
}

func GetInstance() *Config {
	if configInstance == nil {
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
		}
		configInstance.loadTemplates()
	}
	return configInstance
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
