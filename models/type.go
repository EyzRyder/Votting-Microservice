package models

type Option struct {
		Id    string `json:"id"`
		Title string `json:"title"`
        Score float64 `json:"score"`
	}


type Response struct {
		Poll        Poll `json:"poll"`
		PollOptions []Option    `json:"pollOptions"`
	}
