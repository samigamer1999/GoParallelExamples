from ctypes import *
from numpy.ctypeslib import ndpointer
import time

# On charge le module go compilé
lib = cdll.LoadLibrary("./quicksort.so")


pyarr = [100,101,3,5,6,54,84,84,84,8,456,4,564,56,5]  # La liste à trier

# On précise les types des paramètres en entrées et sortie de la fonction Sort sur Go 
lib.Sort.argtypes = [c_double * len(pyarr)]
lib.Sort.restype = ndpointer(dtype = c_double, shape = (len(pyarr),))

# Crée un C array à la base d'un py array 
arr = (c_double * len(pyarr))(*pyarr)

t1 = time.clock()
# On calcule le résultat
result = lib.Sort(arr, len(arr))
t2 = time.clock() - t1
print("Temps d'éxection :", t2)

print(result)
