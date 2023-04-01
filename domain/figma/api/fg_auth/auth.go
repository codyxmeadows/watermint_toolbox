package fg_auth

import (
	"github.com/watermint/toolbox/essentials/api/api_auth"
	"github.com/watermint/toolbox/infra/app"
)

const (
	ScopeFileRead = "file_read"
)

var (
	Figma = api_auth.OAuthAppData{
		AppKeyName:       app.ServiceFigma,
		EndpointAuthUrl:  "https://www.figma.com/oauth",
		EndpointTokenUrl: "https://www.figma.com/api/oauth/token",
		EndpointStyle:    api_auth.AuthStyleAutoDetect,
		UsePKCE:          false,
		RedirectUrl:      "",
	}
)
