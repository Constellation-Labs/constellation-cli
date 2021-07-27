package updater

import (
	"errors"
	"fmt"
	"os"
	"runtime"
)

type Upgrade struct {
	TargetVersion string
	RepositoryUrlPattern string
	BinaryFileChmod os.FileMode // FIXME: do we need this?
	BinaryFilename string
	SignatureFilename string
	BinaryPath string
	PublicKeyPath string
}

func (c Upgrade) binaryFileUrl(version string) string {
	arch := runtime.GOARCH
	if arch == "arm64" {
		arch = "amd64"
	}

	return fmt.Sprintf(c.RepositoryUrlPattern, runtime.GOOS, arch, version, c.BinaryFilename)
}

func (c Upgrade) signatureFileUrl(version string) string {
	arch := runtime.GOARCH
	if arch == "arm64" {
		arch = "amd64"
	}

	return fmt.Sprintf(c.RepositoryUrlPattern, runtime.GOOS, arch, version, c.SignatureFilename)
}


func (c Upgrade) Run() error {
	// FIXME: Verify current privileges for the executable? It does not bring much value though since someone might have tampered it alraedy

	// download binary file
	berr, upgradeFile := DownloadFile(c.binaryFileUrl(c.TargetVersion))
	serr, signatureFile := DownloadFile(c.signatureFileUrl(c.TargetVersion))

	defer func() {
		os.Remove(upgradeFile.Name())
		os.Remove(signatureFile.Name())
	}()

	if berr != nil || serr != nil {
		return errors.New("failed downloading")
	}

	err := Verify(upgradeFile.Name(), signatureFile.Name(), c.PublicKeyPath)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid signature: %s", err))
	}

	// make a binary file backup
	isUpgradeSuccessful := false
	backupBinaryFile := c.BinaryPath + ".bak"

	// TODO: error when file does not exist (first run)
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
