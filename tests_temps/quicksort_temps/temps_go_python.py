from ctypes import *
import matplotlib.pyplot as plt
import numpy as np
import math
import random as rand
import time



def quick_sort(array, start, end):
    if start >= end:
        return

    p = partition(array, start, end)
    quick_sort(array, start, p-1)
    quick_sort(array, p+1, end)

def partition(array, start, end):
    pivot = array[start]
    low = start + 1
    high = end

    while True:
        # If the current value we're looking at is larger than the pivot
        # it's in the right place (right side of pivot) and we can move left,
        # to the next element.
        # We also need to make sure we haven't surpassed the low pointer, since that
        # indicates we have already moved all the elements to their correct side of the pivot
        while low <= high and array[high] >= pivot:
            high = high - 1

        # Opposite process of the one above
        while low <= high and array[low] <= pivot:
            low = low + 1

        # We either found a value for both high and low that is out of order
        # or low is higher than high, in which case we exit the loop
        if low <= high:
            array[low], array[high] = array[high], array[low]
            # The loop continues
        else:
            # We exit out of the loop
            break

    array[start], array[high] = array[high], array[start]

    return high

liste_temps=[]
size=[i for i in range(10,100000,10000)]
for i in size:
    liste=[]
    for k in range(i):
        liste.append(rand.randint(0,9999999) - rand.randint(0,9999999))
    start=time.time()
    quick_sort(liste,0,i-1)
    end=time.time() - start
    liste_temps.append(end)

print(size)
print(liste_temps)
plt.plot(size,liste_temps,color="black",label="temps en secondes")
plt.legend()
plt.show()
