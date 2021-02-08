package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/adler32"
	"io"
)

type Packet struct {
	TotalSize uint32
	Magic [4]byte
	Payload []byte
	Checksum uint32
}

var RPC_MAGIC = [4]byte{'p', 'y', 'x', 'i'}

func encodePacket(w io.Writer, payload []byte) error {
	// len(MAGIC) + len(Checksum) == 8
	totalSize := uint32(len(payload) + 8)

	// write total size
	binary.Write(w, binary.BigEndian, totalSize)

	// write magic byets
	binary.Write(w, binary.BigEndian, RPC_MAGIC)

	// write the payload
	w.Write(payload)

	// calculate checksum
	var buf bytes.Buffer
	buf.Write(RPC_MAGIC[:])
	buf.Write(payload)
	checksum := adler32.Checksum(buf.Bytes())

	// write checksum
	return binary.Write(w, binary.BigEndian, checksum)
}

func encodePackage2(w io.Writer, payload []byte) error {
	// len(Magic) + len(Checksum) == 8
	totalsize := uint32(len(RPC_MAGIC) + len(payload) + 4)
	// write total size
	binary.Write(w, binary.BigEndian, totalsize)

	// write magic bytes
	binary.Write(w, binary.BigEndian, RPC_MAGIC)

	// write payload
	w.Write(payload)

	// calculate checksum
	// 减少butes.Buffer的内存开销
	sum := adler32.New()
	sum.Write(RPC_MAGIC[:])
	sum.Write(payload)
	checksum := sum.Sum32()

	return binary.Write(w, binary.BigEndian, checksum)
}

func encodePackage3(w io.Writer, payload []byte) error {
	// len(Magic) + len(Checksum) == 8
	totalsize := uint32(len(RPC_MAGIC) + len(payload) + 4)
	// write total size
	binary.Write(w, binary.BigEndian, totalsize)

	sum := adler32.New()
	// MultiWriter 是将多个writer包装到一个writer里
	// 此后每次往ww写的话，同时会写到sum和w上
	ww := io.MultiWriter(sum, w)

	// write magic bytes
	binary.Write(ww, binary.BigEndian, RPC_MAGIC)

	// write payload
	ww.Write(payload)

	// calculate checksum
	// 这里sum已经有了magic 和 payload，可以直接计算
	checksum := sum.Sum32()

	// 最终还是要只写到w上
	return binary.Write(w, binary.BigEndian, checksum)
}

func DecodePacket(r io.Reader) ([]byte, error) {
	var totalSize uint32
	err := binary.Read(r, binary.BigEndian, &totalSize)
	if err != nil {
		return nil, err
	}
	if totalSize < 8 {
		return nil, errors.New("total size is too ")
	}

	sum := adler32.New()
	// 从r中读，同时会将读到的写到sum中
	rr := io.TeeReader(r, sum)

	var magic [4]byte
	err = binary.Read(rr, binary.BigEndian, &magic)
	if err != nil {
		return nil, err
	}
	if magic != RPC_MAGIC {

	}

	payload := make([]byte, totalSize-8)
	_, err = io.ReadFull(rr, payload)
	if err != nil {
		return nil, err
	}

	var checksum uint32
	err = binary.Read(r, binary.BigEndian, &checksum)
	if err != nil {
		return nil, err
	}

	// 这里已经将magic+payload都读完了，同时也写入到了sum，直接
	if checksum != sum.Sum32() {

	}

	return payload, nil

}



func main() {
	
}
