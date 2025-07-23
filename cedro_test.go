package cedro

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_multiApi(t *testing.T) {
	t.Run("basic book operations", func(t *testing.T) {
		msgs := []string{
			"Z:WINQ25:U:9:V:135420:598:224:22071453",
			"Z:WINQ25:A:0:A:135385:4:1:22071455",
			"Z:WINQ25:U:9:V:135420:595:224:22071453",
			"Z:WINQ25:A:5:A:135395:40:1:22071455",
			"T:WINQ25:090619:2:134865:6:5:7:5:63:3:106:+:142:090619292!",
		}

		object := Object{}
		for _, msg := range msgs {
			object = Process(msg, object)
		}
		assert.Equal(t, 595, object.Symbols[0].AggregatedBookAsk[9].Volume)
		assert.Equal(t, 4, object.Symbols[0].AggregatedBookBid[0].Volume)
		assert.Equal(t, 6, object.Time.Minute())
	})
}
