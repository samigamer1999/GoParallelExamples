package main 

import(
	"C"
	"unsafe"
	"github.com/skelterjohn/go.matrix"
	"sync"
  "runtime"
  "math" 
)

func QRPara(j int, r *matrix.DenseMatrix, q *matrix.DenseMatrix, ch chan string, wg *sync.WaitGroup) {
   ch2 := make(chan string, 1)
   l,_ := r.GetSize()
   index := 0
   state := "halt"

   for i := l - 2; i > j - 1; i--{
     // Si la goroutine précédente a terminé, on peut continuer librement
	   if state != "ended"{
      // On s'arrête ici tant que nous avons pas la permission d'avancer
	   	for state == "halt" {
		   select{
          // On récupère la permission d'avancer de la goroutine précédente
		      case state = <-ch:
		      default:  
          // La channel est vide, on attend encore
		      state = "halt" 
		   }
		}
	   }

    // l'algorithme des rotation de givens typique
    x := r.Arrays()[i][j]
    y := r.Arrays()[i+1][j]
    if y!=0{
    	t := math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
    	cos := x / t
    	sin := y / t
    	rot := matrix.MakeDenseMatrixStacked([][]float64{
					{cos, -sin},
					{sin, cos},
				    })
        SubMatrixTimesLeft(r.GetMatrix(i, j, 2, l - j),rot.Transpose())
        SubMatrixTimesRight(q.GetMatrix(0, i, l, 2),rot)
        index++
        // On ne commence le processus suivant que si notre goroutine a éliminé 2 coefficients 
        if(index == 2 && j < l){
          // La goroutine s'occupera des éléments de la colonne suivante j+1
        	go QRPara(j+1, r, q, ch2, wg)
        }
    }
  
  // On annonce à la goroutine suivante qu'elle peut avancé si on a éliminé un coefficient
  select {
  	case ch2 <- "go":
  	default:
	}  
  }

  // La gouroutine a terminé
  // On vide notre channel
  select {
  	case  <- ch:
  	default:
	}

  // On annonce à la goroutine suivante que nous avons terminé
  select {
  	case ch2 <- "ended":
  	default:
	}
   
   wg.Done()
}



func SubMatrixTimesLeft(SubA *matrix.DenseMatrix, B *matrix.DenseMatrix){
    // Fonction qui multiplie une sous-matrice de A par B à gauche en changeant les valeurs de A  
    temp := matrix.ParallelProduct(B, SubA)
    b, _ := B.GetSize()
    _, a := SubA.GetSize()
    for i := 0 ; i < b ; i++{
      for j := 0 ; j < a ; j++{
         SubA.Arrays()[i][j] = temp.Arrays()[i][j]
    }
    }
}

func SubMatrixTimesRight(SubA *matrix.DenseMatrix, B *matrix.DenseMatrix){
    // Fonction qui multiplie une sous-matrice de A par B à droite en changeant les valeurs de A
    temp := matrix.ParallelProduct(SubA, B)
    _, b := B.GetSize()
    a, _ := SubA.GetSize()
    for i := 0 ; i < a ; i++{
      for j := 0 ; j < b ; j++{
         SubA.Arrays()[i][j] = temp.Arrays()[i][j]
    }
    }
}

//export QR
func QR(cArray *C.double, cSize C.int, rows C.int, cols C.int)  uintptr {
    // On récupère le py array et on le transforme en GoSlice
    gSlice := (*[1 << 30]C.double)(unsafe.Pointer(cArray))[:cSize:cSize]
    result := make([]float64, cSize)
    for i := 0; i < int(cSize); i++ {
        result[i] = float64(gSlice[i])
    }
    // On crée la matrice à l'aide de notre array
    A := matrix.MakeDenseMatrix(result, int(rows), int(cols))
    
    runtime.GOMAXPROCS(4)

    r := matrix.MakeDenseCopy(A)
    l, _ := A.GetSize()
    q := matrix.Eye(l)
    var wg sync.WaitGroup
  	ch := make(chan string, 1)
    // La première goroutine peut avancer librement
  	ch <- "ended"
  	wg.Add(int(rows) - 1)
    go QRPara(0, r, q, ch, &wg)
    wg.Wait()
    // On retourne une concatenation de r et q
    return uintptr(unsafe.Pointer(&append(r.Array(),q.Array()...)[0]))
}

    
func main(){

}