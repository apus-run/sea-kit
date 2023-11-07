package wg

import (
	"log"
	"testing"
)

func TestWrapper(t *testing.T) {
	var wg = New()
	wg.Wrap(func() {
		log.Println("this is test")
	})

	for i := 0; i < 10; i++ {
		num := i
		// wrap go goroutine without recovery func
		wg.Wrap(func() {
			log.Println("current index:", num)
		})
	}

	wg.WrapWithRecovery(func() {
		log.Println("exec goroutine with recovery func")
		var s = []string{"a", "b", "c"}
		log.Printf("s[3] = %v", s[3])
	}, func(r interface{}) {
		// exec recover:runtime error: index out of range [3] with length 3
		log.Printf("exec recover:%v", r)
	})

	wg.Wait()
}
