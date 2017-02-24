package oauth2

type GithubUser struct {
	Login    string `json:"login"`
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Location string `json:"location"`
}
