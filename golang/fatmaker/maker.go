package fatmaker

import (
	"debug/macho"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	MagicFat64 = macho.MagicFat + 1
	alignBits = 14
	align = 1 << alignBits
)

type input struct {
	data []byte
	cpu uint32
	subcpu uint32
	offset int64
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <output file> <input file 1> <input file 2> ...\n", os.Args[0])
		os.Exit(2)
	}

	var inputs []input
	offset := int64(align)
	for _, item := range os.Args[2:] {
		data, err := ioutil.ReadFile(item)
		if err != nil {
			panic(err)
		}
		if len(data) < 2 {
			panic(fmt.Sprintf("file %s too small", item))
		}

		// all currently supported mac archs are little endian.
		magic := binary.LittleEndian.Uint32(data[0:4])
		if magic != macho.Magic32 && magic != macho.Magic64 {
			panic(fmt.Sprintf("input %s is not a macho file, magic=%x", item, magic))
		}

		cpu := binary.LittleEndian.Uint32(data[4:8])
		subcpu := binary.LittleEndian.Uint32(data[8:12])
		inputs = append(inputs, input{data: data, cpu: cpu, subcpu: subcpu, offset: offset})
		offset += int64(len(data))
		offset = (offset + align - 1) / align * align
	}

	// decide on whether we are doing fat32 or fat64
	sixtyfour := false
	if len(inputs) == 0 {
		panic(fmt.Sprintf("inputs file error"))
	}
	if inputs[len(inputs)-1].offset >= 1<<32 || len(inputs[len(inputs)-1].data) >= 1<<32 {
		sixtyfour = true
		// TODO now, there are still some problems in 64, skip for now
		panic("files too large to fit into a fat binary")
	}

	// make output file
	out, err := os.Create(os.Args[1])
	if err != nil {
		panic(err)
	}
	err = out.Chmod(0755)
	if err != nil {
		panic(err)
	}

	// build a fat_header
	var hdr []uint32
	if sixtyfour {
		hdr = append(hdr, MagicFat64)
	} else {
		hdr = append(hdr, macho.MagicFat)
	}
	hdr = append(hdr, uint32(len(inputs)))

	// build a fat_arc for each input file
	for _, item := range inputs {
		hdr = append(hdr, item.cpu)
		hdr = append(hdr, item.subcpu)
		if sixtyfour {
			hdr = append(hdr, uint32(item.offset>>32)) // big endian
			hdr = append(hdr, uint32(len(item.data)>>32))
		} else {
			hdr = append(hdr, uint32(item.offset))
			hdr = append(hdr, uint32(len(item.data)))
		}
		hdr = append(hdr, alignBits)
		if sixtyfour {
			hdr = append(hdr, 0)
		}
	}

	err = binary.Write(out, binary.BigEndian, hdr)
	if err != nil {
		panic(err)
	}
	offset = int64(4 * len(hdr))

	// write each contained file.
	for _, item := range inputs {
		if offset < item.offset {
			_, err = out.Write(make([]byte, item.offset-offset))
			if err != nil {
				panic(err)
			}
			offset = item.offset
		}
		_, err := out.Write(item.data)
		if err != nil {
			panic(err)
		}
		offset += int64(len(item.data))
	}

	defer func() {
		if err = out.Close(); err != nil {
			panic(err)
		}
	}()
}
