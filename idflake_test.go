package idflake

import (
	"fmt"
	"testing"
)

func Test_NewIdflake(t *testing.T) {
	idflake, _ := NewIdflake(10)

	id, err := idflake.NextId()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println("sort id : ", id)

	if _, err := idflake.SetEpoch(0); err != nil {
		t.Error("SetEpoch error")
	}

	id, err = idflake.NextId()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("long id : ", id)
}

func Test_id(t *testing.T) {
	idflake, err := NewIdflake(10)
	if err != nil {
		t.Error("NewIdflake(10) error(%v)", err)
	}

	var lastid uint64 = 0
	var id uint64 = 0

	for i := 0; i < 10000; i++ {
		id, err = idflake.NextId()
		if err != nil {
			t.FailNow()
		}
		if id <= lastid {
			t.Error("idflake error : ", fmt.Sprintf("%d  <  %d ", id, lastid))
			t.FailNow()

		}
		lastid = id
	}
}

func Benchmark_Idflake(b *testing.B) {
	idflake, err := NewIdflake(10)
	if err != nil {
		b.Error("NewIdflake(10) ", fmt.Sprintf("error(%v)", err))
		b.FailNow()
	}
	for i := 0; i < b.N; i++ {
		if _, err := idflake.NextId(); err != nil {
			b.Error("NewIdflake(10) ", fmt.Sprintf("error(%v)", err))
			b.FailNow()
		}
	}
}
