package main

type commentListItem struct {
	User struct {
		Login string `json:"login"`
	} `json:"user"`
	Body     string `json:"body"`
	Location string `json:"html_url"`
}

func (c *commentListItem) isTarget(user string) bool {
	return c.User.Login == user
}
