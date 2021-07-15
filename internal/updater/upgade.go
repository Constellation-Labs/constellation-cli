package updater

import (
	"bytes"
	"constellation/internal/updater/cli"
	"errors"
	"fmt"
	"os"
	"runtime"
)

// linux/386/1.0/cl_cli
//const RepositoryUrlPattern = "https://example.com/%s/%s/%s/%s"
//const ChecksumSize = 1
//const BinaryFileChmod = 0755

type Upgrade struct {
	TargetVersion string
	RepositoryUrlPattern string
	ChecksumSize int
	BinaryFileChmod os.FileMode
}

// func (n *node) GetNodeMetrics() (*Metrics, error) {
func (c Upgrade) binaryFileUrl(version string) string {
	return fmt.Sprintf(c.RepositoryUrlPattern, runtime.GOOS, runtime.GOARCH, version, cli.binaryFilename)
}

func (c Upgrade) checksumFileUrl(version string) string {
	return fmt.Sprintf(c.RepositoryUrlPattern, runtime.GOOS, runtime.GOARCH, version, cli.checksumFilename)
}

func (c Upgrade) calculateChecksum(file *os.File) []byte {
	return make([]byte, c.ChecksumSize)
}

func (c Upgrade) Run() error {
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
	checksum := make([]byte, c.ChecksumSize)
	checksumSize, checksumErr := checksumFile.Read(checksum)

	if checksumErr != nil || checksumSize != c.ChecksumSize || bytes.Compare(calculatedChecksum, checksum) != 0 {
		return errors.New("invalid checksum")
	}

	// make a binary file backup
	isUpgradeSuccessful := false
	backupBinaryFile := cli.binaryPath + ".bak"

	if createBackupError := os.Rename(cli.binaryPath, backupBinaryFile); createBackupError != nil {
		return errors.New("cannot replace binary file")
	}

	// register revert
	defer func() {
		if isUpgradeSuccessful == false {
			os.Rename(backupBinaryFile, cli.binaryPath)
		}
	}()

	// replace binary
	if upgradeBinaryError := os.Rename(upgradeFile.Name(), cli.binaryPath); upgradeBinaryError != nil {
		return upgradeBinaryError
	}

	if chmodBinaryError := os.Chmod(cli.binaryPath, c.BinaryFileChmod); chmodBinaryError != nil {
		return chmodBinaryError
	}

	isUpgradeSuccessful = true

	return nil
}
