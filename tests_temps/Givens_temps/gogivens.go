package main

import (
	"C"
	"math"
	"sync"
	"time"

	"./go.matrix"
)
import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)

func QRPara(j int, r *matrix.DenseMatrix, q *matrix.DenseMatrix, ch chan string, wg *sync.WaitGroup) {
	ch2 := make(chan string, 1)
	l, _ := r.GetSize()
	index := 0
	state := "halt"
	for i := l - 2; i > j-1; i-- {
		if state != "ended" {
			for state == "halt" {
				select {
				case state = <-ch:
				default:
					state = "halt"
				}
			}
		}
		x := r.Arrays()[i][j]
		y := r.Arrays()[i+1][j]
		if y != 0 {
			t := math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
			cos := x / t
			sin := y / t
			rot := matrix.MakeDenseMatrixStacked([][]float64{
				{cos, -sin},
				{sin, cos},
			})
			SubMatrixTimesLeft(r.GetMatrix(i, j, 2, l-j), rot.Transpose())
			SubMatrixTimesRight(q.GetMatrix(0, i, l, 2), rot)
			index++
			if index == 2 && j < l {
				go QRPara(j+1, r, q, ch2, wg)
			}
		}
		select {
		case ch2 <- "go":
		default:
		}

	}
	select {
	case <-ch:
	default:
	}
	select {
	case ch2 <- "ended":
	default:
	}

	wg.Done()
}

func SubMatrixTimesLeft(SubA *matrix.DenseMatrix, B *matrix.DenseMatrix) {

	temp := matrix.ParallelProduct(B, SubA)
	b, _ := B.GetSize()
	_, a := SubA.GetSize()
	for i := 0; i < b; i++ {
		for j := 0; j < a; j++ {
			SubA.Arrays()[i][j] = temp.Arrays()[i][j]
		}
	}
}

func SubMatrixTimesRight(SubA *matrix.DenseMatrix, B *matrix.DenseMatrix) {

	temp := matrix.ParallelProduct(SubA, B)
	_, b := B.GetSize()
	a, _ := SubA.GetSize()
	for i := 0; i < a; i++ {
		for j := 0; j < b; j++ {
			SubA.Arrays()[i][j] = temp.Arrays()[i][j]
		}
	}
}

func QR(m *matrix.DenseMatrix) (*matrix.DenseMatrix, *matrix.DenseMatrix) {
	l, _ := m.GetSize()
	r := matrix.MakeDenseCopy(m)
	q := matrix.Eye(l)
	for j := 0; j < l-1; j++ {
		for i := l - 2; i > j-1; i-- {
			x := r.Arrays()[i][j]
			y := r.Arrays()[i+1][j]
			if y != 0 {
				t := math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
				cos := x / t
				sin := y / t
				rot := matrix.MakeDenseMatrixStacked([][]float64{
					{cos, -sin},
					{sin, cos},
				})
				SubMatrixTimesLeft(r.GetMatrix(i, j, 2, l-j), rot.Transpose())
				SubMatrixTimesRight(q.GetMatrix(0, i, l, 2), rot)

			}
		}
	}
	return q, r
}

