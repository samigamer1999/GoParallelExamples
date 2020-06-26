from ctypes import *
from numpy.ctypeslib import ndpointer
import pygame
import random
from argparse import ArgumentParser
import time

AU = (149.6e6 * 1000)
SCALE = 50 / AU


def main():
    # On charge le module Go compilé 
    lib = cdll.LoadLibrary("./barneshut.so")

    # On lit les arguments de la ligne de commande
    parser = ArgumentParser()
    parser.add_argument("-x", "--width", dest='width', help="largeur de la fenêtre")
    parser.add_argument("-y", "--height", dest='height', help="hauteur de la fenêtre")
    parser.add_argument("-n", "--number", dest='number', help='Nombre de corps')
    parser.add_argument("-t", "--theta", dest='theta', help='theta')
    parser.add_argument("-f", "--file", dest='file', help='Utiliser les données d\'un fichier', metavar="TEXT")
    args = parser.parse_args()
    # On choisit la largeur et hauteur de notre fenêtre
    width = eval(args.width)
    height = eval(args.height)

    # Nombre de corps
    n = eval(args.number)
    # Le facteur theta
    theta = eval(args.theta)

    # On initilise les paramètres de pygame
    pygame.init()
    screen = pygame.display.set_mode((width, height))
    clock = pygame.time.Clock()

    # On initialise les tableaux de masses, vitesses initiales et positions alétoirement
    if args.file == None:
        mass = [random.randint(1, 10000) * 10 ** 24 for i in range(n)]
        pos = [(-9 + 18 * random.uniform(0, 1) ) * AU for i in range(2*n)]
        vel = [(-10 + 20 * random.uniform(0, 1)) * 1000 for i in range(2*n)]
    else:
        mass =[]
        pos = []
        vel = []
        file = open(args.file, 'r')
        for line in file:
            temp = line.split()
            mass.append(eval(temp[0]))
            pos.append(eval(temp[1])) #posx
            pos.append(eval(temp[2])) #posy
            vel.append(eval(temp[3])) #velx
            vel.append(eval(temp[4])) #vely
    while True:
        # Controler le pas de temps 
        slider = 1
        while True:
            timestep = slider * 24 * 3600
            events = pygame.event.get()

            # Listener modifier la valeur du slider (entre 1 et 20 jours)
            for event in events:
                if event.type == pygame.KEYDOWN:
                    if event.key == pygame.K_LEFT:
                        if slider > 1:
                            slider -= 1
                    if event.key == pygame.K_RIGHT:
                        if slider < 20:
                            slider += 1

            for e in events:
                if e.type == pygame.QUIT:
                    return

            # On précise les types des paramètres en entrées et sortie de la fonction CalcPositions sur Go 
            lib.CalcPositions.argtypes = [c_double * len(pos)] + [c_double * len(vel)] + [c_double * len(mass)] 
            lib.CalcPositions.restype = ndpointer(dtype = c_double, shape = (4*len(mass),))

            # Crée des C arrays à la base de py arrays 
            pos = (c_double * len(pos))(*pos)
            vel = (c_double * len(vel))(*vel)
            mass = (c_double * len(mass))(*mass)
            
            t1 = time.clock()
            
            # Retourne une liste 1D , concatenation de postition et vitesses des corps
            posandvels = lib.CalcPositions(pos, vel, mass, len(mass), timestep, width, height, int(theta * 100))

            # On met à jour les positions et les vitesses
            pos = posandvels[:2*len(mass)]
            vel = posandvels[2*len(mass):]
            t2 = time.clock() - t1
            print("Temps de mise-à-jour :", t2)
            
            # Le traçage
            screen.fill((30, 30, 30))
            for i in range(len(mass)):
                # Les positions sont normalisées afin qu'elle rentre dans le petit cadre de la fenêtre
                pygame.draw.circle(screen, (255,255,255) , (int(width / 2) + int(pos[2 * i] * SCALE), int(height / 2) + int(pos[2 * i + 1] * SCALE)), 5)
            pygame.display.update()

            dt = clock.tick(60)


if __name__ == "__main__":
    main()
