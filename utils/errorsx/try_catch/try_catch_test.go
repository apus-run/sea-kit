package try_catch

import "testing"

func TestTryCatch(t *testing.T) {
	Try(func() (any, error) {
		return 0, nil
	}).Catch(func(_ any, err error) any {
		t.Fatalf("catch error: %v", err)
		return nil
	}).Final(func(result any) {
		t.Log("errors: everything is good")
	})

	Try(func() (any, error) {
		return 1, New("e")
	}).Catch(func(result any, err error) any {
		t.Logf("captured result: %v", result.(int))
		t.Logf("captured error: %v", err)
		return err
	}).Final(func(result any) {
		if result == nil {
			t.Fatalf("cannot capture error")
		}
	})

	Try(func() (any, error) {
		return 1, nil
	}).Final(func(r any) {
		if r.(int) != 1 {
			t.Fatalf("result from try block is not as expected")
		}
	})
}
