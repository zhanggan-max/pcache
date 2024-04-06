package purgekit

import "testing"

func TestArcGet(t *testing.T) {
	arc := NewARCache(128, nil)
	for i := 0; i < 128; i++ {
		arc.Add(i, i)
	}
	if arc.t1.Len() != 128 {
		t.Fatalf("arc should have 128 elements, but got: %v", arc.t1.Len())
	}
	if arc.t2.Len() != 0 {
		t.Fatalf("arc t2 should hold no elements, but got: %v", arc.t2.Len())
	}
	for i := 0; i < 128; i++ {
		arc.Get(i)
	}
	if arc.t2.Len() != 128 {
		t.Fatalf("arc t2 should have 128 elements, but got: %v", arc.t2.Len())
	}
	if arc.t1.Len() != 0 {
		t.Fatalf("arc t1 should hold no element, but got: %v", arc.t1.Len())
	}
	for i := 0; i < 128; i++ {
		arc.Get(i)
	}
	if arc.t2.Len() != 128 {
		t.Fatalf("arc t2 should have 128 elements, but got %v", arc.t2.Len())
	}
	if arc.t1.Len() != 0 {
		t.Fatalf("arc t1 should hold no element, but got %v", arc.t1.Len())
	}
}

func TestArcAdative(t *testing.T) {
	arc := NewARCache(4, nil)
	for i := 0; i < 4; i++ {
		arc.Add(i, i)
	}
	// t1 {0, 1, 2, 3}
	if arc.t1.Len() != 4 {
		t.Fatalf("t1 should have 4 elements, but got %v", arc.t1.Len())
	}
	arc.Get(0)
	arc.Get(1)
	// t1 {2, 3}, t2 {0, 1}
	if arc.t2.Len() != 2 {
		t.Fatalf("t2 should have 2 elements, but got: %v", arc.t2.Len())
	}
	arc.Add(4, 4)
	// t1 {3, 4}, t2 {0, 1}, b1 {2}
	if arc.b1.Len() != 1 {
		t.Fatalf("an element should be envicted from t1 to b1, but got: %v", arc.b1.Len())
	}
	arc.Add(2, 2)
	// t1 {4}, t2 {0, 1, 2}, b1 {3}, p = 1
	if arc.b1.Len() != 1 {
		t.Fatalf("b1 should have 1 element, but got: %v", arc.b1.Len())
	}
	if arc.p != 1 {
		t.Fatalf("p should be 1, but got: %v", arc.p)
	}
	if arc.t2.Len() != 3 {
		t.Fatalf("t2 should have 3 elements, but got: %v", arc.t2.Len())
	}
	arc.Add(4, 4)
	// t1 {}, t2 {0, 1, 2, 4}, b1 {3}
	if arc.t1.Len() != 0 {
		t.Fatalf("t1 should hold no element, but got: %v", arc.t1.Len())
	}
	if arc.t2.Len() != 4 {
		t.Fatalf("t2 should hold 4 elements, but got: %v", arc.t2.Len())
	}
	arc.Add(5, 5)
	// t1 {5}, t2 {1, 2, 4}, b1 {3}, b2 {0}, p = 1
	if arc.t1.Len() != 1 {
		t.Fatalf("t1 should have 1 element, but got: %v", arc.t1.Len())
	}
	if arc.t2.Len() != 3 {
		t.Fatalf("t2 should have 3 elements, but got: %v", arc.t2.Len())
	}
	if arc.b2.Len() != 1 {
		t.Fatalf("b2 should hold 1 element, but got: %v", arc.b2.Len())
	}
	arc.Add(0, 0)
	// t1 {}, t2 {1, 2, 4, 0}, b1 {3, 5}, b2 {}, p = 0
	if arc.t1.Len() != 0 {
		t.Fatalf("t1 should hold 1 element, but got: %v", arc.t1.Len())
	}
	if arc.t2.Len() != 4 {
		t.Fatalf("t2 should have 3 element, but got: %v", arc.t2.Len())
	}
	if arc.b1.Len() != 2 {
		t.Fatalf("b1 should have 1 element, but got: %v", arc.b1.Len())
	}
	if arc.b2.Len() != 0 {
		t.Fatalf("b2 should have 1 element, but got: %v", arc.b2.Len())
	}
	if arc.p != 0 {
		t.Fatalf("arc.p should be 1, but got: %v", arc.p)
	}
}
