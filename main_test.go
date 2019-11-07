package main

import (
	"encoding/xml"
	"testing"

	"github.com/go-test/deep"
	"github.com/jarcoal/httpmock"
)

func Test_fetchRss(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://www.test.com/articles",
		httpmock.NewStringResponder(200, mockRssFeedResponse))

	expected := mockRssFeedExpected
	actual := fetchRss("https://www.test.com/articles")

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Error(diff)
	}
}

var mockRssFeedResponse = `
<?xml version="1.0" encoding="UTF-8"?>
<rss xmlns:cbc="https://www.cbc.ca/rss/cbc" version="2.0">
  <channel>
    <title><![CDATA[CBC | Top Stories News]]></title>
    <link>http://www.cbc.ca/news/?cmp=rss</link>
    <description><![CDATA[FOR PERSONAL USE ONLY]]></description>
    <language>en-ca</language>
    <lastBuildDate>Sat, 2 Nov 2019 17:11:35 EDT</lastBuildDate>
    <copyright><![CDATA[Copyright: (C) Canadian Broadcasting Corporation, http://www.cbc.ca/aboutcbc/discover/termsofuse.html#Rss]]></copyright>
        <docs><![CDATA[http://www.cbc.ca/rss/]]></docs>
    <image>
      <title>CBC.ca</title>
      <url>https://www.cbc.ca/rss/image/cbc_144.gif</url>
      <link>https://www.cbc.ca/news/?cmp=rss</link>
    </image>
    <item cbc:type="story" cbc:deptid="2.633" cbc:syndicate="true">
      <title><![CDATA[How Elizabeth Warren's wealth tax would work]]></title>
            <link>https://www.cbc.ca/news/world/elizabeth-warren-wealth-tax-explained-1.5333534?cmp=rss</link>
      <guid isPermaLink="false">1.5333534</guid>
      <pubDate>Fri, 1 Nov 2019 13:16:54 EDT</pubDate>
                        <author>Steven D&apos;Souza</author>
                                                              <category>News/World</category>
      <description><![CDATA[<img src='https://i.cbc.ca/1.5229831.1564451554!/fileImage/httpImage/image.jpg_gen/derivatives/16x9_460/1163038492.jpg' alt='1163038492' width='460' title='SIOUX CITY, IOWA - JULY 19: Democratic presidential hopeful U.S. Sen. Elizabeth Warren (D-MA) speaks during the AARP and The Des Moines Register Iowa Presidential Candidate Forum on July 19, 2019 in Sioux City, Iowa. Twenty democratic presidential hopefuls are participating in the AARP and Des Moines Register candidate forums that will feature four candidates per forum that are being to be held in cities across Iowa over five days. ' height='259' />                <p>One of Massachusetts Senator Elizabeth Warren's key proposals is a tax on the ultra-rich to pay for some of her big-ticket promises. The idea has support but also a long line of critics who question whether it's practical, effective and even legal.</p>]]></description>
    </item>
  </channel>
</rss>
`

// https://www.theverge.com/apps/rss/index.xml

var mockRssFeedExpected = &RssTwoMessage{
	XMLName: xml.Name{"", "rss"},
	Channels: []Channel{
		Channel{
			Title:       "CBC | Top Stories News",
			Link:        "http://www.cbc.ca/news/?cmp=rss",
			Description: `FOR PERSONAL USE ONLY`,
			Language:    "",
			Items: []Item{
				Item{
					Title:       "How Elizabeth Warren's wealth tax would work",
					Link:        "https://www.cbc.ca/news/world/elizabeth-warren-wealth-tax-explained-1.5333534?cmp=rss",
					Description: `<img src='https://i.cbc.ca/1.5229831.1564451554!/fileImage/httpImage/image.jpg_gen/derivatives/16x9_460/1163038492.jpg' alt='1163038492' width='460' title='SIOUX CITY, IOWA - JULY 19: Democratic presidential hopeful U.S. Sen. Elizabeth Warren (D-MA) speaks during the AARP and The Des Moines Register Iowa Presidential Candidate Forum on July 19, 2019 in Sioux City, Iowa. Twenty democratic presidential hopefuls are participating in the AARP and Des Moines Register candidate forums that will feature four candidates per forum that are being to be held in cities across Iowa over five days. ' height='259' />                <p>One of Massachusetts Senator Elizabeth Warren's key proposals is a tax on the ultra-rich to pay for some of her big-ticket promises. The idea has support but also a long line of critics who question whether it's practical, effective and even legal.</p>`,
					PubDate:     "Fri, 1 Nov 2019 13:16:54 EDT",
					GuID:        "1.5333534",
					Categories: []Category{
						Category{
							Content: "News/World",
						},
					},
				},
			},
		},
	},
}