//export Interfacage1
func Interfacage1() {
	file, err := os.OpenFile("fichier_intermediaire.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close() // on ferme automatiquement à la fin de notre programme
	data, err := ioutil.ReadFile("fichier_intermediaire.txt")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	lenght := 0
	incr := 0
	current_int := 0
	for string(lines[0][incr]) != "S" {
		current_int, err = strconv.Atoi(string(lines[0][incr]))
		if err != nil {
			panic(err)
		}
		lenght += int(math.Pow10(incr)) * current_int
		incr++
	}
	Matrix := matrix.Eye(lenght)
	i, q, float := 0, 0, ""
	for line := 1; line < (len(lines) - 1); line++ {
		q, i = 0, 0
		for i < len(lines[line]) {
			// On s'occupe d'une ligne
			float = ""
			for string(lines[line][i]) != "S" {
				float = float + string(lines[line][i])
				i++
			}
			Matrix.Arrays()[line-1][q], err = strconv.ParseFloat(float, 64)
			if err != nil {
				panic(err)
			}
			q++
			i++
		}
	}

	runtime.GOMAXPROCS(40)
	R := matrix.MakeDenseCopy(Matrix)
	l, _ := Matrix.GetSize()
	Q := matrix.Eye(l)
	var wg sync.WaitGroup
	ch := make(chan string, 1)
	ch <- "ended"
	wg.Add(int(l - 1))
	go QRPara(0, R, Q, ch, &wg)
	wg.Wait()
	os.Remove("fichier_intermediaire.txt")
	fileR, err := os.OpenFile("fichier_intermediaire_R.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	fileQ, err := os.OpenFile("fichier_intermediaire_Q.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	for ligne := 0; ligne < lenght; ligne++ {
		for colonne := 0; colonne < lenght; colonne++ {
			fileQ.WriteString(fmt.Sprintf("%f", Q.Arrays()[ligne][colonne]))
			fileR.WriteString(fmt.Sprintf("%f", R.Arrays()[ligne][colonne]))
			fileR.WriteString("S")
			fileQ.WriteString("S")
		}
		fileR.WriteString("\n")
		fileQ.WriteString("\n")
	}
	fileQ.Close()
	fileR.Close()
}

//export Interfacage2
func Interfacage2(cArray *C.double, cSize C.int, rows C.int, cols C.int) uintptr {
	gSlice := (*[1 << 30]C.double)(unsafe.Pointer(cArray))[:cSize:cSize]
	result := make([]float64, cSize)
	for i := 0; i < int(cSize); i++ {
		result[i] = float64(gSlice[i])
	}
	A := matrix.MakeDenseMatrix(result, int(rows), int(cols))

	runtime.GOMAXPROCS(4)
	r := matrix.MakeDenseCopy(A)
	l, _ := A.GetSize()
	q := matrix.Eye(l)
	var wg sync.WaitGroup
	ch := make(chan string, 1)
	ch <- "ended"
	wg.Add(int(rows) - 1)
	go QRPara(0, r, q, ch, &wg)
	wg.Wait()

	return uintptr(unsafe.Pointer(&append(r.Array(), q.Array()...)[0]))
}

//export Temps
func Temps(cArray *C.double, cSize C.int, rows C.int, cols C.int) {
	gSlice := (*[1 << 30]C.double)(unsafe.Pointer(cArray))[:cSize:cSize]
	result := make([]float64, cSize)
	times := []float64{}
	for i := 0; i < int(cSize); i++ {
		result[i] = float64(gSlice[i])
	}
	A := matrix.MakeDenseMatrix(result, int(rows), int(cols))
	//on calcule le temps avec parallélisation
	runtime.GOMAXPROCS(4)
	r := matrix.MakeDenseCopy(A)
	l, _ := A.GetSize()
	q := matrix.Eye(l)
	start := time.Now()
	var wg sync.WaitGroup
	ch := make(chan string, 1)
	ch <- "ended"
	wg.Add(int(rows) - 1)
	go QRPara(0, r, q, ch, &wg)
	wg.Wait()
	end := time.Now().Sub(start).Seconds()
	times = append(times, float64(end))
	//on calcule le temps sans parallélisation
	start = time.Now()
	r = matrix.MakeDenseCopy(A)
	q = matrix.Eye(l)
	r, q = QR(A)
	end = time.Now().Sub(start).Seconds()
	times = append(times, float64(end))
	fmt.Println(times)
	file, err := os.OpenFile("fichier_intermediaire_temps.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	file.WriteString(fmt.Sprintf("%f", times[0]))
	file.WriteString("S")
	file.WriteString(fmt.Sprintf("%f", times[1]))
	file.Close()
}
func main() {

}
