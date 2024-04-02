package util


type Message struct {
    PollOptionID string `json:"pollOptionId"`
    Votes        int    `json:"votes"`
}

type Subscriber func(message Message)

type VotingPubSub struct {
    channels map[string][]Subscriber
}

func (vp *VotingPubSub) Subscribe(pollID string, subscriber Subscriber) {
    if vp.channels == nil {
        vp.channels = make(map[string][]Subscriber)
    }
    vp.channels[pollID] = append(vp.channels[pollID], subscriber)
}

func (vp *VotingPubSub) Publish(pollID string, message Message) {
    subscribers, ok := vp.channels[pollID]
    if !ok {
        return
    }
    for _, subscriber := range subscribers {
        subscriber(message)
    }
}
