package facades

import pusherGo "github.com/pusher/pusher-http-go/v5"

type pusher struct {
	client *pusherGo.Client
}

type PusherData map[string]string

func newPusher() *pusher {
	client := pusherGo.Client{
		AppID:   Env().GetString("PUSHER_APP_ID"),
		Key:     Env().GetString("PUSHER_APP_KEY"),
		Secret:  Env().GetString("PUSHER_APP_SECRET"),
		Cluster: Env().GetString("PUSHER_APP_CLUSTER"),
		Secure:  true,
	}

	return &pusher{client: &client}
}

func Pusher() *pusherGo.Client {
	return newPusher().client
}
