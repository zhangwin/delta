package delta

import (
	"reflect"
	"testing"
)

func TestPostProcessAdd(t *testing.T) {
	d := &DiffSolution{
		lines: [][3]string{
			{"A", "A", string(LineFromBoth)},
			{"", "B", string(LineFromB)},
			{"", "C", string(LineFromB)},
			{"", "D", string(LineFromB)},
			{"", "E", string(LineFromB)},
			{"B", "B", string(LineFromBoth)},
			{"C", "C", string(LineFromBoth)},
			{"D", "D", string(LineFromBoth)},
		},
	}

	e := &DiffSolution{
		lines: [][3]string{
			{"A", "A", string(LineFromBoth)},
			{"B", "B", string(LineFromBoth)},
			{"C", "C", string(LineFromBoth)},
			{"D", "D", string(LineFromBoth)},
			{"", "E", string(LineFromB)},
			{"", "B", string(LineFromB)},
			{"", "C", string(LineFromB)},
			{"", "D", string(LineFromB)},
		},
	}
	d.PostProcess()
	if !reflect.DeepEqual(d, e) {
		t.Errorf("expected:\n%+v\nbut got:\n%+v", e, d)
	}
}

func TestPostProcessDel(t *testing.T) {
	d := &DiffSolution{
		lines: [][3]string{
			{"A", "A", string(LineFromBoth)},
			{"B", "", string(LineFromA)},
			{"C", "", string(LineFromA)},
			{"D", "", string(LineFromA)},
			{"B", "B", string(LineFromBoth)},
			{"C", "C", string(LineFromBoth)},
			{"D", "D", string(LineFromBoth)},
		},
	}

	e := &DiffSolution{
		lines: [][3]string{
			{"A", "A", string(LineFromBoth)},
			{"B", "B", string(LineFromBoth)},
			{"C", "C", string(LineFromBoth)},
			{"D", "D", string(LineFromBoth)},
			{"B", "", string(LineFromA)},
			{"C", "", string(LineFromA)},
			{"D", "", string(LineFromA)},
		},
	}
	d.PostProcess()
	if !reflect.DeepEqual(d, e) {
		t.Errorf("expected:\n%+v\nbut got:\n%+v", e, d)
	}
}
