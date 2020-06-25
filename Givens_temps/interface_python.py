from ctypes import *
import numpy as np
import random as rd
from numpy.ctypeslib import ndpointer
import matplotlib.pyplot as plt
import os
import time
lib = cdll.LoadLibrary("./gogivens.so")


liste_temps_interfacage1=[]
liste_temps_interfacage2=[]
size=[k for k in range(10,100,5)]
for n in size:
    #on enlève tout les fichiers possibles qui auraient pu rester après un précedent lancement
    if os.path.isfile("./fichier_intermediaire.txt"):
        os.remove("fichier_intermediaire.txt")
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

    # Make C array from py array
    start=time.time()
    arr = (c_double * len(pyarr))(*pyarr)
    qr_array = lib.Interfacage2(arr, len(arr), rows, cols)
    r = qr_array[0:rows * cols].reshape(rows, cols)
    q = qr_array[rows * cols:].reshape(rows, cols)
    end=time.time()-start
    liste_temps_interfacage2.append(end)


    start=time.time()
    fichier=open("fichier_intermediaire.txt","x")
    taille=str(len(matrix))[::-1]
    fichier.write(taille+"S"+"\n")
    for liste in matrix:
        for valeur in liste:
            fichier.write(str(valeur))
            fichier.write("S")
        fichier.write("\n")
    fichier.close()
    lib.Interfacage1()
    fichier_R,fichier_Q=open("fichier_intermediaire_R.txt","r"),open("fichier_intermediaire_Q.txt","r")
    lines_R,lines_Q=fichier_R.readlines(),fichier_Q.readlines()
    taille=len(matrix)
    R,Q=np.eye(taille),np.eye(taille)
    for line in range(taille):
        q,i = 0, 0
        while i < len(lines_R[line])-1:
            float_string=str()
            while str(lines_R[line][i])!="S" :
                float_string += str(lines_R[line][i])
                i+=1
            R[line][q]= float(float_string)
            q+=1
            i+=1
        q,i=0,0
        while i < len(lines_Q[line])-1:
            float_string=str()
            while str(lines_Q[line][i])!="S" :
                float_string += str(lines_Q[line][i])
                i+=1
            Q[line][q]= float(float_string)
            q+=1
            i+=1
    fichier_R.close()
    fichier_Q.close()
    end=time.time()-start
    os.remove("fichier_intermediaire_Q.txt")
    os.remove("fichier_intermediaire_R.txt")
    liste_temps_interfacage1.append(end)
plt.plot(size,liste_temps_interfacage1,color="black",label="Interfaçage avec fichier")
plt.plot(size,liste_temps_interfacage2,color="red",label="Interfaçage avec pointeurs")
plt.xlabel("Taille de la matrice")
plt.ylabel("Temps en seconde")
plt.title("Différence de performance des interfaçages sur la décomposition QR sans parallélisation ")
plt.legend()
plt.show()
