package bindings

import (
	"fmt"
	"testing"

	"github.com/the80srobot/html5/safe"
)

func BenchmarkPopulateGoMap100Keys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := map[string]string{}
		for i := 0; i < 100; i++ {
			k := fmt.Sprintf("key_%d", i)
			v := fmt.Sprintf("val_%d", i)
			m[k] = v
		}
	}
}

func BenchmarkPopulateValueMap100Bindings(b *testing.B) {
	var m Map

	vars := make([]Var, 100)
	for i := 0; i < 100; i++ {
		vars[i] = m.Declare(fmt.Sprintf("binding_%d", i), safe.Default)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		vm := m.MustBind()
		for i := 0; i < 100; i++ {
			vm.Set(vars[i].Bind(safe.Bless(safe.Default, fmt.Sprintf("val_%d", i))))
		}
	}
}

func BenchmarkLookupGoMap100Keys(b *testing.B) {
	m := map[string]string{}
	for i := 0; i < 100; i++ {
		k := fmt.Sprintf("key_%d", i)
		v := fmt.Sprintf("val_%d", i)
		m[k] = v
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := m["key_50"]
		// This check is mainly here to make sure the lookup isn't optimized
		// away.
		if val != "val_50" {
			b.Fatalf("bad lookup: key_50 = %q", val)
		}
	}
}

func BenchmarkLookupValueMap100Bindings(b *testing.B) {
	var m Map
	vm := m.MustBind()

	vars := make([]Var, 100)
	for i := 0; i < 100; i++ {
		vars[i] = m.Declare(fmt.Sprintf("binding_%d", i), safe.Default)
		vm.Set(vars[i].Bind(safe.Bless(safe.Default, fmt.Sprintf("val_%d", i))))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := vm.GetString(vars[50])
		// This check is here because we also have it in the baseline.
		if val != "val_50" {
			b.Fatalf("bad lookup: key_50 = %q", val)

		}
	}
}
