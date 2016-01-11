package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Open itself for reading.
	r, err := zip.OpenReader(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	dir, err := ioutil.TempDir("", "russian-doll-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	log.Print("using temp dir ", dir)

	// Iterate through the files in the archive.
	for _, f := range r.File {
		// On OS X, this works:
		//out, err := exec.Command(filepath.Join(os.Args[0], f.Name)).CombinedOutput()
		//
		// But the more general solution is to explicitly extract the
		// embedded binary out of the main binary:
		out, err := execEmbeddedBinary(dir, f)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(string(out))
	}
}

func execEmbeddedBinary(tmpDir string, f *zip.File) ([]byte, error) {
	tmpName := filepath.Join(tmpDir, f.Name)
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	tmpFile, err := os.OpenFile(tmpName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(tmpFile, rc)
	if err != nil {
		return nil, err
	}
	err = tmpFile.Close()
	if err != nil {
		return nil, err
	}
	return exec.Command(tmpName).CombinedOutput()
}
