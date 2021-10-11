import shutil
from PIL import Image
import os

path = './datacollection/data/images/'

speciesDeleted = 0
corruptedRemoved = 0

# get rid of corrupted images
for root, dirs, files in os.walk(path):
    for dir in dirs:
        for file in os.listdir(path+dir):
            print(path + dir)
            try:
                img = Image.open(path+dir+'/'+file)
            except OSError:
                os.remove(path+dir+'/'+file)
                corruptedRemoved = corruptedRemoved + 1

# get rid of species with less than 10 photos
for root, dirs, files in os.walk(path):
    for dir in dirs:
        length = len(os.listdir(path+dir))
        if length < 10:
            speciesDeleted = speciesDeleted + 1
            shutil.rmtree(path+dir)            

print("Number of species removed: " + str(speciesDeleted))
print("Number of corrupted files removed: " + str(corruptedRemoved))