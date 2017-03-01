package oauth2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type GithubUser struct {
	Login    string `json:"login"`
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type BaiduUser struct {
	UID      string `json:"uid"`
	Uname    string `json:"uname"`
	Portrait string `json:"portrait"`
}

func GetGithubUser(client *http.Client) (GithubUser, error) {
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Fatal(err)
		return GithubUser{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
	var uj GithubUser
	if err = json.Unmarshal(body, &uj); err != nil {
		log.Fatal(err)
		return uj, err
	}
	fmt.Println("login:" + uj.Login + "  ,   name:" + uj.Name)
	return uj, nil
}

func GetBaiduUser(accessToken string) (BaiduUser, error) {

	if accessToken == "" {
		return BaiduUser{}, errors.New("accessToken is not allow nil")
	}
	bu := BaiduUser{}
	url := "https://openapi.baidu.com/rest/2.0/passport/users/getLoggedInUser?access_token=" + accessToken
	err := getJson(url, &bu)
	if err != nil {
		return bu, errors.New("convert userinfo json to struct fail")
	}
	return bu, nil
}

func getJson(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err = json.Unmarshal(body, &target); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
