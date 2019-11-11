package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/grokify/html-strip-tags-go"
	"github.com/mitchellh/go-wordwrap"
	"github.com/urfave/cli"
)

type configFile struct {
	Version  string
	RssFeeds []*rssFeeds
}

type rssFeeds struct {
	Name       string `xml:"Name"`
	RssFeedURL string `xml:"RssFeedURL"`
}

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
			Name:        "list",
			Description: "list all subscribed feeds",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "feedURL, f",
					Usage: "Feed url Needed",
				},
			},
			Action: func(c *cli.Context) error {
				listAll()
				return nil
			},
		},
		cli.Command{
			Name:        "show",
			Description: "show the rss",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "rssName, r",
					Usage: "saved rss name",
				},
			},
			Action: func(c *cli.Context) error {
				message, err := fetchRss(c.String("rssName"))
				if err != nil {
					return err
				}
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

func fetchRss(inputRssName string) (*RssTwoMessage, error) {
	configFileData, err := ioutil.ReadFile(".config.json")
	if err != nil {
		return nil, err
	}

	rssFeedURL := ""
	configFile := &configFile{}
	if err := json.Unmarshal(configFileData, configFile); err != nil {
		return nil, err
	}

	for _, rssFeed := range configFile.RssFeeds {
		if rssFeed.Name == inputRssName {
			rssFeedURL = rssFeed.RssFeedURL
		}
	}

	if rssFeedURL == "" {
		return nil, errors.New("No rss feeds of that name were found")
	}

	res, err := http.Get(rssFeedURL)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	message := new(RssTwoMessage)

	err = xml.Unmarshal(bodyBytes, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func printRssMessagesToScreen(stdin io.Writer, item Item) {
	fmt.Fprintf(stdin, fmt.Sprintf(`
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
	))
}

func printRssMessages(channels []Channel) {
	cmd := exec.Command("less", "-r")
	r, stdin := io.Pipe()
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	c := make(chan struct{})
	go func() {
		defer close(c)
		cmd.Run()
	}()

	for _, channel := range channels {
		for _, item := range channel.Items {
			printRssMessagesToScreen(stdin, item)
		}
	}
	stdin.Close()
	<-c
}

func listAll() error {
	data, err := ioutil.ReadFile(".config.json")
	if err != nil {
		return err
	}
	configFile := &configFile{}
	if err := json.Unmarshal(data, configFile); err != nil {
		return err
	}

	for _, rssFeed := range configFile.RssFeeds {
		fmt.Printf("%s | %s\n", rssFeed.Name, rssFeed.RssFeedURL)
	}
	return nil
}
