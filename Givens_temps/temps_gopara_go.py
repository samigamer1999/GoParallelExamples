from ctypes import *
import numpy as np
import random as rd
from numpy.ctypeslib import ndpointer
import matplotlib.pyplot as plt
import os
import time
lib = cdll.LoadLibrary("./gogivens.so")


liste_para=[]
liste=[]
size=[k for k in range(10,150,4)]
for n in size:
    #on enlève tout les fichiers possibles qui auraient pu rester après un précedent lancement
    if os.path.isfile("./fichier_intermediaire.txt"):
        os.remove("fichier_intermediaire.txt")
    if os.path.isfile("./fichier_intermediaire_temps.txt"):
        os.remove("fichier_intermediaire_temps.txt")
    if os.path.isfile("./fichier_intermediaire_R.txt"):
        os.remove("fichier_intermediaire_R.txt")
    if os.path.isfile("./fichier_intermediaire_Q.txt"):
        os.remove("fichier_intermediaire_Q.txt")
    matrix=np.eye(n)
    pyarr=[]
    for ligne in range(n):
        for colonne in range(n):
            valeur=rd.randint(0,999)
            pyarr.append(valeur)
            matrix[ligne][colonne]=valeur
    rows = n
    cols = n
    lib.Interfacage2.argtypes = [c_double * len(pyarr)]
    lib.Interfacage2.restype = ndpointer(dtype = c_double, shape = (2*len(pyarr),))
    lib.Temps.argtypes = [c_double * len(pyarr)]

    # Make C array from py array
    arr = (c_double * len(pyarr))(*pyarr)
    lib.Temps(arr, len(arr), rows, cols)
    fichier=open("fichier_intermediaire_temps.txt")
    lines=fichier.readlines()
    valeur,i="",0
    while lines[0][i]!="S":
        valeur+=lines[0][i]
        i+=1
    liste_para.append(eval(valeur))
    i+=1
    valeur=""
    while i<len(lines[0]):
        valeur+=lines[0][i]
        i+=1
    liste.append(eval(valeur))
os.remove("fichier_intermediaire_temps.txt")
plt.plot(size,liste,color="black",label="Sans parallélisation")
plt.plot(size,liste_para,color="red",label="Avec parallélisation")
plt.xlabel("Taille de la matrice")
plt.ylabel("Temps en secondes")
plt.title("Différence de performance entre script parallélisé ou non avec GOMAXPROCS=4 pour la rotation de givens")
plt.legend()
plt.show()
