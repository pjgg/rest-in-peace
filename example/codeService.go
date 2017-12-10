package example

// Your app service that use a third party Rest endpoint (above github service)
type codeService struct {
	githubConnector github
}

type codeRepository interface {
	getUserInfo() map[string]interface{}
}

func (codeService *codeService) getUserInfo() map[string]interface{} {
	return codeService.githubConnector.getCurrentUserInfo()
}
