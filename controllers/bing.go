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

		log.Println("first makeRequest")
		res, err := makeRequest(word)
		if err != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, models.ErrorResponse{Message: "音源の取得に失敗"})
			return
		}
		if res.StatusCode == http.StatusUnauthorized {
			log.Println("getToken")
			token = getToken()
			log.Println("second makeRequest")
			res, err = makeRequest(word)
			if err != nil {
				utils.ResponseWithError(w, http.StatusInternalServerError, models.ErrorResponse{Message: "音源の取得に失敗"})
				return
			}
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err.Error())
		}

		w.Write(body)
	}
}

func makeRequest(word string) (*http.Response, error) {
	postString := `<speak version='1.0' xml:lang='en-US'><voice xml:lang='en-US' xml:gender='Female' name='en-US-JessaRUS'><prosody rate='-20.00%'>` + word + `</prosody></voice></speak>`

	log.Println("NewRequest")
	req, err := http.NewRequest("POST", "https://eastasia.tts.speech.microsoft.com/cognitiveservices/v1?", bytes.NewBuffer([]byte(postString)))
	if err != nil {
		log.Fatal(err.Error())
	}
	req.Header.Add("x-microsoft-outputformat", "audio-16khz-32kbitrate-mono-mp3")
	req.Header.Add("authorization", "Bearer "+token)
	req.Header.Add("content-type", "application/ssml+xml")
	req.Header.Add("content-length", strconv.FormatInt(int64(len(postString)), 10))

	log.Println("proxyURL")
	proxyURL, _ := url.Parse("http://" + os.Getenv("VPN_USER") + ":" + os.Getenv("VPN_PASS") + os.Getenv("VPN_HOST") + ":80")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyURL),
		},
	}
	log.Println("client.Do")
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	log.Println("res")
	log.Println(res)
	return res, nil
}

func getToken() string {
	log.Println("agouti.ChromeDriver")
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			"--window-size=300,1200",
			"--disable-dev-shm-usage",
			"no-sandbox",
		}),
	)
	defer driver.Stop()

	log.Println("agouti proxy")
	proxys := agouti.ProxyConfig{
		ProxyType:          "pac",                                                                                           //All type -> {direct|manual|pac|autodetect|system}
		ProxyAutoconfigURL: "http://" + os.Getenv("VPN_USER") + ":" + os.Getenv("VPN_PASS") + os.Getenv("VPN_HOST") + ":80", //This is Your Shadowsocks local pac url
	}
	capabilities := agouti.NewCapabilities().Browser("chrome").Proxy(proxys).Without("javascriptEnabled")

	log.Println("driver start")
	if err := driver.Start(); err != nil {
		log.Fatal(err.Error())
	}

	log.Println("driver NewPage")
	page, err := driver.NewPage(agouti.Desired(capabilities))
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("driver Navigate")
	if err := page.Navigate(`https://www.bing.com/translator`); err != nil {
		log.Fatal(err.Error())
	}

	log.Println("driver Find")
	textarea := page.FindByXPath(`//*[@id="tta_input_ta"]`)
	if err = textarea.Fill("init"); err != nil {
		log.Fatal(err.Error())
	}

	script := `return sessionStorage.getItem("TTSR");`
	var token string
	for {
		log.Println("driver Find 2")
		soundButton := page.FindByXPath(`//*[@id="tta_playiconsrc"]`)
		log.Println("soundButton click")
		if err = soundButton.Click(); err != nil {
			log.Fatal(err.Error())
		}
		log.Println("page run script")
		if err = page.RunScript(script, nil, &token); err != nil {
			log.Fatal(err.Error())
		} else {
			log.Println("trim space")
			if strings.TrimSpace(token) != "" {
				break
			}
		}
	}

	log.Println("token")
	log.Println(token)
	return token
}
