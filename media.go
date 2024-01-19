package ui

import (
	"bytes"
	"image"
	"io"
	"io/fs"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/gen2brain/mpeg"
)

const DefaultAudioSampleRate = 44100

func audioSampleRate() int {
	if audioContext != nil {
		return audioContext.SampleRate()
	}
	return DefaultAudioSampleRate
}

type Samples struct {
	Time  float64
	Bytes []byte
}

type MediaSamples struct {
	*Samples
	Media Media
}

type MediaImage struct {
	*image.RGBA
	Media Media
}

type Media interface {
	HasVideo() bool
	HasAudio() bool
	VideoSize() (int, int)
	SampleRate() int
	Duration() time.Duration
	Time() time.Duration
	Seek(time.Duration) bool
	Decode() chan<- (time.Duration)
	Rewind()
	AudioReader() (io.ReadSeeker, int64)
	Done() chan (bool)
	Video() <-chan (*MediaImage)
	Audio() <-chan (*MediaSamples)
}

type MPEG1MediaAudioReader struct {
	*mpeg.SamplesReader
	Ready chan (struct{})
}

func (r *MPEG1MediaAudioReader) Read(p []byte) (n int, err error) {
	<-r.Ready
	return r.SamplesReader.Read(p)
}

func (r *MPEG1MediaAudioReader) Seek(offset int64, whence int) (int64, error) {
	return r.SamplesReader.Seek(offset, whence)
}

type MPEG1Media struct {
	*mpeg.MPEG
	videoCallback func(Media, *image.RGBA)
	video         chan (*MediaImage)
	audio         chan (*MediaSamples)
	decode        chan (time.Duration)
	reader        *MPEG1MediaAudioReader
}

func (m *MPEG1Media) HasVideo() bool {
	return m.MPEG.NumVideoStreams() > 0
}

func (m *MPEG1Media) HasAudio() bool {
	return m.MPEG.NumAudioStreams() > 0
}

func (m *MPEG1Media) VideoSize() (int, int) {
	width := m.MPEG.Width()
	height := m.MPEG.Height()
	return width, height
}

func (m *MPEG1Media) SampleRate() int {
	return m.MPEG.Samplerate()
}

func (m *MPEG1Media) Done() chan (bool) {
	return m.MPEG.Done()
}

func (m *MPEG1Media) Duration() time.Duration {
	return time.Duration((float64(mpeg.SamplesPerFrame) / float64(m.MPEG.Samplerate())) * float64(time.Second))
}

func (m *MPEG1Media) Time() time.Duration {
	return m.MPEG.Time()
}

func (m *MPEG1Media) Seek(to time.Duration) bool {
	return m.MPEG.Seek(to, false)
}

func (m *MPEG1Media) Decode() chan<- (time.Duration) {
	return m.decode
}

func (m *MPEG1Media) Rewind() {
	m.MPEG.Rewind()
}

func (m *MPEG1Media) Video() <-chan (*MediaImage) {
	return m.video
}

func (m *MPEG1Media) Audio() <-chan (*MediaSamples) {
	return m.audio
}

func (m *MPEG1Media) AudioReader() (io.ReadSeeker, int64) {
	// The mpeg samples decoder hasa fixed size for audio fragments.
	ab := m.MPEG.Audio().Reader().(*mpeg.SamplesReader)
	return ab, mpeg.SamplesPerFrame
}

func (m *MPEG1Media) VideoCallback(_ *mpeg.MPEG, frame *mpeg.Frame) {
	mi := &MediaImage{Media: m}
	if frame != nil {
		mi.RGBA = frame.RGBA()
	}
	go func() { m.video <- mi }()
}

func (m *MPEG1Media) AudioCallback(_ *mpeg.MPEG, msamples *mpeg.Samples) {
	var samples *Samples
	if msamples != nil {
		samples = &Samples{}
		samples.Time = msamples.Time
		samples.Bytes = msamples.Bytes()
		ms := &MediaSamples{Samples: samples, Media: m}
		go func() { m.audio <- ms }()
		go func() { m.reader.Ready <- struct{}{} }()
	}
}

func (m *MPEG1Media) SetVideoCallback(callback func(Media, *image.RGBA)) {
	if m.HasVideo() {
		m.MPEG.SetVideoCallback(func(_ *mpeg.MPEG, frame *mpeg.Frame) {
			var img *image.RGBA
			if frame != nil {
				img = frame.RGBA()
			}
			callback(m, img)
		})
	}
}

