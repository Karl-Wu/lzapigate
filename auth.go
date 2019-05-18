package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

/*
export function* loginRequest() {
  yield takeEvery('LOGIN_REQUEST', function*({payload}) {
    const {username, password} = payload;
      const {response, error} = yield call(loginApi, username, password);
      if (response) {
        console.log('loginRequest:', response);
        yield put({
          type: actions.LOGIN_SUCCESS,
          token: response.token,
          user: response.user,
          payload
        });
      } else {
        console.log('loginRequest error:', error);
        yield put({ type: actions.LOGIN_ERROR, payload });
      }
  });
}
*/

type loginJsonStruct struct {
	Token string `json:"token"`
}

func login(username, password string) error {
	var (
		respJson loginJsonStruct
	)

	client := &http.Client{}
	cred := map[string]string{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(cred)
	resp, err := client.Post(authURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("login error:", err)
		return err
	}
	DEBUG.Println(LOGIN, "login resp:", resp)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	DEBUG.Println(LOGIN, "login resp boday", string(respBody))

	json.Unmarshal(respBody, &respJson)

	DEBUG.Println(LOGIN, "login resp json:", respJson)

	authToken = respJson.Token

	return nil
}
