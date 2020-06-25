Utiliser notre programme de rotation de Givens via Python
----------------------------------------------------------------

Pour utiliser notre programme de rotation de Givens, il va falloir que vous importiez de Github une librairie : la librairie matrix.go (disponible sur [https://github.com/skelterjohn/go.matrix](https://github.com/skelterjohn/go.matrix).

Pour cela veuillez vous placer avec votre terminal sur l’espace de travail où se situe le script Go et le script Python.  
Veuillez tapez :  `go get "github.com/skelterjohn/go.matrix"`

C’est bon la librairie est installée.

Maintenant sur le script Python il faut renseigner la matrice que vous voulez décomposer (ligne 8 sur le script) derrière la variable **matrix**.

A la fin de la fonction vous avez deux matrices **r,q** de type numpy qui sont le résultat de la décomposition de **matrix** (lignes 29 et 30 sur le script).

	
