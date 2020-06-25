package main

import(
	"C"
	"unsafe"
	"math"
    "sync"
    "runtime"
)


var G = 6.67428e-11
var AU = (149.6e6 * 1000)
var SCALE = 5000 / AU


// La "classe" vecteur
type vector struct{ x, y float64 }
 
func (v vector) add(w vector) vector {
    return vector{v.x + w.x, v.y + w.y}
}
 
func (v vector) sub(w vector) vector {
    return vector{v.x - w.x, v.y - w.y}
}
 
func (v vector) scale(m float64) vector {
    return vector{v.x * m, v.y * m}
}
 
func (v vector) mod() float64 {
    return math.Sqrt(v.x*v.x + v.y*v.y)
}

// La "classe" body
type Body struct {
	mass  float64
	acc , vel, pos  vector
}

// La classe "Node"
type Node struct {
    body *Body
    centerOfMass vector
    totalMass float64
    children [4]*Node
    rect [4]int
    typee string
}

// Constructeur de Node
func newNode(rect [4]int) *Node{
	node := Node{nil, vector{0, 0}, 0, [4]*Node{nil, nil, nil, nil}, [4]int{0, 0, 0, 0}, "leaf"}
    node.rect = rect
    return &node
}

// Fonction qui retourne vrai si la (x,z) est dans le noeud
func (n Node) contains(x int, z int) bool{
	    x0, z0, x1, z1 := n.rect[0], n.rect[1], n.rect[2], n.rect[3] 
        if x0 <= x &&  x <= x1 && z0 <= z && z <= z1{
            return true
        }
        return false
}

// Fonction qui divise un noeuds en quatre
func subdivide(n *Node) {
	n.typee = "internal"
    x0, z0, x1, z1 := n.rect[0], n.rect[1], n.rect[2], n.rect[3] 
    h := (x1 - x0) / 2
    rects := [][4]int{}
    rects = append(rects, [4]int{x0, z0, x0 + h, z0 + h}) 
    rects = append(rects, [4]int{x0, z0 + h, x0 + h, z1}) 
    rects = append(rects, [4]int{x0 + h, z0 + h, x1, z1}) 
    rects = append(rects, [4]int{x0 + h, z0, x1, z0 + h}) 
    for i:=0 ; i < len(rects); i++{
        n.children[i] = newNode(rects[i])
    }
}

func updateCenterOfMass(n *Node){
	if n.typee == "leaf"{
		if n.body != nil{
			n.centerOfMass = n.body.pos
			n.totalMass = n.body.mass
		}
	}else{
		    n.centerOfMass = vector{0, 0}
			n.totalMass = 0
			for i := range n.children{
				updateCenterOfMass(n.children[i])
				n.totalMass += n.children[i].totalMass
                n.centerOfMass =  n.centerOfMass.add(n.children[i].centerOfMass.scale(n.children[i].totalMass))
			}
			n.centerOfMass = n.centerOfMass.scale(1 / n.totalMass)
	}
}


// Insertion de corps dans l'arbre récursivement
func insertBody(body *Body, node *Node){
    
    if !node.contains(5000 + int(body.pos.x * SCALE), 5000 + int(body.pos.y * SCALE)){
        return
    }
     if node.body == nil{
     	
    	node.body = body
        return
    }
     if node.typee == "internal"{
       for _,child := range node.children{
            insertBody(body, child)
       }
       updateCenterOfMass(node)
    }
    if node.typee == "leaf" && node.body != nil{
    	subdivide(node)
        for i :=  range node.children{
            insertBody(node.body, node.children[i])
        }
        for i :=  range node.children{
            insertBody(body, node.children[i])
        }
}
}

// L'accélération -> d'un corps sur un autre
func BodyOnBodyAcc(body1 *Body, body2 *Body){
    temp := G * body1.mass / math.Pow(body1.pos.sub(body2.pos).mod(), 3)
    body2.acc = body2.acc.add(body1.pos.sub(body2.pos).scale(temp))
}

