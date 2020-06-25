package main

import (
	"C"
	"math"
	"runtime"
	"sync"
	"time"
	"unsafe"
        "fmt"
)
import (
	"os"
	
)

func quicksort(A []float64) []float64 {
	if len(A) <= 1 {
		return A
	} else {

		low, high := len(A)-1, 0
		a := len(A) / 2
		pivot := int(math.Floor(float64(a)))
		A[low], A[pivot] = A[pivot], A[low]

		for i, _ := range A {
			if A[i] < A[low] {
				A[high], A[i] = A[i], A[high]
				high++
			}
		}
		A[low], A[high] = A[high], A[low]

		quicksort(A[:high])
		quicksort(A[high+1:])

		return A
	}
}

func paraquicksort(liste []float64, semChan chan struct{}) []float64 {

	if len(liste) <= 1 {
		return liste
	}

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

	wg := sync.WaitGroup{}

	select {
	case semChan <- struct{}{}:
		wg.Add(1)
		go func() {
			paraquicksort(liste[:left], semChan)
			<-semChan
			wg.Done()
		}()
	default:
		// Can't create a new goroutine, let's do the job ourselves.
		paraquicksort(liste[:left], semChan)
	}

	paraquicksort(liste[left+1:], semChan)

	wg.Wait()
	return liste

}

func QuickSort(src []float64) []float64 {
	runtime.GOMAXPROCS(40)
	extraGoroutines := 8
	semChan := make(chan struct{}, extraGoroutines)
	defer close(semChan)
	return paraquicksort(src, semChan)
}

//export Sort
func Sort(cArray *C.double, cSize C.int, para string) uintptr {
	gSlice := (*[1 << 30]C.double)(unsafe.Pointer(cArray))[:cSize:cSize]
	result := make([]float64, cSize)
	for i := 0; i < int(cSize); i++ {
		result[i] = float64(gSlice[i])
	}
	if para == "Y" {
		result = QuickSort(result)
	} else {
		result = quicksort(result)
	}
	return uintptr(unsafe.Pointer(&result[0]))
}

//export Temps
func Temps(cArray *C.double, cSize C.int) {
	times := []float64{}
	start := time.Now()
	Sort(cArray, cSize, "Y")
	end := time.Now().Sub(start).Seconds()
	times = append(times, float64(end))
	start = time.Now()
	Sort(cArray, cSize, "N")
	end = time.Now().Sub(start).Seconds()
	times = append(times, float64(end))
        
	file, err := os.OpenFile("fichier_intermediaire_temps.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	file.WriteString(fmt.Sprintf("%f", times[0]))
	file.WriteString("S")
	file.WriteString(fmt.Sprintf("%f", times[1]))
	file.Close()
}

func main() {}
