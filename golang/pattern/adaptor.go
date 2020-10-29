package pattern

import "fmt"

type MusicPlayer interface {
	play(fileType string, fileName string)
}

type OldPlayer struct {
}

func (o *OldPlayer) playMp3(fileName string)  {
	fmt.Println(fileName)
}

func (o *OldPlayer) playWma(fileName string) {
	fmt.Println(fileName)
}

type PlayerAdaptor struct {
	oldPlayer OldPlayer
}

func (player *PlayerAdaptor) play(fileType string, fileName string) {
	switch fileType {
	case "mp3":
		player.oldPlayer.playMp3(fileName)
	case "wma":
		player.oldPlayer.playWma(fileName)
	default:

	}
}
