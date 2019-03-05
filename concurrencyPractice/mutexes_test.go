package mutexes

import "testing"

func TestChanDoSum1000(t *testing.T) {
	want :=1000
	got := ChanDoSum(want)
	if got!=want {
		t.Fatalf("Expected %d, got %d",want,got)
	}
}

func TestMutexDoSum1000(t *testing.T) {
	want :=1000
	got := MutexDoSum(want)
	if got!=want {
		t.Fatalf("Expected %d, got %d",want,got)
	}
}

func BenchmarkMutexDoSum(b *testing.B) {
	for n:=0;n<b.N;n++ {
		MutexDoSum(1000)
	}
}

func BenchmarkChanDoSum(b *testing.B) {
	for n:=0;n<b.N;n++ {
		ChanDoSum(1000)
	}
}