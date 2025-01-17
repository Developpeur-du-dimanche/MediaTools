package medias

type Video struct {
	StreamIndex int    `json:"index"`
	CodecName   string `json:"codec_name"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

type Audio struct {
	StreamIndex int    `json:"index"`
	CodecName   string `json:"codec_name"`
	Channels    int    `json:"channels"`
	Language    string `json:"language"`
}

type Subtitle struct {
	StreamIndex int    `json:"index"`
	CodecName   string `json:"codec_name"`
	Language    string `json:"language"`
}

type Ffprobe struct {
	path string
}

func NewFfprobe(path string) *Ffprobe {
	return &Ffprobe{
		path: path,
	}
}
