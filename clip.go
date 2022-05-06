package logicclipping

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Clip struct {
	s3.Object
	s3Api     *AwsS3ClientAPI
	bucket    string
	LocalPath string
}

func (clip *Clip) GetData() ([]byte, error) {
	return clip.s3Api.GetObject(clip.bucket, *clip.Key)
}

func (clip *Clip) Download(dir string) error {
	file := *clip.Key
	_, file, _ = strings.Cut(file, "/") // key is usually "clips/<fileName>" so we cut the first part
	localPath := fmt.Sprintf("%s/%s", dir, file)

	// check if file exists
	_, err := os.Stat(localPath)
	if err == nil {
		clip.LocalPath = localPath
		return nil
	}

	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	data, err := clip.GetData()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(localPath, data, 0644)
	if err != nil {
		return err
	}
	clip.LocalPath = localPath
	return nil
}

func (clip *Clip) DownloadAs(dir, format string) error {
	err := clip.Download(dir)
	if err != nil {
		return err
	}

	err = clip.Transcode(format)
	if err != nil {
		return err
	}

	return nil
}

func (clip *Clip) Transcode(format string) error {
	if clip.LocalPath == "" {
		return errors.New("clip not downloaded")
	}
	if format == "ts" {
		return nil
	}
	input := clip.LocalPath
	output := fmt.Sprintf("%s.%s", input, format)

	// check if file exists
	_, err := os.Stat(output)
	if err == nil {
		return nil
	}

	err = transcode(input, output)
	if err != nil {
		return err
	}
	clip.LocalPath = output

	// // TODO: Cleanup ?
	// err = os.Remove(input)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (lc *LogicConnection) GetClips() ([]*Clip, error) {
	clips := []*Clip{}

	s3Clips, err := lc.s3.ListObjects(lc.bucketOutputName)
	if err != nil {
		return nil, err
	}
	for _, s3Clip := range s3Clips {
		if !strings.HasSuffix(*s3Clip.Key, ".ts") {
			continue
		}
		clips = append(clips, &Clip{
			Object: *s3Clip,
			s3Api:  lc.s3,
			bucket: lc.bucketOutputName,
		})
	}
	return clips, nil
}

func (lc *LogicConnection) GetClipByAssetName(assetName string) (*Clip, error) {
	allClips, err := lc.GetClips()
	if err != nil {
		return nil, err
	}
	for _, clip := range allClips {
		// checks if clip starts with assetName
		if strings.HasPrefix(*clip.Key, fmt.Sprintf("clips/%s", assetName)) {
			return clip, nil
		}
	}

	return nil, errors.New("clip not found")
}

func transcode(input, output string) error {
	cmd := execute.ExecTask{
		Command:     "ffmpeg",
		Args:        []string{"-i", input, "-c", "copy", "-y", output},
		StreamStdio: false,
	}

	res, err := cmd.Execute()
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return errors.New("non-zero exit code: " + res.Stderr)
	}

	return nil
}
