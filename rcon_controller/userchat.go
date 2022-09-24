package rconcontroller

type userChatMessage struct {
	Messages []msg `json:"messages"`
}

type msg struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}
