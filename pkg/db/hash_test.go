package db

import (
	"testing"
)

func TestHash(t *testing.T) {
	tests := map[string][]string{
		"hello": {
			"HellO",
			"HELlo",
			"hElLO",
		},
	}

	for k, v := range tests {
		one := HashPassword(k)
		two := HashPassword(k)
		if one != two {
			t.Error("HASH didn't equal res in eq value")
		}
		for _, el := range v {
			two = HashPassword(el)
			if one == two {
				t.Errorf("HASH equal %s : %s", k, el)
			}
		}
	}
}