func (m *MPEG1Media) SetAudioCallback(callback func(Media, *Samples)) {
	if m.HasAudio() {
		m.MPEG.SetAudioCallback(func(_ *mpeg.MPEG, msamples *mpeg.Samples) {
			var samples *Samples
			if msamples != nil {
				samples = &Samples{}
				samples.Time = msamples.Time
				samples.Bytes = msamples.Bytes()
			}
			callback(m, samples)
		})
	}
}

func NewMediaFromFile(arg string) (Media, error) {
	r, err := os.Open(arg)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return NewMediaFromReader(r)
}

func NewMediaFromFileSystem(fs fs.FS, arg string) (Media, error) {
	r, err := fs.Open(arg)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return NewMediaFromReader(r)
}

func NewMediaFromReader(rd io.Reader) (Media, error) {
	mpg, err := mpeg.New(rd)
	if err != nil {
		return nil, err
	}

	hasVideo := mpg.NumVideoStreams() > 0
	hasAudio := mpg.NumAudioStreams() > 0
	mpg.SetVideoEnabled(hasVideo)
	mpg.SetAudioEnabled(hasAudio)
	if hasAudio {
		samplerate := mpg.Samplerate()
		mpg.SetAudioFormat(mpeg.AudioS16)
		duration := time.Duration((float64(mpeg.SamplesPerFrame) / float64(samplerate)) * float64(time.Second))
		mpg.SetAudioLeadTime(duration)
	}

	res := &MPEG1Media{
		MPEG:   mpg,
		video:  make(chan (*MediaImage)),
		audio:  make(chan (*MediaSamples)),
		decode: make(chan (time.Duration)),
	}

	if res.HasVideo() {
		res.MPEG.SetVideoCallback(res.VideoCallback)
	}

	if res.HasAudio() {
		res.reader = &MPEG1MediaAudioReader{Ready: make(chan struct{})}
		res.MPEG.SetAudioCallback(res.AudioCallback)
	}

	go func() {
		for {
			dur := <-res.decode
			res.MPEG.Decode(dur)
		}
	}()

	return res, nil
}

type MediaPlayer struct {
	BasicWidget
	Media

	title      string
	pause      bool
	active     bool
	borderless bool
	seekTo     time.Duration

	img    *ebiten.Image
	player *audio.Player
	audio  *bytes.Reader
}

var audioContext *audio.Context

func SetAudioContext(ac *audio.Context) {
	audioContext = ac
}

func (p *MediaPlayer) AudioCallback(m Media, samples *Samples) {

}

func (p *MediaPlayer) getVideo() {
	if !p.Media.HasVideo() {
		return
	}
	select {
	case image := <-p.Media.Video():
		if image.RGBA != nil {
			p.img.WritePixels(image.RGBA.Pix)
		}
	default:
		return
	}
}

func (p *MediaPlayer) getAudio() {
	if !p.Media.HasAudio() {
		return
	}
	select {
	case samples := <-p.Media.Audio():
		p.audio.Reset(samples.Bytes)

		if p.player != nil {
			if p.pause && p.player.IsPlaying() {
				p.player.Pause()
			}
			if !p.pause && !p.player.IsPlaying() {
				p.player.Play()
			}
		}
	default:
		return
	}
}

func newMediaPlayer(m Media) (*MediaPlayer, error) {
	a := &MediaPlayer{}
	a.Media = m
	a.seekTo = -1

	a.width, a.height = a.Media.VideoSize()

	hasAudio := a.Media.HasAudio()
	a.img = ebiten.NewImage(a.width, a.height)

	if hasAudio {
		a.audio = bytes.NewReader(make([]byte, mpeg.SamplesPerFrame))
		samplerate := a.Media.SampleRate()
		if audioContext == nil {
			audioContext = audio.NewContext(DefaultAudioSampleRate)
		}
		var err error

		ar, as := a.Media.AudioReader()
		// as := int64(mpeg.SamplesPerFrame)
		resampled := audio.Resample(ar, as, samplerate, audioSampleRate())

		a.player, err = audioContext.NewPlayer(resampled)
		if err != nil {
			return nil, err
		}

		a.player.SetBufferSize(a.Media.Duration())
	}

	return a, nil
}

