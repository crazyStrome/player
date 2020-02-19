#  Windows下的音乐播放器

使用函数如下

```
func NewMusicInfo(name string) (*MusicInfo, error)
func (music *MusicInfo) CurrPos() (int64, error)
func (music *MusicInfo) Len() (int64, error)
func (music *MusicInfo) MusicPlay() error
func (music *MusicInfo) Pause()
func (music *MusicInfo) Progress() (int64, error)
func (music *MusicInfo) Restart()
```