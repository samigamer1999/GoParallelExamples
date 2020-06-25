from ctypes import *
from numpy.ctypeslib import ndpointer
import random as rand
import os
import matplotlib.pyplot as plt
import sys

# Load compiled go module
lib = cdll.LoadLibrary("./goquicksort.so")


size=[k for k in range(100000,5000000,1000000)]
liste_para=[]
liste=[]
for n in size:
    
    #on enlève tout les fichiers possibles qui auraient pu rester après un précedent lancement
    if os.path.isfile("./fichier_intermediaire.txt"):
        os.remove("fichier_intermediaire.txt")
    if os.path.isfile("./fichier_intermediaire_temps.txt"):
        os.remove("fichier_intermediaire_temps.txt")
    pyarr = [rand.randint(0,9999999) - rand.randint(0,9999999) for k in range(n)]
    lib.Temps.argtypes = [c_double * len(pyarr)]
    arr = (c_double * len(pyarr))(*pyarr)
    lib.Temps(arr, len(arr))
    fichier=open("fichier_intermediaire_temps.txt")
    lines=fichier.readlines()
    valeur,i="",0
    while lines[0][i]!="S":
        valeur+=lines[0][i]
        i+=1
    liste_para.append(eval(valeur)*1000)
    print("para", valeur, n)
    i+=1
    valeur=""
    while i<len(lines[0]):
        valeur+=lines[0][i]
        i+=1
    liste.append(eval(valeur)*1000)
temp1 = [10, 10010, 20010, 30010, 40010, 50010, 60010, 70010, 80010, 90010]
temp2 = [9.775161743164062e-06 * 1000, 0.021709442138671875* 1000, 0.04332160949707031* 1000, 0.06991887092590332 * 1000, 0.09443998336791992 * 1000, 0.12839031219482422 * 1000, 0.14751768112182617 * 1000, 0.1770946979522705 * 1000 , 0.20438885688781738 * 1000, 0.2458024024963379 * 1000]

os.remove("fichier_intermediaire_temps.txt")
plt.plot(size,liste,color="black",label="Go Sans parallélisation")
plt.plot(size,liste_para,color="red",label="Go Avec parallélisation")

plt.xlabel("Taille de la liste")
plt.ylabel("Temps en millis-seconde")
plt.title("Différence de performance entre script parallélisé ou non avec GOMAXPROCS=40 pour le tri de liste")
plt.legend()
plt.show()
