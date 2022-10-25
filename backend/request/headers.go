package request

var genericHeaders = map[string]string{
	"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"accept-language": "en-US,en;q=0.9",
	"cache-control":   "no-cache",
	"pragma":          "no-cache",
	"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",

	// allow 3rd party sites to block atus by filtering out requests with this header
	"software": "ATUS",
}
