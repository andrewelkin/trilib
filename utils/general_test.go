package utils

import "testing"

func TestMakeSureNegativeIf(t *testing.T) {

	a := MakeSureNegativeIf(true, -1.2)
	if a > 0 {
		t.Errorf("error MakeSureNegativeIf: must be negative")
	}
	a = MakeSureNegativeIf(true, 1.2)
	if a > 0 {
		t.Errorf("error MakeSureNegativeIf: must be negative")
	}

	a = MakeSureNegativeIf(false, 1.2)
	if a < 0 {
		t.Errorf("error MakeSureNegativeIf: must be positive")
	}
	a = MakeSureNegativeIf(false, -1.2)
	if a < 0 {
		t.Errorf("error MakeSureNegativeIf: must be positive")
	}

}

func TestGetValueOfDefault(t *testing.T) {

	m := map[string]int{
		"a": 1,
		"b": 2,
	}

	a := GetValueOrDefault(m, "a", 33)
	if a != 1 {
		t.Errorf("error GetValueOfDefault: must return 1")
	}
	a = GetValueOrDefault(m, "z", 33)
	if a != 33 {
		t.Errorf("error GetValueOfDefault: must return 33")
	}

}
