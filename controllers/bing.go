package controllers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"memoriz-en/models"
	"memoriz-en/utils"

	"github.com/sclevine/agouti"
)

var token string

func (c QuestionController) BingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		expired, _ := utils.CheckTokenDate(w, r)
		if expired {
			return
		}

		// POSTの解析
		var post models.BingPost
		json.NewDecoder(r.Body).Decode(&post)

		// wordの取得
		word := post.Word

		res := makeRequest(word)
		if res.StatusCode == http.StatusUnauthorized {
			token = getToken()
			res = makeRequest(word)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err.Error())
		}

		w.Write(body)
	}
}

func makeRequest(word string) *http.Response {
	postString := `<speak version='1.0' xml:lang='en-US'><voice xml:lang='en-US' xml:gender='Female' name='en-US-JessaRUS'><prosody rate='-20.00%'>` + word + `</prosody></voice></speak>`

	req, err := http.NewRequest("POST", "https://eastasia.tts.speech.microsoft.com/cognitiveservices/v1?", bytes.NewBuffer([]byte(postString)))
	if err != nil {
		log.Fatal(err.Error())
	}
	req.Header.Add("x-microsoft-outputformat", "audio-16khz-32kbitrate-mono-mp3")
	req.Header.Add("authorization", "Bearer "+token)
	req.Header.Add("content-type", "application/ssml+xml")
	req.Header.Add("content-length", strconv.FormatInt(int64(len(postString)), 10))

	proxyURL, _ := url.Parse("http://" + os.Getenv("VPN_USER") + ":" + os.Getenv("VPN_PASS") + os.Getenv("VPN_HOST") + ":80")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyURL),
		},
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	return res
}

func getToken() string {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
		}),
	)
	defer driver.Stop()

	proxys := agouti.ProxyConfig{
		ProxyType:          "pac",                                                                                           //All type -> {direct|manual|pac|autodetect|system}
		ProxyAutoconfigURL: "http://" + os.Getenv("VPN_USER") + ":" + os.Getenv("VPN_PASS") + os.Getenv("VPN_HOST") + ":80", //This is Your Shadowsocks local pac url
	}
	capabilities := agouti.NewCapabilities().Browser("chrome").Proxy(proxys).Without("javascriptEnabled")

	if err := driver.Start(); err != nil {
		log.Fatal(err.Error())
	}

	page, err := driver.NewPage(agouti.Desired(capabilities))
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := page.Navigate(`https://www.bing.com/translator`); err != nil {
		log.Fatal(err.Error())
	}

	textarea := page.FindByXPath(`//*[@id="tta_input_ta"]`)
	if err = textarea.Fill("init"); err != nil {
		log.Fatal(err.Error())
	}

	script := `return sessionStorage.getItem("TTSR");`
	var token string
	for {
		soundButton := page.FindByXPath(`//*[@id="tta_playiconsrc"]`)
		if err = soundButton.Click(); err != nil {
			log.Fatal(err.Error())
		}
		if err = page.RunScript(script, nil, &token); err != nil {
			log.Fatal(err.Error())
		} else {
			if strings.TrimSpace(token) != "" {
				break
			}
		}
	}

	return token
}
