/*
Copyright 2018 Iguazio Systems Ltd.

Licensed under the Apache License, Version 2.0 (the "License") with
an addition restriction as set forth herein. You may not use this
file except in compliance with the License. You may obtain a copy of
the License at http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
implied. See the License for the specific language governing
permissions and limitations under the License.

In addition, you may not use the software for any purposes that are
illegal under applicable law, and the grant of the foregoing license
under the Apache 2.0 license is conditioned upon your compliance with
such restriction.
*/

package frames

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func TestMarshal(t *testing.T) {
	frame := createFrame(t)

	marshaler, ok := frame.(Marshaler)
	if !ok {
		t.Fatalf("frame is not Marshaler")
	}

	out, err := marshaler.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	msg, ok := out.(*MapFrameMessage)
	if !ok {
		t.Fatalf("wrong message type - %T", msg)
	}
}

func TestRoundTrip(t *testing.T) {
	frame1 := createFrame(t)

	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	if err := enc.Encode(frame1); err != nil {
		t.Fatal(err)
	}

	dec := NewDecoder(&buf)
	frame2, err := dec.Decode()
	if err != nil {
		t.Fatal(err)
	}

	cols1, cols2 := frame1.Columns(), frame2.Columns()
	if !reflect.DeepEqual(cols1, cols2) {
		t.Fatalf("columns mismatch: %v != %v", cols1, cols2)
	}

}

func createFrame(t *testing.T) Frame {
	var (
		columns []Column
		col     Column
		err     error
	)

	col, err = NewSliceColumn("icol", []int{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}

	columns = append(columns, col)
	col, err = NewSliceColumn("fcol", []float64{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}

	columns = append(columns, col)
	col, err = NewSliceColumn("scol", []string{"1", "2", "3"})
	if err != nil {
		t.Fatal(err)
	}

	columns = append(columns, col)
	col, err = NewSliceColumn("tcol", []time.Time{time.Now(), time.Now(), time.Now()})
	if err != nil {
		t.Fatal(err)
	}

	col, err = NewLabelColumn("lcol", "srv", 3)
	if err != nil {
		t.Fatal(err)
	}

	columns = append(columns, col)
	frame, err := NewMapFrame(columns, nil)
	if err != nil {
		t.Fatal(err)
	}

	return frame
}