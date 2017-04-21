package user

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"globaltedinc/mahjongtester/platform"

	"github.com/golang/glog"
)

type hashInfoT struct {
	Key string `json:"key"`
	IV  string `json:"iv"`
}

type gameInfoT struct {
	AppID       string `json:"app_id"`
	AppName     string `json:"app_name"`
	UID         uint32 `json:"uid"`
	Username    string `json:"username"`
	Timestamp   int64  `json:"timestamp"`
	UserIP      string `json:"user_ip"`
	AccessToken string `json:"access_token"`
	Extra       string `json:"extra"`
}

type balanceInfoT struct {
	Cash uint64 `json:"cash"`
	Coin uint64 `json:"Coin"`
	Nm   uint64 `json:"nm"`
}

type jsonBalanceT struct {
	Code int              `json:"code"`
	Info balanceInfoT `json:"info"`
}

type userPlatformT struct {
	site     string
	username string
	password string
	uid      uint32
	name     string
	balance  balanceInfoT

	webClient *http.Client
	platform  *platform.PlatformT
}

func (user *userPlatformT) initialize(platform *platform.PlatformT) error {
	user.platform = platform

	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return err
	}
	user.webClient = &http.Client{Jar: jar}

	return nil
}

func (user *userPlatformT) getHashInfo() (*hashInfoT, error) {
	hashInfoResp, err := user.webClient.Get(fmt.Sprintf("http://%s/Api/getHashInfo", user.platform.GetURL()))
	if err != nil {
		return nil, err
	}
	hashInfoData, err := ioutil.ReadAll(hashInfoResp.Body)
	hashInfoResp.Body.Close()
	if err != nil {
		return nil, err
	}
	glog.Info("data:", string(hashInfoData))

	hashInfo := &hashInfoT{}
	err = json.Unmarshal(hashInfoData, &hashInfo)
	if err != nil {
		return nil, err
	}

	return hashInfo, nil
}

func (user *userPlatformT) LoginPlatform(username string, password string) error {
	user.username = username
	user.password = password

	hashInfo, err := user.getHashInfo()
	if err != nil {
		return err
	}

	block, err := aes.NewCipher([]byte(hashInfo.Key))
	if err != nil {
		return err
	}

	mode := cipher.NewCBCEncrypter(block, []byte(hashInfo.IV))
	dst := make([]byte, aes.BlockSize)
	passwordTmp := make([]byte, 16)
	copy(passwordTmp, user.password)
	//glog.Info("password:", passwordTmp)
	mode.CryptBlocks(dst, passwordTmp)
	//glog.Info("dst:", dst)

	encoded := base64.StdEncoding.EncodeToString(dst)
	//glog.Info(encoded)

	postdata := url.Values{
		"username":          []string{user.username},
		"password":          []string{encoded},
		"remember_password": []string{"0"},
		"verify_code":       []string{""},
	}

	resp, err := user.webClient.PostForm(fmt.Sprintf("http://%s/accounts/login", user.platform.GetURL()), postdata)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return err
	}
	//glog.Info("data:", string(data))

	loginResult := make(map[string]string)
	err = json.Unmarshal(data, &loginResult)
	if err != nil {
		return err
	}
	//glog.Info(loginResult)
	if _, ok := loginResult["error_code"]; ok {
		// 注册用户
		postdata := url.Values{
			"reg_username":        []string{user.username},
			"reg_password":        []string{encoded},
			"reg_password_repeat": []string{encoded},
			"new_reg_source":      []string{"0"},
		}
		resp, err := user.webClient.PostForm(fmt.Sprintf("http://%s/accounts/register", user.platform.GetURL()), postdata)
		if err != nil {
			return err
		}

		data, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return err
		}
		glog.Info("data:", string(data))
	}

	return nil
}

func (user *userPlatformT) getGameInfo(gameAppID string) (uint32, string, int64, string, string, error) {

	gameInfoResp, err := user.webClient.Get(fmt.Sprintf("http://%s/game/info/%s", user.platform.GetURL(), gameAppID))
	if err != nil {
		return 0, "", 0, "", "", err
	}
	gameInfoData, err := ioutil.ReadAll(gameInfoResp.Body)
	gameInfoResp.Body.Close()
	if err != nil {
		return 0, "", 0, "", "", err
	}
	glog.Info("gameinfo data:", string(gameInfoData))

	var gameInfo gameInfoT
	err = json.Unmarshal(gameInfoData, &gameInfo)
	if err != nil {
		return 0, "", 0, "", "", err
	}

	//user.uid = gameInfo.UID
	//user.username = gameInfo.Username
	//user.timestamp = gameInfo.Timestamp
	//user.extra = gameInfo.Extra
	//user.token = gameInfo.AccessToken

	return gameInfo.UID, gameInfo.Username, gameInfo.Timestamp, gameInfo.Extra, gameInfo.AccessToken, nil
}

type assetsInResponseT struct {
	Amount    uint64 `json:"amount"`
	AppID     string `json:"app_id"`
	Currency  string `json:"currency"`
	Date      string `json:"date"`
	RoomID    string `json:"room_id"`
	Timestamp uint32 `json:"timestamp"`
	Token     string `jsin:"token"`
	UID       uint64 `json:"uid"`
}

func (user *userPlatformT) assetsIn(gameid string, currency string, amount uint64, roomid uint64, roomDesc string, token string) (*assetsInResponseT, error) {
	url := fmt.Sprintf("http://%s/purchase/assetsIn?app_id=%s&currency=%s&amount=%d&room_id=%d&room_name=test&room_desc=%s&access_token=%s",
		user.platform.GetURL(), gameid, currency, amount, roomid, roomDesc, token)
	//glog.Info(url)
	assetsResp, err := user.webClient.Get(url)
	if err != nil {
		return nil, err
	}
	assetsData, err := ioutil.ReadAll(assetsResp.Body)
	assetsResp.Body.Close()
	if err != nil {
		return nil, err
	}
	//glog.Info("assets data:", string(assetsData))

	assetsInResponse := &assetsInResponseT{}
	err = json.Unmarshal(assetsData, assetsInResponse)
	if err != nil {
		return nil, err
	}

	return assetsInResponse, nil
}

func (user *userPlatformT) getBalance() error {
	balanceResp, err := user.webClient.Get(fmt.Sprintf("http://%s/accounts/getUserBalance", user.platform.GetURL()))
	if err != nil {
		return err
	}
	balanceData, err := ioutil.ReadAll(balanceResp.Body)
	balanceResp.Body.Close()
	if err != nil {
		return err
	}
	//glog.Info("balance data:", string(balanceData))
	var balance jsonBalanceT
	err = json.Unmarshal(balanceData, &balance)
	if err != nil || balance.Code != 0 {
		glog.Error(balance)
		return err
	}

	user.balance = balance.Info

	return nil
}
