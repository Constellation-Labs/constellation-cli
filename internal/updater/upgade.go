package updater

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime"
)

//TODO: FIX this size according to the checksum impl
const ChecksumSize = 1

type Upgrade struct {
	TargetVersion string
	RepositoryUrlPattern string
	BinaryFileChmod os.FileMode // FIXME: do we need this?
	BinaryFilename string
	ChecksumFilename string
	BinaryPath string
}

func (c Upgrade) binaryFileUrl(version string) string {
	return fmt.Sprintf(c.RepositoryUrlPattern, runtime.GOOS, runtime.GOARCH, version, c.BinaryFilename)
}

func (c Upgrade) checksumFileUrl(version string) string {
	return fmt.Sprintf(c.RepositoryUrlPattern, runtime.GOOS, runtime.GOARCH, version, c.ChecksumFilename)
}

func (c Upgrade) calculateChecksum(file *os.File) []byte {
	return make([]byte, ChecksumSize)
}

func (c Upgrade) Run() error {
	// FIXME: Verify current privileges for the executable? It does not bring much value though since someone might have tampered it alraedy

	// download binary file
	berr, upgradeFile := DownloadFile(c.binaryFileUrl(c.TargetVersion))
	cerr, checksumFile := DownloadFile(c.checksumFileUrl(c.TargetVersion))

	defer func() {
		os.Remove(upgradeFile.Name())
		os.Remove(checksumFile.Name())
	}()

	if berr != nil || cerr != nil {
		return errors.New("failed downloading")
	}

	// verify checksum

	calculatedChecksum := c.calculateChecksum(upgradeFile)
	checksum := make([]byte, ChecksumSize)
	checksumSize, checksumErr := checksumFile.Read(checksum)

	if checksumErr != nil || checksumSize != ChecksumSize || bytes.Compare(calculatedChecksum, checksum) != 0 {
		return errors.New("invalid checksum")
	}

	// make a binary file backup
	isUpgradeSuccessful := false
	backupBinaryFile := c.BinaryPath + ".bak"

	if createBackupError := os.Rename(c.BinaryPath, backupBinaryFile); createBackupError != nil {
		return errors.New("cannot replace binary file")
	}

	// register revert
	defer func() {
		if isUpgradeSuccessful == false {
			os.Rename(backupBinaryFile, c.BinaryPath)
		}
	}()

	// replace binary
	if upgradeBinaryError := os.Rename(upgradeFile.Name(), c.BinaryPath); upgradeBinaryError != nil {
		return upgradeBinaryError
	}

	if chmodBinaryError := os.Chmod(c.BinaryPath, c.BinaryFileChmod); chmodBinaryError != nil {
		return chmodBinaryError
	}

	isUpgradeSuccessful = true

	return nil
}
