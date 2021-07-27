// +build !windows

package self

const BinaryFilename = "cl_update"
const SignatureFilename = "cl_update.sig"
const BinaryPath = "/usr/local/bin/" + BinaryFilename
const PublicKeyPath = "assets/public.pem"