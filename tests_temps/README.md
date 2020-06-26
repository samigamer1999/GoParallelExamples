Les scripts/programmes de tests des temps d'éxecutions
----------------------------------------------------------------

Pour faire un bencmarking (avec des graphiques) des temps d'éxecutions de nos programmes sur votre propre machine, vous devez changer 
aller sur le fichier .go du programme qui vous intérresse, et renseigner le nombre de cpus de votre machine / cpus que vous souhaitez utiliser
en cherchant l'instruction `runtime.GOMAXPROCS(40)` et en la changeant par `runtime.GOMAXPROCS(nombre de cpus)`.

Éxecutez les scripts Python de chaque algorithme pour avoir vos résultats sous formes de graphiques `matplotlib.pyplot`.

PS: Pour tester l'algorithme de Barnes-Hut vous devez éxecutez le fichier .go comme suit : `go run barneshut.go`. Vous pourrez changer les multiples variables (largeur, hauteur, nombre de corps, theta) dans le fichier directement. 
