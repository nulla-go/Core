package main

import (
	"github.com/strengine/Core/av"
	"github.com/strengine/Core/av/avutil"
	"github.com/strengine/Core/av/transcode"
	"github.com/strengine/Core/cgo/ffmpeg"
	"github.com/strengine/Core/format"
	"github.com/strengine/Core/format/rtmp"
)

// need ffmpeg with libspeex and libfdkaac installed
//
// open http://www.wowza.com/resources/4.4.1/examples/WebcamRecording/FlashRTMPPlayer11/player.html
// click connect and reCored
// input camera H264/SPEEX will converted H264/AAC and saved in out.ts

func init() {
	format.RegisterAll()
}

func main() {
	server := &rtmp.Server{}

	server.HandlePublish = func(conn *rtmp.Conn) {
		file, _ := avutil.Create("out.ts")

		findcodec := func(stream av.AudioCodecData, i int) (need bool, dec av.AudioDecoder, enc av.AudioEncoder, err error) {
			need = true
			dec, _ = ffmpeg.NewAudioDecoder(stream)
			enc, _ = ffmpeg.NewAudioEncoderByName("libfdk_aac")
			enc.SetSampleRate(48000)
			enc.SetChannelLayout(av.CH_STEREO)
			return
		}

		trans := &transcode.Demuxer{
			Options: transcode.Options{
				FindAudioDecoderEncoder: findcodec,
			},
			Demuxer: conn,
		}

		avutil.CopyFile(file, trans)
	}

	server.ListenAndServe()
}
