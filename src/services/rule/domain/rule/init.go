package rule

var (
	SvrRule  *ServiceRule
	RepoRule *RepositoryRule
)

func init() {
	SvrRule = &ServiceRule{}
	RepoRule = &RepositoryRule{}
}
