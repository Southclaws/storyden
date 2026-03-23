package settings

//go:generate go run github.com/Southclaws/enumerator

type clientIPModeEnum string

const (
	clientIPModeRemoteAddr        clientIPModeEnum = "remote_addr"
	clientIPModeSingleHeader      clientIPModeEnum = "single_header"
	clientIPModeXFFTrustedProxies clientIPModeEnum = "xff_trusted_proxies"
)
