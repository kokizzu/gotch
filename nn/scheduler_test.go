package nn

import (
	// "reflect"
	// "fmt"
	"testing"

	"github.com/sugarme/gotch"
)

func TestLambdaLR(t *testing.T) {
	vs := NewVarStore(gotch.CPU)
	opt, err := DefaultAdamConfig().Build(vs, 0.001)
	if err != nil {
		t.Error(err)
	}

	ld1 := func(epoch interface{}) float64 {
		return float64(epoch.(int) / 30)
	}

	var s *LRScheduler
	s = NewLambdaLR(opt, []LambdaFn{ld1}).Build()

	wants := []float64{
		0.001, // initial LR
		0.001, // epoch 30s/30 = 1 * 0.001
		0.002, // epoch 60s/30 = 2 * 0.001 = 0.002
		0.003, // epoch 90s/30 = 3 * 0.001 = 0.003
	}
	i := 0
	for epoch := 0; epoch < 100; epoch++ {
		if epoch%30 == 0 && epoch > 0 {
			s.Step(epoch)
			i += 1
		}

		want := wants[i]
		got := opt.GetLRs()[0]
		if got != want {
			t.Errorf("Epoch %d: Want %v - Got %v", epoch, want, got)
		}
	}
}

func TestMultiplicativeLR(t *testing.T) {
	vs := NewVarStore(gotch.CPU)
	opt, err := DefaultAdamConfig().Build(vs, 0.001)
	if err != nil {
		t.Error(err)
	}

	ld1 := func(epoch interface{}) float64 {
		return float64(epoch.(int) / 30)
	}

	var s *LRScheduler
	s = NewMultiplicativeLR(opt, []LambdaFn{ld1}).Build()

	wants := []float64{
		0.001, // initial LR
		0.001, // epoch 30s/30 = 1 * 0.001
		0.002, // epoch 60s/30 = 2 * 0.001 = 0.002
		0.006, // epoch 90s/30 = 3 * 0.002 = 0.006
	}
	i := 0
	for epoch := 0; epoch < 100; epoch++ {
		if epoch%30 == 0 && epoch > 0 {
			s.Step(epoch)
			i += 1
		}

		want := wants[i]
		got := opt.GetLRs()[0]
		if got != want {
			t.Errorf("Epoch %d: Want %v - Got %v", epoch, want, got)
		}
	}
}

func TestStepLR(t *testing.T) {
	vs := NewVarStore(gotch.CPU)
	opt, err := DefaultAdamConfig().Build(vs, 0.05)
	if err != nil {
		t.Error(err)
	}

	var s *LRScheduler
	s = NewStepLR(opt, 30, 0.1).Build()

	wants := []float64{
		0.05,    // initial LR -> 0.05
		0.005,   // 30 <= epoch < 60 -> 0.05 * gamma = 0.005
		0.0005,  // 60 <= epoch < 90 -> 0.005 * gamma = 0.0005
		0.00005, // 90 <= epoch < 120 -> 0.0005 * gamma = 0.00005
	}
	i := 0
	for epoch := 0; epoch < 100; epoch++ {
		s.Step(epoch)
		if epoch%30 == 0 && epoch > 0 {
			i += 1
		}
		want := wants[i]
		got := opt.GetLRs()[0]
		if got != want {
			t.Errorf("Epoch %d: Want %v - Got %v", epoch, want, got)
		}
	}
}

func TestMultiStepLR(t *testing.T) {
	vs := NewVarStore(gotch.CPU)
	opt, err := DefaultAdamConfig().Build(vs, 0.05)
	if err != nil {
		t.Error(err)
	}

	var s *LRScheduler
	s = NewMultiStepLR(opt, []int{30, 80}, 0.1).Build()

	wants := []float64{
		0.05,   // initial LR -> 0.05
		0.005,  // 30 <= epoch < 80 -> 0.05 * gamma = 0.005
		0.0005, // 80 <= epoch -> 0.005 * gamma = 0.0005
	}
	i := 0
	for epoch := 0; epoch < 100; epoch++ {
		s.Step(epoch)
		if contain(epoch, []int{30, 80}) {
			i += 1
		}
		want := wants[i]
		got := opt.GetLRs()[0]
		if got != want {
			t.Errorf("Epoch %d: Want %v - Got %v", epoch, want, got)
		}
	}
}
