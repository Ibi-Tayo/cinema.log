package domain

type User struct {
	ID          int64  `json:"id"`
	GithubID    string `json:"githubId"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	ProfilePicURL string `json:"profilePicUrl"`
}