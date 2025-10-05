// We can see only exported members of the package 'simple'
package simple_test

import (
	simple "go-code-patterns/testing/simple/blackbox"
	"testing"
)

func TestCalcArea(t *testing.T) {
	var a, b, c uint = 3, 4, 5
	var expectedArea = 6.0

	// we import the func using the package name
	area, err := simple.CalcArea(a, b, c)
	if err != nil {
		t.Fatalf("Should be no errors when calculating the area")
	}
	if area == expectedArea {
		t.Logf("The area of the triangle (%d, %d, %d) should be equal to %f", a, b, c, expectedArea)
	} else {
		t.Errorf("The area of the triangle (%d, %d, %d) should be equal to %f, but got %f", a, b, c, expectedArea, area)
	}
}
