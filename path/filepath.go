package path

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// CheckDirWritable check whether the specified dir writable
func CheckDirWritable(dir string) (err error) {
	logDirFileInfo, statErr := os.Stat(dir)
	if statErr != nil && os.IsNotExist(statErr) {
		// try to mkdir if not exists
		mkdirErr := os.MkdirAll(dir, 0755)
		if mkdirErr != nil {
			err = mkdirErr
			return
		}

		// reread the dir info
		logDirFileInfo, statErr = os.Stat(dir)
		if statErr != nil && os.IsNotExist(statErr) {
			err = statErr
			return
		}
	}

	// check whether the directory or not
	if !logDirFileInfo.IsDir() {
		err = errors.New("should be directory")
		return
	}

	// check create file
	probeWritableFile := filepath.Join(dir, "probe.tmp")
	wErr := ioutil.WriteFile(probeWritableFile, []byte("probe writable"), 0644)
	if wErr != nil {
		err = fmt.Errorf("should be writable directory, %s", wErr.Error())
		return
	}
	defer os.Remove(probeWritableFile)

	// check create directory
	probeMkdirDir := filepath.Join(dir, "probe_tmp")
	mkErr := os.Mkdir(probeMkdirDir, 0755)
	if mkErr != nil {
		err = fmt.Errorf("should be writable directory, %s", mkErr.Error())
		return
	}
	defer os.Remove(probeMkdirDir)

	return
}
