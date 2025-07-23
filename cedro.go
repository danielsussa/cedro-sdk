package cedro

import (
	"slices"
	"strconv"
	"strings"
	"time"
)

type (
	Object struct {
		Symbols []*Symbol
		Time    time.Time
	}

	Symbol struct {
		Symbol          string
		Bid             float64
		Ask             float64
		Last            float64
		TheoreticalRate float64
		Participation   float64

		BookLineAsk []*BookLine
		BookLineBid []*BookLine

		AggregatedBookAsk []*AggregatedBook
		AggregatedBookBid []*AggregatedBook
	}

	AggregatedBook struct {
		Price       float64
		Volume      int
		TotalOrders int
		DT          time.Time
	}

	BookLine struct {
		Price  float64
		Volume int
		Broker int
		DT     time.Time
	}
)

func Process(msg string, object Object) Object {
	msgSpl := strings.Split(msg, ":")
	switch msgSpl[0] {
	case "T":
		var symbol *Symbol
		object.Symbols, symbol = getBySymbol(msgSpl[1], object.Symbols)
		t, err := time.Parse("150405", msgSpl[2])
		if err == nil {
			object.Time = t
		}

		for i := 3; i < len(msgSpl); i += 2 {
			switch msgSpl[i] {
			case "2": // Preço do último negócio
				price := stringToFloat(msgSpl[4])
				symbol.Last = price
			case "3": // Melhor oferta de compra
				price := stringToFloat(msgSpl[4])
				symbol.Bid = price
			case "4": // Melhor oferta de venda
				price := stringToFloat(msgSpl[4])
				symbol.Ask = price
			}
		}
	case "Z":
		switch msgSpl[2] {
		case "U":
			var symbol *Symbol
			object.Symbols, symbol = getBySymbol(msgSpl[1], object.Symbols)

			position := stringToInt(msgSpl[3])
			direction := msgSpl[4]
			price := stringToFloat(msgSpl[5])
			volume := stringToInt(msgSpl[6])
			totalOrders := stringToInt(msgSpl[7])

			bl := &AggregatedBook{
				Price:       price,
				Volume:      volume,
				TotalOrders: totalOrders,
			}
			if direction == "V" {
				symbol.AggregatedBookAsk = addOnPositionAgr(symbol.AggregatedBookAsk, bl, position)
			} else {
				symbol.AggregatedBookBid = addOnPositionAgr(symbol.AggregatedBookBid, bl, position)
			}
		case "A":
			var symbol *Symbol
			object.Symbols, symbol = getBySymbol(msgSpl[1], object.Symbols)

			position := stringToInt(msgSpl[3])
			direction := msgSpl[4]
			price := stringToFloat(msgSpl[5])
			volume := stringToInt(msgSpl[6])
			totalOrders := stringToInt(msgSpl[7])

			bl := &AggregatedBook{
				Price:       price,
				Volume:      volume,
				TotalOrders: totalOrders,
			}
			if direction == "V" {
				symbol.AggregatedBookAsk = addOnPositionAgr(symbol.AggregatedBookAsk, bl, position)
			} else {
				symbol.AggregatedBookBid = addOnPositionAgr(symbol.AggregatedBookBid, bl, position)
			}
		case "D":
			switch msgSpl[3] {
			case "1": // O tipo 1 indica que somente a oferta da posição indicada deve ser cancelada.
				var symbol *Symbol
				object.Symbols, symbol = getBySymbol(msgSpl[1], object.Symbols)

				position := stringToInt(msgSpl[3])
				direction := msgSpl[4]
				if direction == "A" {
					symbol.BookLineBid = removeAt(symbol.BookLineBid, position)
				} else {
					symbol.BookLineAsk = removeAt(symbol.BookLineAsk, position)
				}

			case "2": // O tipo 2 indica que todas as ofertas melhores do que a oferta indicada pela posição, inclusive ela, devem ser canceladas.
				var symbol *Symbol
				object.Symbols, symbol = getBySymbol(msgSpl[1], object.Symbols)

				position := stringToInt(msgSpl[3])
				direction := msgSpl[4]
				if direction == "A" {
					symbol.BookLineBid = removeFirstN(symbol.BookLineBid, position)
				} else {
					symbol.BookLineAsk = removeFirstN(symbol.BookLineAsk, position)
				}
			case "3": // O tipo 3 indica que todas as ofertas, tanto de compra quanto de venda, devem ser canceladas. Neste tipo a mensagem não vem acompanhada de direção e posição.
				var symbol *Symbol
				object.Symbols, symbol = getBySymbol(msgSpl[1], object.Symbols)

				symbol.BookLineBid = make([]*BookLine, 0)
				symbol.BookLineAsk = make([]*BookLine, 0)
			}
		}
	case "B":
		switch msgSpl[2] {
		case "A": // process BQT ADD
			var symbol *Symbol
			object.Symbols, symbol = getBySymbol(msgSpl[1], object.Symbols)

			position := stringToInt(msgSpl[3])
			direction := msgSpl[4]
			price := stringToFloat(msgSpl[5])
			volume := stringToInt(msgSpl[6])
			broker := stringToInt(msgSpl[7])
			dt := stringToTime(msgSpl[8])

			bl := &BookLine{
				Price:  price,
				Volume: volume,
				Broker: broker,
				DT:     dt,
			}
			if direction == "A" {
				symbol.BookLineBid = addOnPosition(symbol.BookLineBid, bl, position)
			} else {
				symbol.BookLineAsk = addOnPosition(symbol.BookLineAsk, bl, position)
			}
		case "D": // process BQT DELETE
			switch msgSpl[3] {
			case "1": // O tipo 1 indica que somente a oferta da posição indicada deve ser cancelada.
				var symbol *Symbol
				object.Symbols, symbol = getBySymbol(msgSpl[1], object.Symbols)

				position := stringToInt(msgSpl[3])
				direction := msgSpl[4]
				if direction == "A" {
					symbol.BookLineBid = removeAt(symbol.BookLineBid, position)
				} else {
					symbol.BookLineAsk = removeAt(symbol.BookLineAsk, position)
				}

			case "2": // O tipo 2 indica que todas as ofertas melhores do que a oferta indicada pela posição, inclusive ela, devem ser canceladas.
				var symbol *Symbol
				object.Symbols, symbol = getBySymbol(msgSpl[1], object.Symbols)

				position := stringToInt(msgSpl[3])
				direction := msgSpl[4]
				if direction == "A" {
					symbol.BookLineBid = removeFirstN(symbol.BookLineBid, position)
				} else {
					symbol.BookLineAsk = removeFirstN(symbol.BookLineAsk, position)
				}
			case "3": // O tipo 3 indica que todas as ofertas, tanto de compra quanto de venda, devem ser canceladas. Neste tipo a mensagem não vem acompanhada de direção e posição.
				var symbol *Symbol
				object.Symbols, symbol = getBySymbol(msgSpl[1], object.Symbols)

				symbol.BookLineBid = make([]*BookLine, 0)
				symbol.BookLineAsk = make([]*BookLine, 0)
			}
		case "U": // process BQT UPDATE
			var symbol *Symbol
			object.Symbols, symbol = getBySymbol(msgSpl[1], object.Symbols)

			newPosition := stringToInt(msgSpl[3])
			oldPosition := stringToInt(msgSpl[4])
			direction := msgSpl[5]
			price := stringToFloat(msgSpl[6])
			volume := stringToInt(msgSpl[7])
			broker := stringToInt(msgSpl[8])
			dt := stringToTime(msgSpl[9])

			bl := &BookLine{
				Price:  price,
				Volume: volume,
				Broker: broker,
				DT:     dt,
			}
			if direction == "V" {
				symbol.BookLineAsk = removeAt(symbol.BookLineAsk, oldPosition)
				symbol.BookLineAsk = addOnPosition(symbol.BookLineAsk, bl, newPosition)
			} else {
				symbol.BookLineBid = removeAt(symbol.BookLineBid, oldPosition)
				symbol.BookLineBid = addOnPosition(symbol.BookLineBid, bl, newPosition)
			}

		}
	}

	return object
}

