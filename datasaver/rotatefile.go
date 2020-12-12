package datasaver

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path"
	"sync/atomic"
	"time"

	"github.com/kisom/goutils/fileutil"
	"github.com/miolini/datacounter"
	cntlib "github.com/orian/counters"
	"github.com/orian/pbio"
	"github.com/orian/utils/file"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type StringGenerator func() string

// TimeBasedGenerator generates a filepath based on date and provided directory.
// If a file exists it adds a suffix before extention with order number.
func TimeBasedGenerator(baseDir string, perm os.FileMode) StringGenerator {
	return func() string {
		t := time.Now().UTC()
		dir := t.Format("2006/01/02")
		err := os.MkdirAll(path.Join(baseDir, dir), perm)
		if err != nil {
			logrus.Errorf("cannot create PB log directory: %s", err)
		}
		f := t.Format("150405") + ".pb.gz"
		f = path.Join(baseDir, dir, f)
		for i := 0; fileutil.FileDoesExist(f) && i < 100; i++ {
			f = fmt.Sprintf("%s.%.3d.pb.gz", t.Format("150405"), i)
			f = path.Join(baseDir, dir, f)
		}
		return f
	}
}

type Encoder interface {
	// Encode encodes a message. Expect a correct format (e.g. if PbEnc is used we want Protbuf message).
	Encode(m interface{}) error
	// Count returns a number of bytes written.
	Count() uint64
	// Created returns when encoder was created.
	Created() time.Time
	// Used returns how many times the encoder was used.
	Used() uint64

	io.Closer
}

type gobEnc struct {
	closer io.Closer
	writer *datacounter.WriterCounter
	enc    *gob.Encoder
	used   uint64
	t      time.Time
}

func (g *gobEnc) Encode(m interface{}) error {
	atomic.AddUint64(&g.used, 1)
	return g.enc.Encode(m)
}

func (g *gobEnc) Close() error {
	return g.closer.Close()
}

func (g *gobEnc) Count() uint64 {
	return g.writer.Count()
}

func (g *gobEnc) Created() time.Time {
	return g.t
}

func (g *gobEnc) Used() uint64 {
	return atomic.LoadUint64(&g.used)
}

func GobGen(s string) (Encoder, error) {
	f, err := file.NewFileWriter(s)
	if err != nil {
		return nil, err
	}
	wc := datacounter.NewWriterCounter(f)
	return &gobEnc{f, wc, gob.NewEncoder(wc), 0, time.Now()}, nil
}

type EncoderGenerator func(s string) (Encoder, error)

type pbEnc struct {
	closer io.Closer
	enc    pbio.WriteCloser
	writer *datacounter.WriterCounter
	used   uint64
	t      time.Time
}

func (g *pbEnc) Encode(m interface{}) error {
	atomic.AddUint64(&g.used, 1)
	return g.enc.WriteMsg(m.(proto.Message))
}

func (g *pbEnc) Close() error {
	return g.closer.Close()
}

func (g *pbEnc) Count() uint64 {
	return g.writer.Count()
}

func (g *pbEnc) Created() time.Time {
	return g.t
}

func (g *pbEnc) Used() uint64 {
	return atomic.LoadUint64(&g.used)
}

func PbGen(s string) (Encoder, error) {
	f, err := file.NewFileWriter(s)
	if err != nil {
		return nil, err
	}
	wc := datacounter.NewWriterCounter(f)
	return &pbEnc{f, pbio.NewDelimitedWriter(wc), wc, 0, time.Now()}, nil
}

type RotateFile struct {
	generator StringGenerator
	enc       Encoder
	filePath  string

	Cnt            cntlib.Counters
	FileClosedHook func(string, time.Time)
	EncGen         EncoderGenerator
}

func NewRotateFile(
	cnt cntlib.Counters,
	generator StringGenerator,
	encGen EncoderGenerator,
) *RotateFile {

	if cnt == nil {
		cnt = cntlib.New()
	}
	return &RotateFile{
		generator: generator,

		Cnt:    cnt,
		EncGen: encGen,
	}
}

func (r *RotateFile) close() (err error) {
	if r.enc == nil {
		return nil
	}
	if err := r.enc.Close(); err == nil {
		if r.enc.Used() == 0 {
			r.Cnt.Get("file-not-used-delete").Increment()
			logrus.Debugf("empty file will be removed: %s", r.filePath)
			if err = os.Remove(r.filePath); err != nil {
				logrus.Errorf("cannot remove empty file: %s", err)
			}
		} else if r.FileClosedHook != nil {
			r.Cnt.Get("file-close-hook").Increment()
			r.FileClosedHook(r.filePath, r.enc.Created())
		}
	}
	r.enc = nil
	return err
}

func (r *RotateFile) Refresh() (err error) {
	r.Cnt.Get("file-refresh").Increment()
	logrus.Debug("refresh file")
	if err = r.close(); err != nil {
		return err
	}
	f := r.generator()
	r.filePath = f
	logrus.Debugf("new file %s", f)
	r.enc, err = r.EncGen(f)
	return err
}

func (r *RotateFile) Close() (err error) {
	return r.close()
}

func (r *RotateFile) Encode(v interface{}) error {
	r.Cnt.Get("encode-value").Increment()
	return r.enc.Encode(v)
}
