package a2022

import (
	"encoding/csv"
	"io"

	"github.com/gocarina/gocsv"
)

func day20_ll(in io.Reader) (int, int) {
	type intCol struct {
		V int
	}
	var rows []intCol
	if err := gocsv.UnmarshalCSVWithoutHeaders(csv.NewReader(in), &rows); err != nil {
		panic(err)
	}
	type linkedVal struct {
		prev, next *linkedVal
		val        int
	}
	type ogIndexToStateLL map[int]*linkedVal

	checkWithKey := func(decryptionKey, mixRounds int) int {
		var current *linkedVal
		var init *linkedVal
		fileMap := ogIndexToStateLL{}
		var zero *linkedVal
		for i, v := range rows {
			vs := v.V * decryptionKey
			lv := &linkedVal{current, nil, vs}
			fileMap[i] = lv
			if v.V == 0 {
				zero = lv
			}
			if init == nil {
				init = lv
			}
			if current != nil {
				current.next = lv
			}
			current = lv
		}
		current.next = init
		init.prev = current

		for mix := 0; mix < mixRounds; mix++ {
			for i := 0; i < len(fileMap); i++ {
				processing := fileMap[i]
				distance := processing.val % (len(fileMap) - 1)
				d := sign(processing.val)
				if d == 0 {
					continue
				}
				processing.next.prev, processing.prev.next = processing.prev, processing.next
				prev, next := processing.prev, processing.next
				for j := 0; j < abs(distance); j++ {
					prev, next = aElseB(d == -1, prev.prev, prev.next), aElseB(d == -1, next.prev, next.next)
				}
				prev.next, next.prev = processing, processing
				processing.prev, processing.next = prev, next
			}
		}
		coords := make([]int, 0, 3)
		current = zero
		for i := 0; i < 3; i++ {
			for j := 0; j < 1000; j++ {
				current = current.next
			}
			coords = append(coords, current.val)
		}
		return coords[0] + coords[1] + coords[2]
	}

	return checkWithKey(1, 1), checkWithKey(811589153, 10)
}
