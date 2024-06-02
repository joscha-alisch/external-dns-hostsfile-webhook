package api

const (
	mediaTypeFormat        = "application/external.dns.webhook+json;"
	contentTypeHeader      = "Content-Type"
	contentTypePlaintext   = "text/plain"
	acceptHeader           = "Accept"
	varyHeader             = "Vary"
	supportedMediaVersions = "1"
	healthPath             = "/health"
	logFieldRequestPath    = "requestPath"
	logFieldRequestMethod  = "requestMethod"
	logFieldError          = "error"
)

var mediaTypeVersion1 = mediaTypeVersion("1")

type mediaType string

func mediaTypeVersion(v string) mediaType {
	return mediaType(mediaTypeFormat + "version=" + v)
}