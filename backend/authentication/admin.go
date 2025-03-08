package authentication

func (a *apiKeyAuthenticator) IsAdmin(apiKey string) bool {
	return a.adminKey == apiKey
}
