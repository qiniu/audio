module github.com/qiniu/audio

go 1.14

require (
	github.com/hajimehoshi/go-mp3 v0.2.1
	github.com/hajimehoshi/oto v0.5.4
	github.com/qiniu/x v1.8.4
)

replace (
	github.com/hajimehoshi/oto v0.3.4 => github.com/qiniu/oto v0.5.4-fixed
	github.com/hajimehoshi/oto v0.5.4 => github.com/qiniu/oto v0.5.4-fixed
)
