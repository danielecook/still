package validate

import (
	"os"

	"github.com/c2h5oh/datasize"
	"github.com/danielecook/still/src/utils"
	"github.com/gabriel-vasile/mimetype"
)

func fileExists(args ...interface{}) (interface{}, error) {
	// Checks for an element present in a set.
	if m, _ := isMissing(args[0]); m.(bool) {
		return (bool)(true), nil
	}
	if _, err := os.Stat(args[0].(string)); err == nil || os.IsExist(err) {
		return (bool)(true), nil
	}
	return (bool)(false), nil
}

func parseSize(s string) int64 {
	var v datasize.ByteSize
	err := v.UnmarshalText([]byte(s))
	utils.Check(err)
	return int64(v)
}

func fileMinSize(args ...interface{}) (interface{}, error) {
	// Checks for file greater than a given size
	if m, _ := isMissing(args[0]); m.(bool) {
		return (bool)(true), nil
	}
	if stat, err := os.Stat(args[0].(string)); err == nil || os.IsExist(err) {
		if stat.Size() >= parseSize(args[1].(string)) {
			return (bool)(true), nil
		}
	}
	return (bool)(false), nil
}

func fileMaxSize(args ...interface{}) (interface{}, error) {
	// Checks for file greater than a given size
	if m, _ := isMissing(args[0]); m.(bool) {
		return (bool)(true), nil
	}
	if stat, err := os.Stat(args[0].(string)); err == nil || os.IsExist(err) {
		if stat.Size() <= parseSize(args[1].(string)) {
			return (bool)(true), nil
		}
	}
	return (bool)(false), nil
}

func mimeTypeIs(args ...interface{}) (interface{}, error) {
	if m, _ := isMissing(args[0]); m.(bool) {
		return (bool)(true), nil
	}
	mime, err := mimetype.DetectFile(args[0].(string))
	utils.Check(err)
	return (bool)(mime.String() == args[1].(string)), nil
}

// has_extension

// md5 checksum

// sha1 checksum

// sha256 checksum

// is mimetype --> image

// Directory exists
