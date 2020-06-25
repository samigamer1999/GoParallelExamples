package main

import(
    "math"
    "time"
    "math/rand"
    "fmt"
    "sync"
    "runtime"
    "os"
    "strconv"
)


var G = 6.67428e-11
var AU = (149.6e6 * 1000)
var SCALE = 500000000000 / AU
var timestep float64


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

type Body struct {
	mass  float64
	acc , vel, pos, force  vector
}

type Node struct {
    body *Body
    centerOfMass vector
    totalMass float64
    children [4]*Node
    rect [4]int
    typee string
}

// Constructor
func newNode(rect [4]int) *Node{
	node := Node{nil, vector{0, 0}, 0, [4]*Node{nil, nil, nil, nil}, [4]int{0, 0, 0, 0}, "leaf"}
    node.rect = rect
    return &node
}

func (n Node) contains(x int, z int) bool{
	    x0, z0, x1, z1 := n.rect[0], n.rect[1], n.rect[2], n.rect[3] 
        if x0 <= x &&  x <= x1 && z0 <= z && z <= z1{
            return true
        }
        return false
}

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

func insertBody(body *Body, node *Node){
    
    if !node.contains(500000000000 + int(body.pos.x * SCALE), 500000000000 + int(body.pos.y * SCALE)){
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

func BodyOnBodyAcc(body1 *Body, body2 *Body){
    temp := G * body1.mass / math.Pow(body1.pos.sub(body2.pos).mod(), 3)
    body2.acc = body2.acc.add(body1.pos.sub(body2.pos).scale(temp))
}


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

func printQuadTree(node *Node){
    fmt.Println("------------")
    fmt.Println(node.rect)
    fmt.Println(node.typee)
    
    if node.body != nil{
    	fmt.Println(node.body.mass)
    }else{
    	fmt.Println("0")
    }
    
    if node.typee == "internal"{
    	for _,child := range node.children{
           printQuadTree(child)
    }
    }
}

func barnespara(size int){
    timestep = 24 * 3600
    runtime.GOMAXPROCS(40)

    n := size
    root := newNode([4]int{0, 0, 1000000000000, 1000000000000})
    // Generate Bodies
    bodies := []Body{}
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    for i := 0; i < n; i++{
    s1 = rand.NewSource(time.Now().UnixNano())
    r1 = rand.New(s1)
    mass := (1 + r1.Float64() * (10000 - 1) )* math.Pow(10, 24)
    vel := -1000 + 2000 * r1.Float64()
    posx := -9 + 18 * r1.Float64()
    posy := -9 + 18 * r1.Float64()
    bodies = append(bodies, Body{mass , vector{0, 0}, vector{0, vel * 1000}, vector{posx * AU, posy * AU}, vector{0, 0}})

    

}

 // Add paralelization here, race conditions aren't a problem because it doesn't matter when/which body we add to the tree
 var wg sync.WaitGroup
extraGoroutines := 1000// On limite le nombre de Goroutines à utiliser

semChan := make(chan int, extraGoroutines)
defer close(semChan)

start := time.Now()
for i := 0; i < len(bodies); i++ {
        insertBody(&bodies[i], root)

    }
 
 
 
    
    for i:= 1; i <= extraGoroutines; i++{
        wg.Add(1)
        go func(i int) {
            for j:= (i - 1) * len(bodies) / extraGoroutines; j < i * len(bodies) / extraGoroutines; j++{
                    bodies[j].acc = vector{0, 0}
                    calculateAcceleration(&bodies[j], root, 0.5)
                    bodies[j].vel = bodies[j].vel.add(bodies[j].acc.scale(timestep))
                    bodies[j].pos = bodies[j].pos.add(bodies[j].vel.scale(timestep))

                }
            wg.Done()
        }(i)
        
}
// Add wait here or either there will be race conditions
 wg.Wait()

 elapsed := time.Now().Sub(start)
 fmt.Println(elapsed)

 f, err := os.OpenFile("test.txt", os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        fmt.Println(err)
        return
    }
    
            _, err = f.WriteString(strconv.Itoa(n)+","+fmt.Sprintf("%f", elapsed.Seconds())+",")
            if err != nil {
                fmt.Println(err)
                f.Close()
                return
            }
            
            
            
            
    
    err = f.Close()
            if err != nil {
                fmt.Println(err)
                return
            }

    start = time.Now()
    root = newNode([4]int{0, 0, 1000000000000, 1000000000000})
    for i := 0; i < len(bodies); i++ {
        insertBody(&bodies[i], root)

    }

    for j:= 0 ; j < len(bodies) ; j++{
                    bodies[j].acc = vector{0, 0}
                    calculateAcceleration(&bodies[j], root, 0.5)
                    bodies[j].vel = bodies[j].vel.add(bodies[j].acc.scale(timestep))
                    bodies[j].pos = bodies[j].pos.add(bodies[j].vel.scale(timestep))

          }

 elapsed = time.Now().Sub(start)
 fmt.Println(elapsed)

 f, err = os.OpenFile("test.txt", os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        fmt.Println(err)
        return
    }
    
            _, err = f.WriteString(fmt.Sprintf("%f", elapsed.Seconds()))
            if err != nil {
                fmt.Println(err)
                f.Close()
                return
            }
            
            _, err = f.WriteString("\n")
            if err != nil {
                fmt.Println(err)
                f.Close()
                return
            }
            
            
    
    err = f.Close()
            if err != nil {
                fmt.Println(err)
                return
            }
}
func main(){
    
  for i := 0; i < 1000000; i = i + 2000 {
    barnespara(i)
  }
// I CHANGED THE FRAME TO 100000X1000000
 /*
  Things to consider:
    - Adding the parameters when executing the file (add file with body data + timestep, possibility to generate random stuff (number of bodies) and also the width/height of the frame)
    - Giving the user the possibility of tuning the extragoroutines parameter

 */
}

