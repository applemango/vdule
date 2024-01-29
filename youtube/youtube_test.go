package youtube

import "testing"

func TestCreateTube(t *testing.T) {
	_, err := CreateTube()
	if err != nil {
		t.Fatalf("%v\n", err)
	}
}
