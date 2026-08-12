package vm

import (
	"testing"
)

func TestMemory(t *testing.T) {
	mem := NewMemory(1024)

	t.Run("general", func(t *testing.T) {
		_, err := mem.Load(0)

		if err == nil {
			t.Fatal("wanted error for Load(0), found nil")
		}

		_, err = mem.Load(1)

		if err == nil {
			t.Fatal("wanted error for Load(1), found nil")
		}

		_, err = mem.Load(2048)

		if err == nil {
			t.Fatal("wanted error for Load out of bounds, found nil")
		}

		var valueToStore uint32 = 1023
		err = mem.Store(4, valueToStore)

		if err != nil {
			t.Fatalf("Store failed %s", err)
		}

		val, err := mem.Load(4)

		if err != nil {
			t.Fatalf("Load failed %s", err)
		}

		if val != valueToStore {
			t.Fatalf("wanted %d, got %d", valueToStore, val)
		}
	})

	t.Run("endianness", func(t *testing.T) {
		err := mem.Store(8, 0x01020304)

		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}

		tests := [][]uint32{
			{8, 4},
			{9, 3},
			{10, 2},
			{11, 1},
		}

		for _, test := range tests {
			addr := test[0]
			want := test[1]
			got := mem.data[addr]

			if want != uint32(got) {
				t.Fatalf("data[%d]: wanted %d, got %d", addr, want, got)
			}
		}
	})
}
