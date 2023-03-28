package log

type Pairs map[string]any

const (
	pairAppId      = "app.id"
	pairAppVersion = "app.version"
)

func AppInfo(appId string, version string) Pairs {
	return Pairs{pairAppId: appId, pairAppVersion: version}
}
func AppInfoWithFields(appId string, version string, fields map[string]any) Pairs {
	p := Pairs{pairAppId: appId, pairAppVersion: version}
	if fields != nil {
		for k, v := range fields {
			p[k] = v
		}
	}
	return p
}
