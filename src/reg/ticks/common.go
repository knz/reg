package ticks

import (
	"reg/t"
)

func TeeTicks(dst1 chan t.Ticks, dst2 chan t.Ticks, src chan t.Ticks) {
	for a := range src {
		dst1 <- a
		dst2 <- a
	}
}
