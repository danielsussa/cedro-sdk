package process

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
		}

		symbols := make([]*Symbol, 0)
		for _, msg := range msgs {
			symbols = Process(msg, symbols)
		}
		assert.Equal(t, 595, symbols[0].AggregatedBookAsk[9].Volume)
		assert.Equal(t, 4, symbols[0].AggregatedBookBid[0].Volume)
	})
}
