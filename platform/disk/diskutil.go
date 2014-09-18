package disk

import (
	"path/filepath"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshdevutil "github.com/cloudfoundry/bosh-agent/platform/deviceutil"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
)

type diskUtil struct {
	diskPath string
	mounter  Mounter
	fs       boshsys.FileSystem
	logger   boshlog.Logger
	logTag   string
}

func NewDiskUtil(diskPath string, mounter Mounter, fs boshsys.FileSystem, logger boshlog.Logger) boshdevutil.DeviceUtil {
	return diskUtil{
		diskPath: diskPath,
		mounter:  mounter,
		fs:       fs,
		logger:   logger,
		logTag:   "diskUtil",
	}
}

func (util diskUtil) GetFileContents(fileName string) ([]byte, error) {
	if !util.fs.FileExists(util.diskPath) {
		return []byte{}, bosherr.New("Failed to get file contents, disk path '%s' does not exist", util.diskPath)
	}

	tempDir, err := util.fs.TempDir("diskutil")
	if err != nil {
		return []byte{}, bosherr.WrapError(err, "Creating temporary disk mount point")
	}
	defer util.fs.RemoveAll(tempDir)

	err = util.mounter.Mount(util.diskPath, tempDir)
	if err != nil {
		return []byte{}, bosherr.WrapError(err, "Mounting disk path %s to %s", util.diskPath, tempDir)
	}
	util.logger.Debug(util.logTag, "Mounted disk path %s to %s", util.diskPath, tempDir)

	diskFilePath := filepath.Join(tempDir, fileName)
	util.logger.Debug(util.logTag, "Reading contents of %s", diskFilePath)
	contents, err := util.fs.ReadFile(diskFilePath)
	if err != nil {
		return []byte{}, bosherr.WrapError(err, "Reading from disk file %s", diskFilePath)
	}
	util.logger.Debug(util.logTag, "Got contents of %s: %s", diskFilePath, string(contents))

	_, err = util.mounter.Unmount(util.diskPath)
	util.logger.Debug(util.logTag, "Unmounting disk path", util.diskPath)
	if err != nil {
		return []byte{}, bosherr.WrapError(err, "Unmounting path %s", util.diskPath)
	}

	return contents, nil
}
