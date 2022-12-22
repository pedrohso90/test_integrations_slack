package main

import (
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"github.com/slack-go/slack"
)

type usersLookupByEmailResponse struct {
	Ok    bool   `json:"ok"`
	User  user   `json:"user"`
	Error string `json:"error"`
}

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {

	var userid string
	
	userid, _ = slackFindUserId()
	_ = setUserInChannel(userid)

}

func slackFindUserId() (string, error){

	token := os.Getenv("SLACK_AUTH_TOKEN")

	userEmail := os.Args[1]

	var client = &http.Client{}
	url := "https://slack.com/api/users.lookupByEmail?email=" + userEmail
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	
	req.Header = http.Header{
		"Authorization": {"Bearer " + token},
	}
	res, err := client.Do(req)
	var data usersLookupByEmailResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
	}
	
	response := data.User.ID
	
	defer res.Body.Close()

	return response, err

}

func setUserInChannel(userid string) (error){

	token := os.Getenv("SLACK_AUTH_TOKEN")
	
	channelID := os.Args[2]

	clientSlack := slack.New(token, slack.OptionDebug(true))

	_, err := clientSlack.InviteUsersToConversation(channelID, userid)
	if err != nil {
		fmt.Println(err)
	}

	return err
}