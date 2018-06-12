package main

import (
	"fmt"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type MappedBody struct {
	Name string `json:"name"`
}

type Tweet struct {
	Username string `json:"username"`
	Link string `json:"link"`
	Content string `json:"content"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))
	fmt.Println(request.Body)

	var body = MappedBody{}
	arrB := []byte(request.Body)
	if err := json.Unmarshal(arrB, &body); err != nil {
		log.Fatal(err)
	}

	if len(body.Name) <= 0 {
		return events.APIGatewayProxyResponse{Body: "Got an invalid name", StatusCode: 400}, nil
	}

	tl := getTwitts(body.Name)
	returnJSON, err := json.Marshal(tl)
	if err != nil {
		log.Fatal(err)
	}

	return events.APIGatewayProxyResponse{Body: string(returnJSON), StatusCode: 200}, nil
}

func getTwitts(name string) []Tweet {
	var tweets = make([]Tweet, 0)
	targetURL := fmt.Sprintf("https://mobile.twitter.com/search?q=%s&s=typd&x=0&y=0", url.QueryEscape(name))
	res, err := http.Get(targetURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}  // Find the review items
	doc.Find("table[class^='tweet']").Each(func(i int, s *goquery.Selection) {

		rawLink, _ := s.Attr("href")
		link := "https://twitter.com" + rawLink
		username := s.Find(".username").Text()
		content, err := s.Find(".tweet-text").Html()
		if err != nil {
			log.Fatal(err)
		}
		t := Tweet{ Username: strings.TrimSpace(username), Link: strings.TrimSpace(link), Content: content }
		tweets = append(tweets, t)
	})

	return tweets
}

func main() {
	lambda.Start(handleRequest)
}
