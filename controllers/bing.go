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

		// 音声のパスを設定
		strid := strconv.Itoa(post.ID)
		path := "audios/" + strid + ".mp3"

		// サーバにキャッシュされているか確認
		if isSoundCached(path) {
			// ファイル読み込み
			sound, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseWithError(w, http.StatusInternalServerError, models.ErrorResponse{Message: "音声の読み込みに失敗"})
				return
			}
			w.Write(sound)
			return
		}

		// トークンが空なら、トークンを取得
		if token == "" {
			t, err := getToken()
			token = t
			if err != nil {
				log.Println(err.Error())
				utils.ResponseWithError(w, http.StatusInternalServerError, models.ErrorResponse{Message: "トークンの取得に失敗"})
				return
			}
		}

		// キャッシュされていなければ、音声をbing/translatorから取得
		res, err := getAudioSource(post.Word)
		if err != nil {
			log.Println(err.Error())
			utils.ResponseWithError(w, http.StatusInternalServerError, models.ErrorResponse{Message: "音源の取得に失敗"})
			return
		}
		if res.StatusCode == http.StatusUnauthorized {
			token, err = getToken()
			if err != nil {
				log.Println(err.Error())
				utils.ResponseWithError(w, http.StatusInternalServerError, models.ErrorResponse{Message: "トークンの取得に失敗"})
				return
			}

			res, err = getAudioSource(post.Word)
			if err != nil {
				log.Println(err.Error())
				utils.ResponseWithError(w, http.StatusInternalServerError, models.ErrorResponse{Message: "音源の取得に失敗"})
				return
			}
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err.Error())
			utils.ResponseWithError(w, http.StatusInternalServerError, models.ErrorResponse{Message: "ボディの読み込みに失敗"})
			return
		}

		// 音声をキャッシュ
		ioutil.WriteFile(path, body, 0664)

		w.Write(body)
	}
}

func getAudioSource(word string) (*http.Response, error) {
	postString := `<speak version='1.0' xml:lang='en-US'><voice xml:lang='en-US' xml:gender='Female' name='en-US-JessaRUS'><prosody rate='-20.00%'>` + word + `</prosody></voice></speak>`

	req, err := http.NewRequest("POST", "https://southeastasia.tts.speech.microsoft.com/cognitiveservices/v1?", bytes.NewBuffer([]byte(postString)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-microsoft-outputformat", "audio-16khz-32kbitrate-mono-mp3")
	req.Header.Add("authorization", "Bearer "+token)
	req.Header.Add("content-type", "application/ssml+xml")
	req.Header.Add("content-length", strconv.FormatInt(int64(len(postString)), 10))

	proxyURL, err := url.Parse("http://" + os.Getenv("VPN_RAW_USER") + ":" + os.Getenv("VPN_RAW_PASS") + os.Getenv("VPN_HOST") + ":80")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyURL),
		}}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return res, nil
}

func getToken() (string, error) {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			"--window-size=300,1200",
			"--disable-dev-shm-usage",
			"no-sandbox",
		}),
	)
	defer driver.Stop()

	proxys := agouti.ProxyConfig{
		ProxyType:          "pac",                                                                                           //All type -> {direct|manual|pac|autodetect|system}
		ProxyAutoconfigURL: "http://" + os.Getenv("VPN_USER") + ":" + os.Getenv("VPN_PASS") + os.Getenv("VPN_HOST") + ":80", //This is Your Shadowsocks local pac url
	}
	capabilities := agouti.NewCapabilities().Browser("chrome").Proxy(proxys).Without("javascriptEnabled")

	if err := driver.Start(); err != nil {
		return "", err
	}

	page, err := driver.NewPage(agouti.Desired(capabilities))
	if err != nil {
		return "", err
	}
	if err := page.Navigate(`https://www.bing.com/translator`); err != nil {
		return "", err
	}

	textarea := page.FindByXPath(`//*[@id="tta_input_ta"]`)
	if err = textarea.Fill("init"); err != nil {
		return "", err
	}

	script := `return sessionStorage.getItem("TTSR");`
	var token string
	for {
		soundButton := page.FindByXPath(`//*[@id="tta_playicontgt"]`)
		if err = soundButton.Click(); err != nil {
			return "", err
		}
		if err = page.RunScript(script, nil, &token); err != nil {
			return "", err
		} else {
			if strings.TrimSpace(token) != "" {
				break
			}
		}
	}
	return token, nil
}

func isSoundCached(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
