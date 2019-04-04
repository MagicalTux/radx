package radx

import (
	"errors"
	"io"
	"log"
	"os"
	"sync"
)

type File struct {
	f    io.ReadWriteSeeker
	l    sync.Mutex // file access lock
	w    sync.Mutex // write lock
	last int        // id of last block
}

func Open(filename string) (*File, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	sz, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		// shouldn't happen
		return nil, err
	}

	last := sz / RadxBlockSize
	if sz%RadxBlockSize != 0 {
		// got a last non-full block
		last += 1
	}

	if last >= 0xffffffff {
		return nil, errors.New("file too large")
	}

	return &File{f: f, last: int(last)}, nil
}

// readBlock loads a data block from the file
func (f *File) readBlock(id int) []byte {
	if id >= f.last {
		// no such block
		return nil
	}

	f.l.Lock()
	defer f.l.Unlock()

	_, err := f.f.Seek(int64(id)*RadxBlockSize, io.SeekStart)
	if err != nil {
		// wat?
		log.Printf("radx: failed to seek in file: %s", err)
		return nil
	}

	buf := make([]byte, RadxBlockSize)
	_, err = io.ReadFull(f.f, buf)
	if err != nil && err != io.EOF {
		return nil
	}

	return buf
}

func (f *File) writeBlock(id int, data []byte) error {
	if len(data) > RadxBlockSize {
		return errors.New("block too large")
	}

	f.l.Lock()
	defer f.l.Unlock()

	// TODO also write to log file for disaster recovery?

	_, err := f.f.Seek(int64(id)*RadxBlockSize, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = f.f.Write(data)
	if err != nil {
		return err
	}

	if id >= f.last {
		f.last = id + 1
	}
	return nil
}
