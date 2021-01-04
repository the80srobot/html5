package bindings

import (
	"testing"

	"github.com/the80srobot/html5/safe"
)

func TestVarAttach(t *testing.T) {
	v := Declare("foo")
	var m Map

	v = m.Attach(v, safe.Default)

	// The variable should now be valid for the map. Test that out by using it
	// to bind a value.
	vm, err := m.Bind(v.SetConst("bar"))
	if err != nil {
		t.Fatal(err)
	}

	s := vm.GetString(&v)
	if s != "bar" {
		t.Errorf("ValueMap.GetString(%v) => %q, wanted %q", v, s, "bar")
	}
}

func TestTrustClimbing(t *testing.T) {
	var m Map
	v := m.Declare("comment_text", safe.TextSafe)
	if _, err := v.Set(safe.UntrustedString("Hello!")); err == nil {
		t.Error("Var.Set() should fail to set an untrusted string on TextSafe Var")
	}

	// Declaring the var again with the default trust shouldn't change anything
	// - text should still be the right level.
	v = m.Declare("comment_text", safe.Default)
	if _, err := v.Set(safe.Bless(safe.TextSafe, "Hello!")); err != nil {
		t.Errorf("Var.Set() of a TextSafe string: %v", err)
	}

	// This should climb all the way up to fully trusted, as the only way to
	// reconcile the two trust levels.
	v = m.Declare("comment_text", safe.URLSafe)
	if _, err := v.Set(safe.Bless(safe.TextSafe, "Hello!")); err == nil {
		t.Error("Var.Set() should have refused to set TextSafe after trust climbing to fully trusted")
	}

	// Fully trusted strings should still be accepted.
	if _, err := v.Set(safe.Bless(safe.FullyTrusted, "Hello!")); err != nil {
		t.Errorf("Var.Set() of a FullyTrusted string: %v", err)
	}
}

// func TestMapString(t *testing.T) {
// 	var m Map
// 	t1 := m.DeclareString("foo", Untrusted)
// 	m.DeclareString("bar", Untrusted)

// 	if t3 := m.DeclareString("foo", AttributeSafe); t3 != t1 {
// 		t.Errorf("Different tag when declaring the same value a second time (%d vs %d)", t3, t1)
// 	}

// 	if s := m.String(t1); s != s1 {
// 		t.Errorf("Lookup after declare yielded the wrong string %q (wanted %q)", s.Name, s1.Name)
// 	}
// }

// func BenchmarkMap100LookupBaseline(b *testing.B) {
// 	m := map[string]Var{}
// 	for i := 0; i < 100; i++ {
// 		k := fmt.Sprintf("key_%d", i)
// 		m[k] = Var{Name: k}
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		v := m["key_50"]
// 		// This check is mainly here to make sure the lookup isn't optimized
// 		// away.
// 		if v.Name == "" {
// 			b.Fatal("bad lookup")
// 		}
// 	}
// }

// func BenchmarkMap100Lookup(b *testing.B) {
// 	m := NewMap()
// 	tags := []Tag{}
// 	for i := 0; i < 100; i++ {
// 		k := fmt.Sprintf("key_%d", i)
// 		tags = append(tags, m.DeclareString(Var{Name: k}))
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		v := m.String(tags[50])
// 		// This check is here because we also have it in the baseline.
// 		if v.Name == "" {
// 			b.Fatal("bad lookup")
// 		}
// 	}
// }
