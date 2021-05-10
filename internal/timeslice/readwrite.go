package timeslice

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

func Load(filename string) ([]int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	if size == 0 {
		return nil, errors.New("empty storage")
	}
	if size%8 != 0 {
		return nil, fmt.Errorf("invalid file size, x8 expected, got %d", size)
	}
	ints := make([]int64, size/8)
	err = binary.Read(bufio.NewReader(file), binary.LittleEndian, ints)
	return ints, err
}

func Save(buf []int64, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = binary.Write(file, binary.LittleEndian, buf)
	file.Close()
	return err
}
