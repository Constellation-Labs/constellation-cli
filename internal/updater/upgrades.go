package updater

import "constellation/internal/updater/cli"
import "constellation/internal/updater/self"

func CommandlineUpgrade(version string) Upgrade {
	return  Upgrade {
		version,
		"https://constellationlabs-cli.s3.us-west-1.amazonaws.com/%s/%s/%s/%s",
		0755,
		cli.BinaryFilename,
		cli.ChecksumFilename,
		cli.BinaryPath,
	}
}

func SelfUpgrade(version string) Upgrade {
	return  Upgrade {
		version,
		"https://constellationlabs-cli.s3.us-west-1.amazonaws.com/%s/%s/%s/%s",
		0700,
		self.BinaryFilename,
		self.ChecksumFilename,
		self.BinaryPath,
	}
}