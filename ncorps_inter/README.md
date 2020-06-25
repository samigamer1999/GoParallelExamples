Utiliser notre programme de résolution de problèmes à N corps via Python
------------------------------------

Pour utiliser notre programme de résolution à N corps il va falloir que vous importiez la librairie **pygame** sur Python 3.  
Pour cela veuillez taper sur votre terminal : `pip3 install pygame`

Pour éxecuter le programme, vous devez lui donner les arguments nécessaires dans la ligne de commande. Voici un exemple :  
`python3 barneshut.py -x 1000 -y 1000 -n 100 -t 0.5`

Ceci permettra de créer une fenêtre de résolution 1000x1000 avec 100 corps aléatoires et un theta = 0.5.

Afin d'utiliser vos propres données des corps, utilisez `-f votrefichier.txt`. Le fichier devra contenir (masse, posx, posy, velx, vely) dans cet ordre pour chaque ligne (chaque corps) avec un espace entre chaque valeur.

Pour contrôler le pas de temps, utilisez les flêches gauche/droite pour diminuer/augmenter sa valeur.

En plus d’avoir une sortie graphique, vous avez la sortie de toute les données nécessaires
avec **pos** (liste des positions) et **vels** (liste des vitesses) (liste des   (ligne 88, 89 du script).
