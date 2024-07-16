package ebitenanimator

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	defaultImage *ebiten.Image
	frames       []*Frame
}

func NewAnimation() *Animation {
	return &Animation{}
}

func NewAnimationFixed(delta time.Duration, images ...*ebiten.Image) *Animation {
	if len(images) == 0 {
		return nil
	}

	frames := make([]*Frame, 0, len(images))
	for _, img := range images {
		frames = append(frames, &Frame{
			Duration: delta,
			Image:    img,
		})
	}

	return &Animation{
		defaultImage: images[0],
		frames:       frames,
	}
}

func (animation *Animation) SetDefault(img *ebiten.Image) *Animation {
	animation.defaultImage = img

	return animation
}

func (animation *Animation) Add(duration time.Duration, img *ebiten.Image, hooks ...func()) *Animation {
	animation.frames = append(animation.frames, &Frame{
		Duration: duration,
		Image:    img,
		Hooks:    hooks,
	})
	return animation
}

type Frame struct {
	Duration time.Duration
	Image    *ebiten.Image
	Hooks    []func()
}

type Animator struct {
	frameIndex  int
	lastUpdated time.Time
	lastImage   *ebiten.Image
	track       *Animation
}

func New(track *Animation) *Animator {
	return &Animator{
		frameIndex: -1,
		track:      track,
	}
}

func (animator *Animator) Update() {
	if animator.track == nil ||
		animator.frameIndex < 0 ||
		animator.frameIndex >= len(animator.track.frames) {
		return
	}

	currentFrame := animator.track.frames[animator.frameIndex]
	animator.lastImage = currentFrame.Image

	dt := currentFrame.Duration
	if time.Since(animator.lastUpdated) < dt {
		return
	}

	animator.lastUpdated = time.Now()
	animator.frameIndex = (animator.frameIndex + 1) % len(animator.track.frames)
	frame := animator.track.frames[animator.frameIndex]

	for _, runHook := range frame.Hooks {
		runHook()
	}
}

func (animator *Animator) Duration() time.Duration {
	if animator.track == nil {
		return 0
	}

	var duration time.Duration
	for _, frame := range animator.track.frames {
		duration += frame.Duration
	}

	return duration
}

func (animator *Animator) Reset(animation *Animation) {
	animator.lastUpdated = time.Time{}
	animator.track = animation
	animator.frameIndex = -1
}

func (animator *Animator) Stop() {
	animator.frameIndex = -1
}

func (animator *Animator) Restart() {
	animator.frameIndex = 0
}

func (animator *Animator) Draw(dst *ebiten.Image, opts *ebiten.DrawImageOptions) bool {
	img := animator.Image()
	if img != nil {
		dst.DrawImage(img, opts)
		return true
	}
	return false
}

func (animator *Animator) Image() *ebiten.Image {
	if animator.track == nil || len(animator.track.frames) == 0 {
		return nil
	}

	if animator.frameIndex < 0 {
		return animator.track.defaultImage
	}

	frame := animator.track.frames[animator.frameIndex]
	if frame.Image == nil {
		return animator.lastImage
	}

	return frame.Image
}
