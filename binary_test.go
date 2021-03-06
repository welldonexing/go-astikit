package astikit

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBitsWriter(t *testing.T) {
	// TODO Need to test LittleEndian
	bw := &bytes.Buffer{}
	w := NewBitsWriter(BitsWriterOptions{Writer: bw})
	err := w.Write("000000")
	if err != nil {
		t.Errorf("expected no error, got %+v", err)
	}
	if e, g := 0, bw.Len(); e != g {
		t.Errorf("expected %d, got %d", e, g)
	}
	err = w.Write(false)
	if err != nil {
		t.Errorf("expected no error, got %+v", err)
	}
	err = w.Write(true)
	if err != nil {
		t.Errorf("expected no error, got %+v", err)
	}
	if e, g := []byte{1}, bw.Bytes(); !reflect.DeepEqual(e, g) {
		t.Errorf("expected %+v, got %+v", e, g)
	}
	err = w.Write([]byte{2, 3})
	if err != nil {
		t.Errorf("expected no error, got %+v", err)
	}
	if e, g := []byte{1, 2, 3}, bw.Bytes(); !reflect.DeepEqual(e, g) {
		t.Errorf("expected %+v, got %+v", e, g)
	}
	err = w.Write(uint8(4))
	if err != nil {
		t.Errorf("expected no error, got %+v", err)
	}
	if e, g := []byte{1, 2, 3, 4}, bw.Bytes(); !reflect.DeepEqual(e, g) {
		t.Errorf("expected %+v, got %+v", e, g)
	}
	err = w.Write(uint16(5))
	if err != nil {
		t.Errorf("expected no error, got %+v", err)
	}
	if e, g := []byte{1, 2, 3, 4, 0, 5}, bw.Bytes(); !reflect.DeepEqual(e, g) {
		t.Errorf("expected %+v, got %+v", e, g)
	}
	err = w.Write(uint32(6))
	if err != nil {
		t.Errorf("expected no error, got %+v", err)
	}
	if e, g := []byte{1, 2, 3, 4, 0, 5, 0, 0, 0, 6}, bw.Bytes(); !reflect.DeepEqual(e, g) {
		t.Errorf("expected %+v, got %+v", e, g)
	}
	err = w.Write(uint64(7))
	if err != nil {
		t.Errorf("expected no error, got %+v", err)
	}
	if e, g := []byte{1, 2, 3, 4, 0, 5, 0, 0, 0, 6, 0, 0, 0, 0, 0, 0, 0, 7}, bw.Bytes(); !reflect.DeepEqual(e, g) {
		t.Errorf("expected %+v, got %+v", e, g)
	}
	err = w.Write(1)
	if err == nil {
		t.Error("expected error")
	}
	bw.Reset()
	err = w.WriteN(uint8(8), 3)
	if err != nil {
		t.Errorf("expected no error, got %+v", err)
	}
	err = w.WriteN(uint16(4096), 13)
	if err != nil {
		t.Errorf("expected no error, got %+v", err)
	}
	if e, g := []byte{136, 0}, bw.Bytes(); !reflect.DeepEqual(e, g) {
		t.Errorf("expected %+v, got %+v", e, g)
	}
}

func testByteHamming84Decode(i uint8) (o uint8, ok bool) {
	p1, d1, p2, d2, p3, d3, p4, d4 := i>>7&0x1, i>>6&0x1, i>>5&0x1, i>>4&0x1, i>>3&0x1, i>>2&0x1, i>>1&0x1, i&0x1
	testA := p1^d1^d3^d4 > 0
	testB := d1^p2^d2^d4 > 0
	testC := d1^d2^p3^d3 > 0
	testD := p1^d1^p2^d2^p3^d3^p4^d4 > 0
	if testA && testB && testC {
		// p4 may be incorrect
	} else if testD && (!testA || !testB || !testC) {
		return
	} else {
		if !testA && testB && testC {
			// p1 is incorrect
		} else if testA && !testB && testC {
			// p2 is incorrect
		} else if testA && testB && !testC {
			// p3 is incorrect
		} else if !testA && !testB && testC {
			// d4 is incorrect
			d4 ^= 1
		} else if testA && !testB && !testC {
			// d2 is incorrect
			d2 ^= 1
		} else if !testA && testB && !testC {
			// d3 is incorrect
			d3 ^= 1
		} else {
			// d1 is incorrect
			d1 ^= 1
		}
	}
	o = uint8(d4<<3 | d3<<2 | d2<<1 | d1)
	ok = true
	return
}

func TestByteHamming84Decode(t *testing.T) {
	for i := 0; i < 256; i++ {
		v, okV := ByteHamming84Decode(uint8(i))
		e, okE := testByteHamming84Decode(uint8(i))
		if !okE {
			if okV {
				t.Error("expected false, got true")
			}
		} else {
			if !okV {
				t.Error("expected true, got false")
			}
			if !reflect.DeepEqual(e, v) {
				t.Errorf("expected %+v, got %+v", e, v)
			}
		}
	}
}

func testByteParity(i uint8) bool {
	return (i&0x1)^(i>>1&0x1)^(i>>2&0x1)^(i>>3&0x1)^(i>>4&0x1)^(i>>5&0x1)^(i>>6&0x1)^(i>>7&0x1) > 0
}

func TestByteParity(t *testing.T) {
	for i := 0; i < 256; i++ {
		v, okV := ByteParity(uint8(i))
		okE := testByteParity(uint8(i))
		if !okE {
			if okV {
				t.Error("expected false, got true")
			}
		} else {
			if !okV {
				t.Error("expected true, got false")
			}
			if e := uint8(i) & 0x7f; e != v {
				t.Errorf("expected %+v, got %+v", e, v)
			}
		}
	}
}
