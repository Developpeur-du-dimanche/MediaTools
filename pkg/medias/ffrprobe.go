package medias

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/Developpeur-du-dimanche/MediaTools/pkg/logger"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/proc"
)

// ErrTagNotFound is a sentinel error used when a queried tag does not exist
var ErrTagNotFound = errors.New("tag not found")

// tags is the map of tag names to values
type tags map[string]interface{}

// GetInt returns a tag value as int64 and an error if one occurred.
// ErrTagNotFound will be returned if the key can't be found, ParseError if
// a parsing error occurs.
func (t tags) GetInt(tag string) (int64, error) {
	v, found := t[tag]
	if !found || v == nil {
		return 0, ErrTagNotFound
	}

	switch v := v.(type) {
	case string:
		return valToInt64(v)
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	}

	str := fmt.Sprintf("%v", v)
	return valToInt64(str)
}

func valToInt64(str string) (int64, error) {
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("int64 parsing error (%v): %w", str, err)
	}
	return val, nil
}

// GetString returns a tag value as string and an error if one occurred.
// ErrTagNotFound will be returned if the key can't be found
func (t tags) GetString(tag string) (string, error) {
	v, found := t[tag]
	if !found || v == nil {
		return "", ErrTagNotFound
	}
	return valToString(v), nil
}

func valToString(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int64:
		return strconv.FormatInt(v, 10)
	}

	return fmt.Sprintf("%v", v)
}

// GetFloat returns a tag value as float64 and an error if one occurred.
// ErrTagNotFound will be returned if the key can't be found.
func (t tags) GetFloat(tag string) (float64, error) {
	v, found := t[tag]
	if !found || v == nil {
		return 0, ErrTagNotFound
	}

	switch v := v.(type) {
	case string:
		return valToFloat64(v)
	case float64:
		return v, nil
	case int64:
		return float64(v), nil
	}

	str := fmt.Sprintf("%v", v)
	return valToFloat64(str)
}

func valToFloat64(str string) (float64, error) {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, fmt.Errorf("float64 parsing error (%v): %w", str, err)
	}
	return val, nil
}

// formatTags is a json data structure to represent format tags
// Deprecated, use the tags of TagList instead
type formatTags struct {
	MajorBrand       string `json:"major_brand"`
	MinorVersion     string `json:"minor_version"`
	CompatibleBrands string `json:"compatible_brands"`
	CreationTime     string `json:"creation_time"`
}

func (f *formatTags) setFrom(tags tags) {
	f.MajorBrand, _ = tags.GetString("major_brand")
	f.MinorVersion, _ = tags.GetString("minor_version")
	f.CompatibleBrands, _ = tags.GetString("compatible_brands")
	f.CreationTime, _ = tags.GetString("creation_time")
}

// streamTags is a json data structure to represent stream tags
// Deprecated, use the tags of TagList instead
type streamTags struct {
	Rotate       int    `json:"rotate,string,omitempty"`
	CreationTime string `json:"creation_time,omitempty"`
	Language     string `json:"language,omitempty"`
	Title        string `json:"title,omitempty"`
	Encoder      string `json:"encoder,omitempty"`
	Location     string `json:"location,omitempty"`
}

func (s *streamTags) setFrom(tags tags) {
	rotate, _ := tags.GetInt("rotate")
	s.Rotate = int(rotate)

	s.CreationTime, _ = tags.GetString("creation_time")
	s.Language, _ = tags.GetString("language")
	s.Title, _ = tags.GetString("title")
	s.Encoder, _ = tags.GetString("encoder")
	s.Location, _ = tags.GetString("location")
}

// streamType represents a media stream type like video, audio, subtitles, etc
type streamType string

const (
	// StreamAny means any type of stream
	StreamAny streamType = ""
	// StreamVideo is a video stream
	StreamVideo streamType = "video"
	// StreamAudio is an audio stream
	StreamAudio streamType = "audio"
	// StreamSubtitle is a subtitle stream
	StreamSubtitle streamType = "subtitle"
	// StreamData is a data stream
	StreamData streamType = "data"
	// StreamAttachment is an attachment stream
	StreamAttachment streamType = "attachment"
)

// probeData is the root json data structure returned by an ffprobe.
type probeData struct {
	Streams []*stream `json:"streams"`
	Format  *format   `json:"format"`
}

// format is a json data structure to represent formats
type format struct {
	Filename         string      `json:"filename"`
	NBStreams        int         `json:"nb_streams"`
	NBPrograms       int         `json:"nb_programs"`
	StartTimeSeconds float64     `json:"start_time,string"`
	DurationSeconds  float64     `json:"duration,string"`
	Size             string      `json:"size"`
	BitRate          string      `json:"bit_rate"`
	ProbeScore       int         `json:"probe_score"`
	TagList          tags        `json:"tags"`
	tags             *formatTags `json:"-"` // Deprecated: Use TagList instead
}