// L'accélération -> d'un noeud sur un corps
func nodeOnBodyAcc(mass float64, pos vector, body *Body){
    temp := G * mass / math.Pow(pos.sub(body.pos).mod(), 3)
    body.acc = body.acc.add(pos.sub(body.pos).scale(temp))

}
func calculateAcceleration(body *Body, node *Node, theta float64){
    if node.typee == "leaf"{
    	 if node.body == nil{
    	 	return
    	 }else if node.body != body{
            BodyOnBodyAcc(node.body, body)
        }
    }else{
    	// Compare s/d with theta, s : width of the region, d : distance between body and center of mass of the node
        s := float64(node.rect[2]  - node.rect[0]) / SCALE
        d := node.centerOfMass.sub(body.pos).mod()
        if s / d < theta{
            nodeOnBodyAcc(node.totalMass, node.centerOfMass, body)
        }else{
            for i := range node.children{
                calculateAcceleration(body, node.children[i], theta)
            }
        }
    }
}


//export CalcPositions
func CalcPositions(posArray *C.double, velArray *C.double, massArray *C.double, Size C.int, timestep C.int, width C.int, height C.int, theta C.int) uintptr{

    // On récupère toutes les listes python, et on les transforme en GoSlices
	posSlice := (*[1 << 30]C.double)(unsafe.Pointer(posArray))[:2 * Size: 2 * Size]
    pos := make([]float64, 2 * Size)
    for i := 0; i < 2 * int(Size); i++ {
        pos[i] = float64(posSlice[i])
    }

    velSlice := (*[1 << 30]C.double)(unsafe.Pointer(velArray))[:2 * Size: 2 * Size]
    vel := make([]float64, 2 * Size)
    for i := 0; i < 2 * int(Size); i++ {
        vel[i] = float64(velSlice[i])
    }

    massSlice := (*[1 << 30]C.double)(unsafe.Pointer(massArray))[:Size:Size]
    mass := make([]float64, Size)
    for i := 0; i < int(Size); i++ {
        mass[i] = float64(massSlice[i])
    }

    
    // On génére des objets Body à l'aide des données
    bodies := []Body{}

    for i := 0; i < int(Size); i++{
	    bmass := mass[i]
	    bvelx := vel[2 * i]
	    bvely := vel[2 * i + 1]
	    bposx := pos[2 * i]
	    bposy := pos[2 * i + 1]
	    bodies = append(bodies, Body{bmass , vector{0, 0}, vector{bvelx, bvely}, vector{bposx , bposy}})
     }

    runtime.GOMAXPROCS(4)
    var wg sync.WaitGroup

    // Le noeud père
    root := newNode([4]int{0, 0, int(width), int(height)})

    // Insertion de tous les coprs dans l'arbre
	for i:= 0; i < len(bodies); i++{
	        insertBody(&bodies[i], root)   
	 } 
	 
	extraGoroutines := 100
	semChan := make(chan int, extraGoroutines)
	defer close(semChan)
	
    // Calcul des positions
 	for i:= 1; i <= extraGoroutines; i++{
       wg.Add(1)
	   go func(i int) {
		    for j:= (i - 1) * len(bodies) / extraGoroutines; j < i * len(bodies) / extraGoroutines; j++{
		            bodies[j].acc = vector{0, 0}
		            calculateAcceleration(&bodies[j], root, float64(theta)/100)
		            bodies[j].vel = bodies[j].vel.add(bodies[j].acc.scale(float64(timestep)))
		            bodies[j].pos = bodies[j].pos.add(bodies[j].vel.scale(float64(timestep)))

		        }
		            wg.Done()
		        }(i) 
	// Add wait here or either there will be race conditions
	 wg.Wait()
	}

    // La liste 1D des positions et vitesses
	posandvel := []float64{}
	for i := 0; i < int(Size); i++{
        posandvel = append(posandvel, bodies[i].pos.x)
        posandvel = append(posandvel, bodies[i].pos.y)
	}
	for i := 0; i < int(Size); i++{
        posandvel = append(posandvel, bodies[i].vel.x)
        posandvel = append(posandvel, bodies[i].vel.y)
	}
	
	return uintptr(unsafe.Pointer(&posandvel[0]))
}
func main(){}