package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/twmb/murmur3"
)

const version = "1.0.0"

var fingerprints = map[int32]string{
	99395752:    "slack-instance",
	116323821:   "spring-boot",
	81586312:    "Jenkins",
	-235701012:  "Cnservers LLC",
	743365239:   "Atlassian",
	2128230701:  "Chainpoint",
	-1277814690: "LaCie",
	246145559:   "Parse",
	628535358:   "Atlassian",
	855273746:   "JIRA",
	1318124267:  "Avigilon",
	-305179312:  "Atlassian – Confluence",
	786533217:   "OpenStack",
	432733105:   "Pi Star",
	705143395:   "Atlassian",
	-1255347784: "Angular IO (AngularJS)",
	-1275226814: "XAMPP",
	-2009722838: "React",
	981867722:   "Atlassian – JIRA",
	-923088984:  "OpenStack",
	1405460984:  "pfSense",
	1278323681:  "Gitlab",
	-1010568750: "phpMyAdmin",
	1015545776:  "pfSense",
	1993518473:  "cPanel Login",
	-895890586:  "PLEX Server",
	1544230796:  "cPanel Login",
	1244636413:  "cPanel Login",
	-127886975:  "Metasploit",
	1139788073:  "Metasploit",
	-1235192469: "Metasploit",
	516963061:   "Gitlab",
	-38580010:   "Magento",
	-1437701105: "XAMPP",
	86919334:    "ServiceNow",
	-1015932800: "Ghost (CMS)",
	-1231681737: "Ghost (CMS)",
	1232159009:  "Apple",
	1382324298:  "Apple",
	-1498185948: "Apple",
	-1252041730: "Vue.js",
	180732787:   "Apache Flink",
}

type FaviconResult struct {
	URL  string
	Hash int32
	Err  error
}

func init() {
	log.SetTimeFormat("15:04:05")
	log.SetLevel(log.InfoLevel)
}

func main() {
	showVersion := flag.Bool("v", false, "show version")
	showShodan := flag.Bool("shodan", false, "show Shodan dorks")
	showHelp := flag.Bool("h", false, "show help")
	flag.Parse()

	if *showVersion {
		displayVersion()
		return
	}

	if *showHelp {
		displayHelp()
		return
	}

	urls := readURLs()
	if len(urls) == 0 {
		log.Error("no URLs provided via stdin")
		return
	}

	results := fetchFavicons(urls)

	hashGroups := make(map[int32][]string)
	for _, result := range results {
		if result.Err == nil {
			hashGroups[result.Hash] = append(hashGroups[result.Hash], result.URL)
		}
	}

	if *showShodan {
		displayShodanDorks(hashGroups)
	} else {
		for hash := range hashGroups {
			fmt.Println(hash)
		}
	}
}

func displayHelp() {
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	flagStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

	fmt.Println()
	fmt.Println(successStyle.Render(" example:"))
	fmt.Printf("    cat urls.txt | %s\n", cmdStyle.Render("faviqon"))
	fmt.Printf("    cat urls.txt | %s -shodan\n\n", cmdStyle.Render("faviqon"))

	fmt.Println(successStyle.Render(" options:"))
	fmt.Printf("    %s \tshow Shodan dorks\n", flagStyle.Render("-shodan"))
	fmt.Printf("    %s \t\tshow version\n", flagStyle.Render("-v"))
	fmt.Printf("    %s \t\tshow this help message\n\n", flagStyle.Render("-h"))
}

func displayVersion() {
	highlightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	fmt.Println()
	fmt.Printf("%s %s\n", highlightStyle.Render("faviqon"), dimStyle.Render("v"+version))
	fmt.Println()
}

func readURLs() []string {
	var urls []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			if strings.HasSuffix(url, "/") {
				urls = append(urls, url+"favicon.ico")
			} else {
				urls = append(urls, url+"/favicon.ico")
			}
		}
	}

	return urls
}

func fetchFavicons(urls []string) []FaviconResult {
	results := make([]FaviconResult, len(urls))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 20)

	for i, url := range urls {
		wg.Add(1)
		go func(idx int, u string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			hash, err := fetchFavicon(u)
			baseURL := strings.TrimSuffix(u, "/favicon.ico")

			if err != nil {
				log.Error("failed to fetch", "url", baseURL)
				results[idx] = FaviconResult{URL: baseURL, Err: err}
			} else {
				results[idx] = FaviconResult{URL: baseURL, Hash: hash}
			}
		}(i, url)
	}

	wg.Wait()
	return results
}

func fetchFavicon(url string) (int32, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	hash := murmur3.Sum32([]byte(b64))

	return int32(hash), nil
}

func displayShodanDorks(hashGroups map[int32][]string) {
	for hash := range hashGroups {
		if hash != 0 {
			fmt.Printf("org:\"target\" http.favicon.hash:%d\n", hash)
		}
	}
}
