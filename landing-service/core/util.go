package core

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func sendCode(number, code string) error {

	apiURL := "https://apis.aligo.in/send/"
	data := url.Values{}
	data.Set("key", os.Getenv("API_KEY"))
	data.Set("user_id", os.Getenv("USER_ID"))
	data.Set("sender", os.Getenv("SENDER"))
	data.Set("receiver", number)
	data.Set("msg", "인증번호는 ["+code+"]"+" 입니다.")

	// HTTP POST 요청 실행
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		fmt.Printf("HTTP Request Failed: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Println(fmt.Errorf("server returned non-200 status: %d, body: %s", resp.StatusCode, string(body)))

	return nil

}
