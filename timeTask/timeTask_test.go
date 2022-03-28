package timeTask

import (
	"testing"
	"time"
)

func TestTimeTask(t *testing.T) {
	count := 0
	NewTask().SetInterval(time.Millisecond).SetConsumer(func() {
		count+=1
	})
	time.Sleep(2*time.Second)
	if count<1 {
		t.Error("time task error",count)
	}
}