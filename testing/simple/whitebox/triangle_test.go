// We can see all members of the package 'simple'
package simple

import (
	"testing"
)

func TestCalcArea(t *testing.T) {
	var a, b, c uint = 3, 4, 5
	var expectedArea = 6.0

	area, err := CalcArea(a, b, c)

	if err != nil {
		t.Fatalf("Should be no errors when calculating the area")
	}
	if area == expectedArea {
		t.Logf("The area of the triangle (%d, %d, %d) should be equal to %f", a, b, c, expectedArea)
	} else {
		t.Errorf("The area of the triangle (%d, %d, %d) should be equal to %f, but got %f", a, b, c, expectedArea, area)
	}
}

func TestSemiperimeter(t *testing.T) {
	var a, b, c uint = 1, 2, 3
	var expectedSp = 3.0

	sp := calcSemiperimeter(a, b, c)

	if sp == expectedSp {
		t.Logf("The semiperimeter of the triangle (%d, %d, %d) should be equal to %f", a, b, c, expectedSp)
	} else {
		t.Errorf("The semiperimeter of the triangle (%d, %d, %d) should be equal to %f, but got %f", a, b, c, expectedSp, sp)
	}
}
