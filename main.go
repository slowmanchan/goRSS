package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/grokify/html-strip-tags-go"
	"github.com/mitchellh/go-wordwrap"
	"github.com/urfave/cli"
)

type RssTwoMessage struct {
	XMLName  xml.Name  `xml:"rss"`
	Channels []Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"langugae"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string     `xml:"title"`
	Link        string     `xml:"link"`
	Description string     `xml:"description"`
	PubDate     string     `xml:"pubDate"`
	GuID        string     `xml:"guid"`
	Creator     string     `xml:"dc:creator"`
	Categories  []Category `xml:"category"`
}

type Category struct {
	Content string `xml:",chardata"`
}

func main() {
	app := cli.NewApp()
	app.Name = "show"
	app.Usage = "show the rss index feed for a given rss feed url"
	app.Commands = []cli.Command{
		cli.Command{
			Name:        "show",
			Description: "show the rss",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "feedURL, f",
					Usage: "Feed url needed",
				},
			},
			Action: func(c *cli.Context) error {
				message := fetchRss(c.String("feedURL"))
				printRssMessages(message.Channels)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchRss(rssFeedURL string) *RssTwoMessage {
	res, err := http.Get(rssFeedURL)
	if err != nil {
		log.Println(err)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	message := new(RssTwoMessage)

	err = xml.Unmarshal(bodyBytes, message)
	if err != nil {
		log.Println(err)
	}

	return message
}

func formatItem(item Item) string {
	return fmt.Sprintf(`
%s

%s

%s

%s

%s
`,
		strings.Repeat("-", 80),
		wordwrap.WrapString(color.BlueString(item.Title), 80),
		wordwrap.WrapString(strings.TrimSpace(strip.StripTags(item.Description)), 80),
		color.YellowString(item.Link),
		strings.Repeat("-", 80),
	)
}

func printRssMessages(channels []Channel) {
	for _, channel := range channels {
		for _, item := range channel.Items {
			fmt.Println(formatItem(item))
		}
	}
}
