package magnetConvert

import (
	"fmt"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"
	"io"
)

//
// FileToMagnet
//  @Description: 通过文件获取磁链
//  @param file
//  @return string
//  @return error
//
func FileToMagnet(file string) (string, error) {
	mi, err := metainfo.LoadFromFile(file)
	if err != nil {
		return "", fmt.Errorf("Cannot read the metainfo from file: %s. %v", file, err)
	}

	info, err := mi.UnmarshalInfo()
	if err != nil {
		return "", fmt.Errorf("Cannot unmarshal the metainfo from file: %s. %v", file, err)
	}
	hs := mi.HashInfoBytes()

	if info.Name == "" {
		return "", nil
	}

	return mi.Magnet(&hs, &info).String(), nil
}

//
// IO2Magnet
//  @Description: 通过io.Reader获取磁链
//  @param r
//  @return string
//  @return error
//
func IO2Magnet(r io.Reader) (string, error) {
	mi, err := metainfo.Load(r)
	if err != nil {
		return "", errors.New("读取磁链meta信息错误")
	}

	info, err := mi.UnmarshalInfo()
	if err != nil {
		return "", errors.New("磁链解析错误")
	}
	hs := mi.HashInfoBytes()

	if info.Name == "" {
		return "", nil
	}

	return mi.Magnet(&hs, &info).String(), nil
}
