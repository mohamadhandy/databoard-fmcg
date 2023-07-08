package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func FetchFirebaseImage(urlFirebase string) (string, error) {
	url := urlFirebase
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	firebaseImage := FirebaseImage{}
	err = json.Unmarshal(body, &firebaseImage)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s?alt=media&token=%s", urlFirebase, firebaseImage.Token), nil
}
