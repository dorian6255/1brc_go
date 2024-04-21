package handling

import (
	"dorian6255_1brc/internals/dataType"
)

func HandleMax(in chan dataType.StoreUnit, maximum map[string]float32) {

	for elem := range in {
		max, ok := maximum[elem.Name]

		// f, _ := max.(float32)

		if !ok || *elem.Value > max {
			maximum[elem.Name] = *elem.Value

		}
	}

	dataType.Wgstore.Done()
	return
}
func HandleAvg(in chan dataType.StoreUnit, average, number map[string]float32) {

	for elem := range in {
		nb, ok := number[elem.Name]

		if ok {

			prevAvg, _ := average[elem.Name]

			Value := ((prevAvg * nb) + *elem.Value) / (nb + 1)
			average[elem.Name] = Value
		} else {
			var one float32 = 1.0
			number[elem.Name] = one

			average[elem.Name] = *elem.Value
		}
	}

	dataType.Wgstore.Done()
	return
}
func HandleMin(in chan dataType.StoreUnit, minimum map[string]float32) {

	for elem := range in {
		// min, ok := minimum.LoadOrStore(elem.Name, elem.Value)

		min, ok := minimum[elem.Name]

		if !ok || *elem.Value < min {
			minimum[elem.Name] = *elem.Value

		}
	}

	dataType.Wgstore.Done()
	return
}

func HandleValues(in chan dataType.ValueIn, minimum, average, maximum chan dataType.StoreUnit) {

	for elem := range in {
		minimum <- dataType.StoreUnit{Name: elem.Name, Value: elem.Value}
		maximum <- dataType.StoreUnit{Name: elem.Name, Value: elem.Value}
		average <- dataType.StoreUnit{Name: elem.Name, Value: elem.Value}
	}
	close(minimum)
	close(maximum)
	close(average)
	dataType.Wgstore.Done()
}