func getBySymbol(symbol string, list []*Symbol) ([]*Symbol, *Symbol) {
	idx := slices.IndexFunc(list, func(i *Symbol) bool {
		return i.Symbol == symbol
	})
	if idx == -1 {
		s := &Symbol{Symbol: symbol}
		list = append(list, s)
		return list, s
	}

	return list, list[idx]
}

func stringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func stringToFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func stringToTime(s string) time.Time {
	// 22071052
	s = "2025" + s
	dt, _ := time.Parse("200602011504", s)
	return dt
}

func addOnPosition(list []*BookLine, bl *BookLine, pos int) []*BookLine {
	if list == nil {
		list = make([]*BookLine, 0)
	}
	for {
		if len(list) > pos {
			break
		}
		list = append(list, nil)
	}
	list[pos] = bl
	return list
}

func addOnPositionAgr(list []*AggregatedBook, bl *AggregatedBook, pos int) []*AggregatedBook {
	if list == nil {
		list = make([]*AggregatedBook, 0)
	}
	for {
		if len(list) > pos {
			break
		}
		list = append(list, nil)
	}
	list[pos] = bl
	return list
}

func removeFirstN(list []*BookLine, pos int) []*BookLine {
	for i := pos - 1; i >= 0; i-- {
		list[i] = nil
	}

	return list
}
func removeAt(list []*BookLine, pos int) []*BookLine {
	if len(list) <= pos {
		return list
	}
	list[pos] = nil
	return list
}
