package main

import (
   "math"
   "sync"
   "runtime"
   "C"
   "unsafe"
)

func paraquicksort(liste []float64, semChan chan struct{}) []float64 {

    if len(liste) <= 1 {
		return liste
	}
	
	// Le quicksort typique
	right, left := len(liste)-1, 0
	a := len(liste) / 2
	pivot := int(math.Floor(float64(a)))
	liste[right], liste[pivot] = liste[pivot], liste[right]

	for i, _ := range liste {
		if liste[i] < liste[right] {
			liste[left], liste[i] = liste[i], liste[left]
			left++
		}
	}
	liste[right], liste[left] = liste[left], liste[right]
	
    
    // On intialize le waitgroup
	wg := sync.WaitGroup{}

	select {
	case semChan <- struct{}{}:
		// On alloue le processus à une nouvelle goroutine
		wg.Add(1)
		go func() {
			paraquicksort(liste[:left], semChan)
			<-semChan
			wg.Done()
		}()
	default:
		// On ne peut plus créer de goroutine. Cette même goroutine s'occupera de le faire
		paraquicksort(liste[:left], semChan)
	}

    // Cette même goroutine s'occupera de la partie gauche de la liste
	paraquicksort(liste[left+1:], semChan)
    
    wg.Wait()
	return liste
	
}

func QuickSort(src []float64) []float64{
	// On précise le nombre de cpus et goroutines à utiliser
	runtime.GOMAXPROCS(4)
	extraGoroutines := 8
	// On initialize le channel qui limitera le nombre de goroutines
	semChan := make(chan struct{}, extraGoroutines)
	defer close(semChan)
	return paraquicksort(src, semChan)
}

//export Sort
func Sort(cArray *C.double, cSize C.int)  uintptr {
    gSlice := (*[1 << 30]C.double)(unsafe.Pointer(cArray))[:cSize:cSize]
    result := make([]float64, cSize)
    for i := 0; i < int(cSize); i++ {
        result[i] = float64(gSlice[i])
    }
    result = QuickSort(result)
    return uintptr(unsafe.Pointer(&result[0]))
}

func main() {}