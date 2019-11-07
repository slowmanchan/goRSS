package main

import "github.com/fatih/color"

func Example_printRss() {
	color.NoColor = true
	printRssMessages(mockRssFeedExpected.Channels)
	// Output: --------------------------------------------------------------------------------
	//
	// How Elizabeth Warren's wealth tax would work
	//
	// One of Massachusetts Senator Elizabeth Warren's key proposals is a tax on the
	// ultra-rich to pay for some of her big-ticket promises. The idea has support but
	// also a long line of critics who question whether it's practical, effective and
	// even legal.
	//
	// https://www.cbc.ca/news/world/elizabeth-warren-wealth-tax-explained-1.5333534?cmp=rss
	//
	// --------------------------------------------------------------------------------
}
