package providers

type GitlabPushPayload struct {
	Kind          string `json:"push"`
	Before        string `json:"before"`
	After         string `json:"after"`
	Ref           string `json:"ref"`
	CheckoutSha   string `json:"checkout_sha"`
	UserId        string `json:"user_id"`
	Name          string `json:"user_name"`
	Username      string `json:"user_username"`
	Email         string `json:"user_email"`
	UserAvaterUrl string `json:"user_avatar"`
	ProjectId     int    `json:"project_id"`
	Project       struct {
		ProjectId          int64  `json:"id"`
		ProjectName        string `json:"name"`
		ProjectDescription string `json:"description"`
		ProjectWebUrl      string `json:"web_url"`
		GitSshUrl          string `json:"git_ssh_url"`
		GitHttpUrl         string `json:"git_http_url"`
		Namespace          string `json:"namespace"`
		VisibilityLevel    string `json:"visibility_level"`
		NamespacePath      string `json:"path_with_namespace"`
		DefaultBranch      string `json:"default_branch"`
		HomepageUrl        string `json:"homepage"`
		ProjectUrl         string `json:"url"`
		ProjectSshUrl      string `json:"ssh_url"`
		ProjectHttpUrl     string `json:"http_url"`
	} `json:"project"`
	Repository struct {
		RepoName            string `json:"name"`
		RepoUrl             string `json:"url"`
		RepoDescription     string `json:"description"`
		RepoHomepageUrl     string `json:"homepage"`
		RepoHttpUrl         string `json:"git_http_url"`
		RepoSshUrl          string `json:"git_shh_url"`
		RepoVisibilityLevel int64  `json:"visibility_level"`
	} `json:"repository"`
	Commits []struct {
		CommitId        string `json:"Commits"`
		CommitMessage   string `json:"message"`
		CommitTimestamp string `json:"timestamp"`
		CommitUrl       string `json:"url"`
		Author          struct {
			AuthorName  string `json:"name"`
			AuthorEmail string `json:"email"`
		} `json:"author"`
	} `json:"commits"`
	TotalCommitsCount int64 `json:"total_commits_count"`
}
