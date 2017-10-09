package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Your hipchat URL Please change
const hipchat = "https://HIPCHAT-SERVER/v2/room/YOUR_ROOM/notication"

// Support Atlassian's status page only Please change
const statuspage = "https://STATUSPAGE"

// Your hipcaht key Please change
const hipchatkey = "HIPCHAT KEY"

var StatusPageChildItems = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "status_page",
		Name:      "child_items",
		Help:      "status_page_child_items status",
	},
	[]string{"colo", "status"},
)

func sendMsg(msg string, color string) bool {
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": fmt.Sprintf("%s %s", "Bearer ", hipchatkey),
	}
	param := req.Param{
		"message": msg,
		"notify":  "True",
		"format":  "text",
		"color":   color,
	}
	r, err := req.Post(hipchat, header, param)
	if err != nil {
		log.Fatal(err)
		return false
	}
	log.Printf("%+v", r) // print info (try it, you may surprise)
	return true
}
func linkScrape() map[string][]string {
	doc, err := goquery.NewDocument(statuspage)
	result := make(map[string][]string)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".status-orange").Each(func(index int, item *goquery.Selection) {
		colo := item.Find(".name").Text()
		if strings.Contains(colo, " - ") {
			result["all"] = append(result["all"], strings.TrimSpace(strings.Split(colo, " - ")[0]))
			result["orange"] = append(result["orange"], strings.TrimSpace(strings.Split(colo, " - ")[0]))
		} else if len(colo) > 0 {
			result["all"] = append(result["all"], strings.TrimSpace(colo))
			result["orange"] = append(result["orange"], strings.TrimSpace(colo))
		}
	})
	doc.Find(".status-yellow").Each(func(index int, item *goquery.Selection) {
		colo := item.Find(".name").Text()
		if strings.Contains(colo, " - ") {
			result["all"] = append(result["all"], strings.TrimSpace(strings.Split(colo, " - ")[0]))
			result["yellow"] = append(result["yellow"], strings.TrimSpace(strings.Split(colo, " - ")[0]))
		} else if len(colo) > 0 {
			result["all"] = append(result["all"], strings.TrimSpace(colo))
			result["yellow"] = append(result["yellow"], strings.TrimSpace(colo))
		}
	})
	doc.Find(".status-green").Each(func(index int, item *goquery.Selection) {
		colo := item.Find(".name").Text()
		if strings.Contains(colo, " - ") {
			result["all"] = append(result["all"], strings.TrimSpace(strings.Split(colo, " - ")[0]))
			result["green"] = append(result["green"], strings.TrimSpace(strings.Split(colo, " - ")[0]))
		} else if len(colo) > 0 {
			result["all"] = append(result["all"], strings.TrimSpace(colo))
			result["green"] = append(result["green"], strings.TrimSpace(colo))
		}
	})
	doc.Find(".status-red").Each(func(index int, item *goquery.Selection) {
		colo := item.Find(".name").Text()
		if strings.Contains(colo, " - ") {
			result["all"] = append(result["all"], strings.TrimSpace(strings.Split(colo, " - ")[0]))
			result["red"] = append(result["red"], strings.TrimSpace(strings.Split(colo, " - ")[0]))
		} else if len(colo) > 0 {
			result["all"] = append(result["all"], strings.TrimSpace(colo))
			result["red"] = append(result["red"], strings.TrimSpace(colo))
		}
	})
	return result
}

var oldmetric map[string][]string

func main() {
	flag.Parse()

	log.Printf("Starting Server: %s", "0.0.0.0")
	http.HandleFunc("/metrics", handleMetricsRequest)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
func handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(StatusPageChildItems)
	StatusPageChildItems.Reset()
	result := linkScrape()
	if len(result["all"]) > 0 {
		for i := range result["all"] {
			StatusPageChildItems.WithLabelValues(result["all"][i], "orange").Set(0)
			StatusPageChildItems.WithLabelValues(result["all"][i], "green").Set(0)
			StatusPageChildItems.WithLabelValues(result["all"][i], "yellow").Set(0)
			StatusPageChildItems.WithLabelValues(result["all"][i], "red").Set(0)
		}
	}
	if len(result["orange"]) > 0 {
		msg := "Statuspage partial Outage item [orange] :"
		for i := range result["orange"] {
			msg = fmt.Sprintf("%s %s", msg, result["orange"][i])
		}
		sendMsg(msg, "purple")
	}
	if len(result["yellow"]) > 0 {
		msg := "Statuspage partial Outage item [yellow] :"
		for i := range result["yellow"] {
			msg = fmt.Sprintf("%s %s", msg, result["yellow"][i])
		}
		sendMsg(msg, "yellow")
	}
	if len(result["red"]) > 0 {
		msg := "Statuspage Outage item [red] :"
		for i := range result["red"] {
			msg = fmt.Sprintf("%s %s", msg, result["red"][i])
		}
		sendMsg(msg, "red")
	}
	for i := range result["orange"] {
		StatusPageChildItems.WithLabelValues(result["orange"][i], "orange").Set(1)
	}
	for i := range result["yellow"] {
		StatusPageChildItems.WithLabelValues(result["yellow"][i], "yellow").Set(1)
	}
	for i := range result["red"] {
		StatusPageChildItems.WithLabelValues(result["red"][i], "red").Set(1)
	}
	for i := range result["green"] {
		StatusPageChildItems.WithLabelValues(result["green"][i], "green").Set(1)
	}
	oldmetric = result
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