// Stream is a json data structure to represent streams.
// A stream can be a video, audio, subtitle, etc type of stream.
type stream struct {
	Index              int               `json:"index"`
	ID                 string            `json:"id"`
	CodecName          string            `json:"codec_name"`
	CodecLongName      string            `json:"codec_long_name"`
	CodecType          string            `json:"codec_type"`
	CodecTimeBase      string            `json:"codec_time_base"`
	CodecTagString     string            `json:"codec_tag_string"`
	CodecTag           string            `json:"codec_tag"`
	RFrameRate         string            `json:"r_frame_rate"`
	AvgFrameRate       string            `json:"avg_frame_rate"`
	TimeBase           string            `json:"time_base"`
	StartPts           int               `json:"start_pts"`
	StartTime          string            `json:"start_time"`
	DurationTs         uint64            `json:"duration_ts"`
	Duration           string            `json:"duration"`
	BitRate            string            `json:"bit_rate"`
	BitsPerRawSample   string            `json:"bits_per_raw_sample"`
	NbFrames           string            `json:"nb_frames"`
	Disposition        StreamDisposition `json:"disposition,omitempty"`
	TagList            tags              `json:"tags"`
	tags               streamTags        `json:"-"` // Deprecated: Use TagList instead
	FieldOrder         string            `json:"field_order,omitempty"`
	Profile            string            `json:"profile,omitempty"`
	Width              int               `json:"width"`
	Height             int               `json:"height"`
	HasBFrames         int               `json:"has_b_frames,omitempty"`
	SampleAspectRatio  string            `json:"sample_aspect_ratio,omitempty"`
	DisplayAspectRatio string            `json:"display_aspect_ratio,omitempty"`
	PixFmt             string            `json:"pix_fmt,omitempty"`
	Level              int               `json:"level,omitempty"`
	ColorRange         string            `json:"color_range,omitempty"`
	ColorSpace         string            `json:"color_space,omitempty"`
	SampleFmt          string            `json:"sample_fmt,omitempty"`
	SampleRate         string            `json:"sample_rate,omitempty"`
	Channels           int               `json:"channels,omitempty"`
	ChannelLayout      string            `json:"channel_layout,omitempty"`
	BitsPerSample      int               `json:"bits_per_sample,omitempty"`
}

// StreamDisposition is a json data structure to represent stream dispositions
type StreamDisposition struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karaoke         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_impaired"`
	CleanEffects    int `json:"clean_effects"`
	AttachedPic     int `json:"attached_pic"`
}

// StartTime returns the start time of the media file as a time.Duration
func (f *format) StartTime() (duration time.Duration) {
	return time.Duration(f.StartTimeSeconds * float64(time.Second))
}

// Duration returns the duration of the media file as a time.Duration
func (f *format) Duration() (duration time.Duration) {
	return time.Duration(f.DurationSeconds * float64(time.Second))
}

// streamType returns all streams which are of the given type
func (p *probeData) streamType(streamType streamType) (streams []stream) {
	for _, s := range p.Streams {
		if s == nil {
			continue
		}
		switch streamType {
		case StreamAny:
			streams = append(streams, *s)
		default:
			if s.CodecType == string(streamType) {
				streams = append(streams, *s)
			}
		}
	}
	return streams
}

// FirstVideoStream returns the first video stream found
func (p *probeData) FirstVideoStream() *stream {
	return p.firstStream(StreamVideo)
}

// FirstAudioStream returns the first audio stream found
func (p *probeData) FirstAudioStream() *stream {
	return p.firstStream(StreamAudio)
}

// FirstSubtitleStream returns the first subtitle stream found
func (p *probeData) FirstSubtitleStream() *stream {
	return p.firstStream(StreamSubtitle)
}

// FirstDataStream returns the first data stream found
func (p *probeData) FirstDataStream() *stream {
	return p.firstStream(StreamData)
}

// FirstAttachmentStream returns the first attachment stream found
func (p *probeData) FirstAttachmentStream() *stream {
	return p.firstStream(StreamAttachment)
}

func (p *probeData) firstStream(streamType streamType) *stream {
	for _, s := range p.Streams {
		if s == nil {
			continue
		}
		if s.CodecType == string(streamType) {
			return s
		}
	}
	return nil
}

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

type FfprobeData struct {
	Filename        string        `json:"filename"`
	DurationSeconds time.Duration `json:"duration,string"`
	Size            string        `json:"size"`
	BitRate         string        `json:"bit_rate"`
}

type FfprobeResult struct {
	Format    FfprobeData `json:"format"`
	Videos    []Video     `json:"video"`
	Audios    []Audio     `json:"audio"`
	Subtitles []Subtitle  `json:"subtitle"`
}

type FfprobeOptions struct {
	option string
	value  string
}

