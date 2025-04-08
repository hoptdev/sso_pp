/*package main

import (
	"fmt"
	jwt_helper "main/helpers"
	"testing"
)

func TestNewTokenWithValidate(t *testing.T) {
	testdata := []struct {
		userId   int
		role     int
		expected bool
	}{
		{
			userId:   1,
			role:     1,
			expected: true,
		},
		{
			role:     1,
			expected: true,
		},
	}

	for _, value := range testdata {
		result, err := jwt_helper.NewRefreshToken(value.userId, value.role)
		if err != nil {
			t.Errorf("Error is not nil: %s | %s", result, err)
		}

		if r, _, e := jwt_helper.ValidateToken(result); !r || e != nil {
			t.Errorf("Incorrect: %s | %s", result, err)
		}

		fmt.Println(result)
	}
}
*/