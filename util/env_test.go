package util_test

import (
	"testing"

	"github.com/dc-dc-dc/cheetah/util"
)

func TestGetEnv(t *testing.T) {
	envChecking := "xxxxxx-fake-env"
	if util.GetEnv(envChecking, "default") != "default" {
		t.Errorf("expected default value")
	}
	t.Setenv(envChecking, "value")
	if util.GetEnv(envChecking, "default") != "value" {
		t.Errorf("expected value")
	}
}

func TestIsTruthy(t *testing.T) {
	trueVals := []string{"true", "True", "TRUE", "1", "t", "T"}
	falseVals := []string{"false", "False", "FALSE", "0", "f", "F"}
	envChecking := "xxxxxx-fake-env"
	for _, val := range trueVals {
		t.Setenv(envChecking, val)
		if !util.IsTruthy(envChecking) {
			t.Errorf("expected %s to be true", val)
		}
	}
	for _, val := range falseVals {
		t.Setenv(envChecking, val)
		if util.IsTruthy(envChecking) {
			t.Errorf("expected %s to be false", val)
		}
	}

}
