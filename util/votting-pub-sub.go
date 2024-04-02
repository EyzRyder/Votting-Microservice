package util


type Message struct {
    PollOptionID string `json:"pollOptionId"`
    Votes        float64    `json:"votes"`
}

type Subscriber func(message Message)

type VotingPubSub struct {
    channels map[string][]Subscriber
}

func NewVotingPubSub() *VotingPubSub {
	return &VotingPubSub{
		channels: make(map[string][]Subscriber),
	}
}

func (vp *VotingPubSub) Subscribe(pollID string, subscriber Subscriber) {
    if _, ok := vp.channels[pollID]; !ok {
        vp.channels[pollID] = []Subscriber{}
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
