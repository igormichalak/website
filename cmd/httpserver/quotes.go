package main

import (
	"html/template"
	"math/rand/v2"
	"strings"
)

var Quotes = []string{
	"The truth will set you free. -- John 8:32",
	"Let all that you do be done in love. -- 1 Corinthians 16:14",
	"Do not be conformed to this world. -- Romans 12:2",
	"For where your treasure is, there your heart will be also. -- Matthew 6:21",
	"Blessed are the pure in heart. -- Matthew 5:8",
	"If you have faith as small as a mustard seed, you can say " +
		"to this mountain, 'Move from here to here', and it will move. -- Matthew 17:20",
	"Love your neighbor as yourself. -- Mark 12:31",
	"Keep your lives free from the love of money " +
		"and be content with what you have. -- Hebrews 13:5",
	"God opposes the proud but shows favor to the humble. -- James 4:6",
	"God loves a cheerful giver. -- 2 Corinthians 9:7",
	"Love your enemies and pray for those who persecute you. -- Matthew 5:44",
	"Bear one another's burdens. -- Galatians 6:2",
	"The one who is in you is greater than the one who is in the world. -- 1 John 4:4",
	"Even the winds and the waves obey him! -- Matthew 8:27",
	"Be strong and do not give up, for your work will be rewarded. -- 2 Chronicles 15:7",
	"If God is for us, who can be against us? -- Romans 8:31",
	"Faith without works is dead. -- James 2:26",
	"For what will it profit a man if he gains the whole world " +
		"and forfeits his soul. -- Matthew 16:26",
}

func getRandomQuote() template.HTML {
	quoteCount := uint(len(Quotes))
	quoteIdx := rand.UintN(quoteCount)
	quote := Quotes[quoteIdx]

	parts := strings.Split(quote, " --")
	if len(parts) != 2 {
		panic(`quote contains more than two occurrences of "--"`)
	}

	text := "&ldquo;" + parts[0] + "&rdquo; "
	source := strings.ReplaceAll(parts[1], " ", "&nbsp;")
	result := text + `<span class="no-break">&ndash;` + source + "</span>"

	return template.HTML(result)
}
