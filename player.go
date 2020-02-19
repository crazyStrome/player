package player

import (
	"fmt"
	"github.com/hajimehoshi/oto"
	"os"
	"errors"
	"github.com/hajimehoshi/go-mp3"
)
type MusicInfo struct {
	name    string
	file    *os.File
	dec     *mp3.Decoder
	play    *oto.Player
	chPause chan bool
	data    []byte
	pos     int64
}

//构造一个MusicInfo
func NewMusicInfo(name string) (*MusicInfo, error) {
	music := new(MusicInfo)
	music.name = name
	f, err := os.Open(name)
	music.file = f
	if err != nil {
		return nil, err
	}
	//defer f.Close()

	music.dec, err = mp3.NewDecoder(music.file)
	if err != nil {
		return nil, err
	}

	//d.Seek(d.Length()/4, io.SeekEnd)
	music.play, err = oto.NewPlayer(music.dec.SampleRate(), 2, 2, 8192)
	if err != nil {
		return nil, err
	}
	music.chPause = make(chan bool)

	music.data = make([]byte, 51200)
	return music, nil
}

//获取音乐文件的长度
func (music *MusicInfo) Len() (int64, error) {
	if music.dec != nil {
		return music.dec.Length(), nil
	}
	return int64(0), errors.New("no decoder")
}

//获取音乐读取的位置
func (music *MusicInfo) CurrPos() (int64, error) {
	if music == nil {
		return 0, errors.New("no musicInfo")
	}
	return music.pos, nil
}

//获取音乐播放的进度:80%
func (music *MusicInfo) Progress() (int64, error) {
	n, err := music.CurrPos()
	if err != nil {
		return 0, err
	}
	total, err := music.Len()
	if err != nil {
		return 0, err
	}
	return n * 100 / total, nil
}

//播放音乐
func (music *MusicInfo) MusicPlay() error {
	var flag bool = false
	for {
		select {
		case flag = <-music.chPause:
			fmt.Println("Paused?:", flag)
		default:
			if flag {
				continue
			} else {
				n, err := music.dec.Read(music.data)
				if err != nil {
					return err
				}
				tmp := music.data[:n]
				n, err = music.play.Write(tmp)
				if err != nil {
					return err
				}
				music.pos += int64(n)
			}
		}

	}
}

//暂停
func (music *MusicInfo) Pause() {
	music.chPause <- true
}
//继续开始
func (music *MusicInfo) Restart() {
	music.chPause <- false
}

