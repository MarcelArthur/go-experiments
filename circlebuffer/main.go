package main

import (
	"errors"
	"fmt"
	"io"
	"time"
)

type CircleByteBuffer struct {
	io.Reader
	io.Writer
	io.Closer
	datas []byte


	start   int
	end     int
	size    int
	isClose bool
	isEnd   bool
}


func NewCircleByteBuffer(len int) *CircleByteBuffer{
	buffer := &CircleByteBuffer{
		start: 0,
		end:   0,
		datas: make([]byte, len),
		size:  len,
		isClose: false,
		isEnd:   false,
	}
	return buffer
}



func (buffer *CircleByteBuffer)getLen() int{
	if buffer.start == buffer.end{
		return 0
	}else if buffer.start < buffer.end{
		return buffer.end - buffer.start
	}else {
		return buffer.start - buffer.end
	}
}


func (buffer *CircleByteBuffer)getFree() int{
	return buffer.size - buffer.getLen()
}

func (buffer *CircleByteBuffer)putByte(b byte) error{
	if buffer.isClose {
		return io.EOF
	}
	buffer.datas[buffer.end] = b
	pos := buffer.end + 1
	for pos == buffer.end {
		if buffer.isClose{
			return io.EOF
			time.Sleep(time.Millisecond)

		}
	}
	if pos == buffer.size {
		buffer.end = 0
	}else{
		buffer.end = pos
	}
	return nil
}


func (buffer *CircleByteBuffer)getByte() (byte, error){
	if buffer.isClose{
		return 0, io.EOF
	}
	if buffer.isEnd && buffer.getLen() <= 0{
		return 0, io.EOF
	}
	if buffer.getLen() <= 0{
		return 0, errors.New("no data")
	}
	ret := buffer.datas[buffer.start]
	buffer.start++
	if buffer.start == buffer.size{
		buffer.start = 0
	}
	return ret, nil
}

func (buffer *CircleByteBuffer)geti(i int)byte{
	if i >= buffer.getLen(){
		panic("out buffer")
	}
	pos := buffer.start+i
	if pos >= buffer.size {
		pos -= buffer.size
	}
	return buffer.datas[pos]
}

func (buffer *CircleByteBuffer)Close() error{
	buffer.isClose = true
	return nil
}

func (buffer *CircleByteBuffer)Read(bts []byte) (int, error){
	if buffer.isClose {
		return 0, io.EOF
	}
	if bts == nil {
		return 0, errors.New("bts is nil")
	}
	ret := 0
	for i:=0;i<len(bts);i++{
		b,err := buffer.getByte()
		if err != nil {
			if err == io.EOF{
				return ret, err
			}
			return ret, nil
		}
		ret++
		bts[i] = b
	}

	if buffer.isClose{
		return ret, io.EOF
	}
	return ret, nil

}


func (buffer *CircleByteBuffer)Write(bts []byte) (int, error){
	if buffer.isClose{
		return 0, io.EOF
	}
	if bts == nil {
		buffer.isEnd = true
		return 0, io.EOF
	}
	ret := 0
	for i:=0; i<len(bts);i++{
		err := buffer.putByte(bts[i])
		if err != nil {
			fmt.Println("Write bts err:", err)
			return ret, err
		}
		ret++
	}

	if buffer.isClose{
		return ret, io.EOF
	}
	return ret, nil

}