var (
	FFPROBE_LOGLEVEL_QUIET   FfprobeOptions = FfprobeOptions{"-loglevel", "quiet"}
	FFPROBE_LOGLEVEL_PANIC   FfprobeOptions = FfprobeOptions{"-loglevel", "panic"}
	FFPROBE_LOGLEVEL_FATAL   FfprobeOptions = FfprobeOptions{"-loglevel", "fatal"}
	FFPROBE_LOGLEVEL_ERROR   FfprobeOptions = FfprobeOptions{"-loglevel", "error"}
	FFPROBE_LOGLEVEL_WARNING FfprobeOptions = FfprobeOptions{"-loglevel", "warning"}
	FFPROBE_LOGLEVEL_INFO    FfprobeOptions = FfprobeOptions{"-loglevel", "info"}
	FFPROBE_LOGLEVEL_VERBOSE FfprobeOptions = FfprobeOptions{"-loglevel", "verbose"}
	FFPROBE_LOGLEVEL_DEBUG   FfprobeOptions = FfprobeOptions{"-loglevel", "debug"}

	SHOW_FORMAT  FfprobeOptions = FfprobeOptions{"-show_format", ""}
	SHOW_STREAMS FfprobeOptions = FfprobeOptions{"-show_streams", ""}

	PRINT_FORMAT_JSON FfprobeOptions = FfprobeOptions{"-print_format", "json"}
	EXPERIMENTAL      FfprobeOptions = FfprobeOptions{"-strict", "experimental"}
)

type Ffprobe struct {
	binary  string
	path    string
	options []FfprobeOptions
}

func NewFfprobe(path string, options ...FfprobeOptions) *Ffprobe {
	return &Ffprobe{
		path:    path,
		options: options,
	}
}

func (f *Ffprobe) hasJsonOption() bool {
	for _, option := range f.options {
		if option == PRINT_FORMAT_JSON {
			return true
		}
	}
	return false
}

func (f *Ffprobe) Probe(ctx context.Context) (*FfprobeResult, error) {

	args := []string{}

	if !f.hasJsonOption() {
		f.options = append(f.options, PRINT_FORMAT_JSON)
	}

	for _, option := range f.options {
		args = append(args, option.option)
		if option.value != "" {
			args = append(args, option.value)
		}
	}

	args = append(args, f.path)

	if f.binary == "" {
		// get binary from PATH
		binary, err := exec.LookPath("ffprobe")

		if err != nil {
			logger.Errorf("ffprobe binary not found in PATH: %v", err)
			return nil, fmt.Errorf("ffprobe not found: %w", err)
		}

		f.binary = binary
	}

	cmd := exec.CommandContext(ctx, "ffprobe", args...)
	cmd.SysProcAttr = proc.ProcAttributes()

	data, err := runCmd(cmd)

	if err != nil {
		logger.Errorf("ffprobe command failed for %s: %v", f.path, err)
		return nil, err
	}

	result := &FfprobeResult{
		Format: FfprobeData{
			Filename:        data.Format.Filename,
			DurationSeconds: data.Format.Duration(),
			Size:            data.Format.Size,
			BitRate:         data.Format.BitRate,
		},
		Videos:    make([]Video, len(data.streamType(StreamVideo))),
		Audios:    make([]Audio, len(data.streamType(StreamAudio))),
		Subtitles: make([]Subtitle, len(data.streamType(StreamSubtitle))),
	}

	for i, stream := range data.streamType(StreamVideo) {
		result.Videos[i] = Video{
			StreamIndex: stream.Index,
			CodecName:   stream.CodecName,
			Width:       stream.Width,
			Height:      stream.Height,
		}
	}

	for i, stream := range data.streamType(StreamAudio) {
		result.Audios[i] = Audio{
			StreamIndex: stream.Index,
			CodecName:   stream.CodecName,
			Channels:    stream.Channels,
			Language:    stream.tags.Language,
		}
	}

	for i, stream := range data.streamType(StreamSubtitle) {
		result.Subtitles[i] = Subtitle{
			StreamIndex: stream.Index,
			CodecName:   stream.CodecName,
			Language:    stream.tags.Language,
		}
	}

	return result, nil

}

func runCmd(cmd *exec.Cmd) (*probeData, error) {
	var outputBuf bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stdout = &outputBuf
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("ffprobe error: %w", err)
	}

	if stdErr.Len() > 0 {
		return nil, fmt.Errorf("ffprobe error: %s", stdErr.String())
	}

	var data = &probeData{}

	err = json.Unmarshal(outputBuf.Bytes(), data)
	if err != nil {
		return nil, fmt.Errorf("ffprobe error: %w", err)
	}

	if data.Format == nil {
		return nil, errors.New("ffprobe error: no format found")
	}

	if len(data.Format.TagList) > 0 {
		data.Format.tags = &formatTags{}
		data.Format.tags.setFrom(data.Format.TagList)
	}

	for _, stream := range data.Streams {
		stream.tags.setFrom(stream.TagList)
	}

	return data, nil

}
