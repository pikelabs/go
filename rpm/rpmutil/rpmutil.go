package rpmutil

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"code.pikelabs.net/go/archive/cpio"
	"code.pikelabs.net/go/rpm"
)

type Package struct {
	SigHeader *rpm.Header
	Header    *rpm.Header
	r         *readCounter
}

func OpenFile(fname string) (*Package, error) {
	d, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	f := bytes.NewReader(d)
	return ReadPackage(f)
}

func ReadPackage(r io.Reader) (*Package, error) {
	rc := &readCounter{r: r}
	// first lets get rid of Lead, for sanity check
	// will still use first LeadMagic to make sure
	// file is considered to be RPM, but rpm itself
	// no longe uses lead structure. it's there for
	// `file` command

	lead := make([]byte, rpm.LeadSize)
	if _, err := io.ReadFull(rc, lead); err != nil {
		return nil, err
	}

	magic := binary.BigEndian.Uint32(lead[0:4])
	if magic&0xFFFFFFFF != rpm.LeadMagic {
		return nil, errors.New("bad lead magic")
	}

	sigHeader, err := rpm.ReadHeader(rc)
	if err != nil {
		return nil, err
	}

	// signature header padded to align to 8 bytes
	psize := (rc.n + 7) / 8 * 8
	skip := int64(psize - rc.n)

	if _, err := io.CopyN(ioutil.Discard, rc, skip); err != nil {
		return nil, err
	}

	header, err := rpm.ReadHeader(rc)

	if err != nil {
		return nil, err
	}
	pkg := &Package{
		SigHeader: sigHeader,
		Header:    header,
		r:         rc,
	}
	return pkg, nil
}

func (pkg *Package) Payload() (cpio.Reader, error) {
	plRdr, err := decompressPkgPayload(pkg)
	if err != nil {
		return nil, err
	}
	return cpio.NewReader(plRdr)
}

func (pkg *Package) Files() ([]FileInfo, error) {
	var (
		paths     []string
		err, err1 error
	)
	paths, err = pkg.Header.GetStrings(rpm.TagFilenames)
	if err != nil {
		paths, err1 = pkg.Header.GetStrings(rpm.TagOldFilenames)
		if err1 != nil {
			return nil, err1
		}
	}
	files := make([]FileInfo, len(paths))
	for i := 0; i < len(files); i++ {
		files[i].Name = paths[i]
	}
	return files, nil
}

func (pkg *Package) Dump(w io.Writer) error {
	tags := pkg.Header.AvailableTags()

	for i, tag := range tags {

		if name, ok := rpm.HeaderNames[tag]; ok {
			fmt.Fprintf(w, "%d: (%d) %s\n", i, tag, name)
		} else {
			fmt.Fprintf(w, "%d: (%d) \n", i, tag)
		}
		ttype, tdata, err := pkg.Header.GetTag(tag)
		if err != nil {
			fmt.Fprintf(w, "err: %s\n", err.Error())
		}
		switch ttype {
		case rpm.DataTypeString:
			fmt.Fprintf(w, "\t %s\n", tdata)
		case rpm.DataTypeStringArray:
			fmt.Fprintf(w, "\t %v\n", rpm.CStringArrayToSlice(tdata))
		}
	}
	return nil
}

type readCounter struct {
	n int
	r io.Reader
}

func (rc *readCounter) Read(b []byte) (n int, err error) {
	n, err = rc.r.Read(b)
	rc.n += n
	return
}
