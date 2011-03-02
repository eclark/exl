package gridfs

import (
	"crypto/md5"
	"fmt"
	"github.com/eclark/exl/bson"
	"github.com/eclark/exl/bson/bsoncompat"
	"github.com/eclark/gomongo/mongo"
	"hash"
	"os"
	"time"
)

const (
	filesSuffix  = ".files"
	chunksSuffix = ".chunks"
	defaultPrefix = "fs"
	defaultChunksize = 256 * 1024
)

type File struct {
	Id mongo.BSON
	Chunksize int
	Size int64
	Nchunk int

	buf []byte
	nextc int
	pos int

	md5 hash.Hash
	db *mongo.Database
	prefix string
	filename string
}

func Open(filename string, db *mongo.Database, prefix string) (*File, os.Error) {
	query, err := mongo.Marshal(map[string]string{"filename": filename})
	if err != nil {
		return nil, err
	}

	file := new(File)
	file.db = db
	if prefix == "" {
		file.prefix = defaultPrefix
	} else {
		file.prefix = prefix
	}
	filem, err := db.GetCollection(file.prefix + filesSuffix).FindOne(query)
	if err != nil {
		return nil, err
	}

	switch filem.Get("length").Kind() {
	case mongo.IntKind:
		file.Size = int64(filem.Get("length").Int())
	case mongo.LongKind:
		file.Size = filem.Get("length").Long()
	default:
		return nil, os.NewError("No length for file!")
	}

	file.Chunksize = int(filem.Get("chunkSize").Int())
	file.Nchunk = int(file.Size / int64(file.Chunksize))
	if file.Size % int64(file.Chunksize) > 0 {
		file.Nchunk++
	}

	file.Id = filem.Get("id_")

	return file, nil
}

func (f *File) Read(b []byte) (n int, err os.Error) {
	for {
		switch {
		case len(b) == 0:
			return
		case len(f.buf) > 0:
			m := copy(b, f.buf)
			f.pos += m
			n += m
			b = b[m:]
			f.buf = f.buf[m:]
		case f.nextc < f.Nchunk:
			query, err := mongo.Marshal(map[string]interface{}{"files_id": f.Id, "n": int32(f.nextc)})
			if err != nil {
				return n, err
			}

			chunkm, err := f.db.GetCollection(f.prefix + chunksSuffix).FindOne(query)
			if err != nil {
				return n, err
			}

			f.buf = chunkm.Get("data").Binary()
			f.nextc++
		case f.nextc == f.Nchunk:
			err = os.EOF
			return
		default:
			panic("should never be reached!")
		}
	}
	return
}

func New(filename string, db *mongo.Database, prefix string) (file *File, err os.Error) {
	file = new(File)
	file.filename = filename
	file.db = db
	if prefix == "" {
		file.prefix = defaultPrefix
	} else {
		file.prefix = prefix
	}

	file.Chunksize = defaultChunksize
	file.Id, err = mongo.NewOID()
	if err != nil {
		return
	}
	file.buf = make([]byte, 0, file.Chunksize)
	file.md5 = md5.New()

	return
}

func (f *File) Write(b []byte) (n int, err os.Error) {
	for {
		switch {
			case len(b) == 0:
				return
			case len(f.buf) < cap(f.buf):
				// copy some bytes from b to f.buf
				amt := cap(f.buf) - len(f.buf)
				if amt > len(b) {
					amt = len(b)
				}
				pe := len(f.buf)
				f.buf = f.buf[:pe+amt]
				m := copy(f.buf[pe:], b[:amt])
				n += m
				b = b[m:]
			case len(f.buf) == f.Chunksize:
				// insert chunk in mongo
				f.writeChunk()
		}
	}
	return
}

func (f *File) writeChunk() os.Error {
	chunk, err := mongo.Marshal(map[string]interface{}{"files_id": f.Id, "n": f.nextc, "data": f.buf})
	if err != nil {
		return err
	}

	err = f.db.GetCollection(f.prefix + chunksSuffix).Insert(chunk)
	if err != nil {
		return err
	}

	f.md5.Write(f.buf)
	f.nextc++
	f.pos += len(f.buf)
	f.buf = f.buf[0:0]

	return nil
}

type filemd5Query struct {
	filemd5 mongo.BSON
	root string
}

func (f *File) Close() os.Error {
	f.writeChunk()

	md5cmddoc := new(bson.Document)
	oid := bson.ObjectId(f.Id.Bytes())
	pre := bson.String(f.prefix)
	md5cmddoc.Append("filemd5", &oid)
	md5cmddoc.Append("root", &pre)

	res, err := f.db.Command(bsoncompat.Wrap(md5cmddoc))
	if err != nil {
		return err
	}
	fmt.Println(res.Bytes())
	fmt.Println(res.Get("errmsg"))

	file, err := mongo.Marshal(map[string]interface{}{"_id": f.Id, "length": int32(f.pos), "chunkSize": int32(f.Chunksize), "uploadDate": time.LocalTime(), "md5": res.Get("md5").String(), "filename": f.filename})
	if err != nil {
		return err
	}

	err = f.db.GetCollection(f.prefix + filesSuffix).Insert(file)
	if err != nil {
		return err
	}

	return nil
}
