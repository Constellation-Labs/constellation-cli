package updater

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type WriteMeter struct {
	Total uint64
}

func (wc WriteMeter) PrintProgressMeter() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading update file... %s complete", humanize.Bytes(wc.Total))
}

func (wc *WriteMeter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgressMeter()
	return n, nil
}


func DownloadFile(url string) (error, *os.File) {

	out, err := ioutil.TempFile("", "constellation-*")

	if err != nil || out.Chmod(600) != nil { // TODO: handle chmod error
		return err, nil // TODO: This is potentially unsecure since we treat chmod error as no error but body is invalid anyway
	}

	resp, err := http.Get(url)

	if err != nil {
		out.Close()
		return err, nil
	}

	defer resp.Body.Close()
	defer out.Close()

	counter := &WriteMeter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		return err, nil
	}

	return nil, out
}

