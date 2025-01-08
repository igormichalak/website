package main

import (
	"fmt"
	"strings"
	"time"
)

var StartupTime = time.Now()

func generateSitemap() string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

	sb.WriteString("<url>")
	sb.WriteString(fmt.Sprintf("<loc>%s</loc>", BaseURL))
	sb.WriteString(fmt.Sprintf("<lastmod>%s</lastmod>", StartupTime.Format("2006-01-02")))
	sb.WriteString("<changefreq>weekly</changefreq>")
	sb.WriteString("<priority>1.0</priority>")
	sb.WriteString("</url>")

	for _, writing := range Writings {
		locValue := BaseURL + "/" + writing.Slug
		lastmodValue := writing.LastModified.Format("2006-01-02")

		loc := fmt.Sprintf("<loc>%s</loc>", locValue)
		lastmod := fmt.Sprintf("<lastmod>%s</lastmod>", lastmodValue)

		sb.WriteString("<url>")
		sb.WriteString(loc)
		sb.WriteString(lastmod)
		sb.WriteString("<priority>0.9</priority>")
		sb.WriteString("</url>")
	}

	sb.WriteString("</urlset>\n")
	return sb.String()
}

func generateRSS() string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0"?>`)
	sb.WriteString(`<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">`)
	sb.WriteString("<channel>")

	sb.WriteString("<title>Igor Michalak's Writings</title>")
	sb.WriteString("<link>" + BaseURL + "</link>")
	sb.WriteString("<description>A RSS feed with the Igor Michalak's Writings</description>")
	sb.WriteString("<language>en-us</language>")
	sb.WriteString("<ttl>60</ttl>")
	sb.WriteString(`<atom:link href="` + BaseURL + `/rss.xml" rel="self" type="application/rss+xml" />`)

	for _, writing := range Writings {
		sb.WriteString("<item>")
		sb.WriteString("<title>" + writing.Title + "</title>")
		sb.WriteString(fmt.Sprintf("<link>%s/%s</link>", BaseURL, writing.Slug))
		sb.WriteString("<description>" + writing.Description + "</description>")
		sb.WriteString("<pubDate>" + writing.PublishedAt.Format(time.RFC1123Z) + "</pubDate>")
		sb.WriteString(fmt.Sprintf("<guid>%s/%s</guid>", BaseURL, writing.Slug))
		sb.WriteString("</item>")
	}

	sb.WriteString("</channel></rss>\n")
	return sb.String()
}
