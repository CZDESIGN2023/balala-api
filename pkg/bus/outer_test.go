package bus

import "testing"

func TestSubscribe(t *testing.T) {
	On("event:fruit", "1", t.Log)

	On("event:fruit", "2", func(data string) error {
		t.Log("2", data)
		return nil
	})

	Emit("event:fruit", "apple")

	Off("event:fruit", "1")

	Emit("event:fruit", "banana")

}
