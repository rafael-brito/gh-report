package reports

import "time"

type RepositoryRef struct {
	Owner string
	Name  string
}

type FileHistoryMode string

const (
	FileHistoryModeCommits FileHistoryMode = "commits"
	FileHistoryModePRs     FileHistoryMode = "prs"
)

type FileHistoryParams struct {
	Repo   RepositoryRef
	File   string
	Limit  int
	Mode   FileHistoryMode
	UserID string // para chave de cache, se necessário
}

type FileHistoryCommit struct {
	SHA            string    `json:"commit_sha"`
	Message        string    `json:"commit_message"`
	URL            string    `json:"commit_url"`
	AuthorLogin    string    `json:"author_login"`
	AuthorAvatar   string    `json:"author_avatar_url"`
	CommittedAt    time.Time `json:"committed_at"`
	AssociatedPRNo *int      `json:"pr_number,omitempty"`
}

type FileHistoryPREntry struct {
	PRNumber    int                 `json:"pr_number"`
	PRTitle     string              `json:"pr_title"`
	PRURL       string              `json:"pr_url"`
	PRMergedAt  *time.Time          `json:"pr_merged_at,omitempty"`
	Commits     []FileHistoryCommit `json:"commits"`
	AuthorLogin string              `json:"author_login"`
}

type FileHistoryEntryType string

const (
	FileHistoryEntryTypeCommit FileHistoryEntryType = "commit"
	FileHistoryEntryTypePR     FileHistoryEntryType = "pr"
)

type FileHistoryEntry struct {
	Type    FileHistoryEntryType `json:"type"`
	Commit  *FileHistoryCommit   `json:"commit,omitempty"`
	PR      *FileHistoryPREntry  `json:"pr,omitempty"`
	OrderTs time.Time            `json:"order_ts"`
}

type TopAuthorStat struct {
	Login   string `json:"login"`
	Commits int    `json:"commits"`
}

type FileHistoryStats struct {
	TotalCommits int             `json:"total_commits"`
	TotalPRs     int             `json:"total_prs"`
	TopAuthors   []TopAuthorStat `json:"top_authors"`
}

type FileHistoryReport struct {
	Repository  RepositoryRef      `json:"repository"`
	FilePath    string             `json:"file_path"`
	Mode        FileHistoryMode    `json:"mode"`
	Limit       int                `json:"limit"`
	GeneratedAt time.Time          `json:"generated_at"`
	Entries     []FileHistoryEntry `json:"entries"`
	Stats       FileHistoryStats   `json:"stats"`
}

// ----- Release diff -----

type ReleaseDiffParams struct {
	Repo   RepositoryRef
	From   string
	To     string
	UserID string // para cache
}

type ReleasePRType string

const (
	ReleasePRTypeFeature     ReleasePRType = "feature"
	ReleasePRTypeBugfix      ReleasePRType = "bugfix"
	ReleasePRTypeImprovement ReleasePRType = "improvement"
	ReleasePRTypeTechnical   ReleasePRType = "technical"
	ReleasePRTypeUnknown     ReleasePRType = "unknown"
)

type ReleasePR struct {
	Number             int           `json:"number"`
	Title              string        `json:"title"`
	URL                string        `json:"url"`
	AuthorLogin        string        `json:"author_login"`
	AuthorAvatarURL    string        `json:"author_avatar_url"`
	MergedAt           *time.Time    `json:"merged_at,omitempty"`
	Labels             []string      `json:"labels"`
	TypeClassification ReleasePRType `json:"type_classification"`
}

type ReleaseSummary struct {
	TotalPRs int                   `json:"total_prs"`
	ByType   map[ReleasePRType]int `json:"by_type"`
}

type ReleaseDiffReport struct {
	Repository  RepositoryRef  `json:"repository"`
	FromTag     string         `json:"from_tag"`
	ToTag       string         `json:"to_tag"`
	GeneratedAt time.Time      `json:"generated_at"`
	Summary     ReleaseSummary `json:"summary"`
	PRs         []ReleasePR    `json:"prs"`
}