func (p *MediaPlayer) HandleWidget(ev Event) {
	if _, ok := ev.(*AwayEvent); ok {
		p.active = false
		p.pause = true
		if p.pause && p.player != nil {
			p.player.Pause()
		}
		return
	}

	if _, ok := ev.(*MouseClickEvent); ok {
		if p.active {
			p.pause = !p.pause
		}
		p.active = true
	}
	if p.active {
		if ke, ok := ev.(*KeyPressEvent); ok {
			p.HandleKeyPress(ke)
		}
		if ue, ok := ev.(*UpdateEvent); ok {
			p.HandleUpdate(ue)
		}
	}
}

func (p *MediaPlayer) HandleKeyPress(kp *KeyPressEvent) {
	// Allow copy pasting of images.
	switch kp.Key {
	case KeySpace, KeyP:
		p.pause = !p.pause
	case KeyArrowRight:
		p.seekTo = p.Media.Time() + 3
	case KeyArrowLeft:
		p.seekTo = p.Media.Time() - 3
	}
}

func (a *MediaPlayer) HandleUpdate(ue *UpdateEvent) error {
	if a.seekTo >= 0 {
		a.Media.Seek(a.seekTo)
		a.seekTo = -1
	} else if a.pause && a.player != nil {
		a.player.Pause()
	} else {
		a.Media.Decode() <- ue.Duration
		a.getVideo()
		a.getAudio()
	}

	select {
	case done := <-a.Media.Done():
		if done {
			a.pause = true
			if a.player != nil {
				a.player.Pause()
				a.Media.Rewind()
			}
		}
	default:
	}

	return nil
}

func NewMediaPlayer(title string, media Media) *MediaPlayer {
	mp, err := newMediaPlayer(media)
	mp.SetTitle(title)
	if err != nil {
		mp.SetTitle(title + ":" + err.Error())
	}
	mp.SetStyle(theme.Media)
	return mp
}

func (w *MediaPlayer) SetTitle(title string) {
	w.title = title
}

func (w *MediaPlayer) SetMedia(m Media) {
	w.SetTitle("SetMedia not supported yet.")
}

func (w *MediaPlayer) Title() string {
	return w.title
}

func (g *MediaPlayer) LayoutWidget(width, height int) {
	g.width, g.height = 0, 0
	if g.img != nil {
		g.width, g.height = g.Media.VideoSize()
	}

	txt := g.Title()
	tw, th := oneLineTextSize(g.Style().Font.Face, txt)
	minh := g.Style().Font.Face.Metrics().Height.Round()
	if tw > g.width {
		g.width = tw
	}
	g.height += th
	if g.height < minh {
		g.height = minh
	}

	margin := g.Style().Margin.Int()
	g.GrowToStyleSize()
	g.width += margin * 2
	g.height += margin * 2
	g.ClipTo(width, height)
}

func (w *MediaPlayer) Destroy() {
	// first hide ourselves
	w.Hide()
	// now destroy the media
	if w.img != nil {
		w.img.Dispose()
		w.img = nil
	}
	if w.player != nil {
		w.player.Close()
		w.player = nil
	}
	// and finally free ourselves
}

func (w *MediaPlayer) Borderless() bool {
	return w.borderless
}

func (w *MediaPlayer) SetBorderless(borderless bool) {
	w.borderless = borderless
}

func (w *MediaPlayer) DrawWidget(screen *Graphic) {
	dx, dy := w.WidgetAbsolute()
	ww, wh := w.WidgetSize()

	var (
		fillColor = w.Style().Fill.Color.RGBA()
		textColor = w.Style().Color.RGBA()
		textFace  = w.Style().Font.Face
		mp        = w.Style().Margin.Int()
	)

	if !w.borderless {
		FillFrameStyle(screen, dx, dy, ww, wh, w.Style())
	}
	if w.title != "" {
		FillRect(screen, dx+mp, dy, ww-2*mp, wh, fillColor)
		TextDrawOffset(screen, w.title, textFace, dx+mp, dy, textColor)
	}

	if w.img != nil {
		iw, ih := w.img.Size()
		_, th := oneLineTextSize(w.Style().Font.Face, w.title)
		ix, iy := dx+mp, dy+mp+th
		sx, sy := float64(ww)/float64(iw), float64(wh)/float64(ih)
		DrawGraphicAtScale(screen, w.img, ix, iy, sx, sy)
	}

	w.DrawDebug(screen, "PIC %d %d", dx, dy)
}
