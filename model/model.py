from fastai import *
from fastai.vision import *
from fastai.callback import *
from fastai.metrics import error_rate
from PIL import Image
# from fastai.utils import *
# show_install()

path = r'../datacollection/data/images/'
data =  ImageDataBunch.from_folder(path, train=".", valid_pct=0.1,
        ds_tfms=get_transforms(), size=224, num_workers=4).normalize(imagenet_stats)
learn = cnn_learner(data, models.resnet34, metrics=error_rate)

img = Image.open('../datacollection/data/input/chanty.png')

is_chanty,_,probs = learn.predict(img)
print(f"Is this a chanty?: {is_chanty}.")
print(f"Probability it's a cat: {probs[1].item():.6f}")