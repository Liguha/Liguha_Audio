import torch
import librosa
import json
import numpy as np
import torch.nn as nn

import resampy
from torchvision import models

SR = 16000
N_MELS = 128
MAX_TIME_FRAMES = 3000
BATCH_SIZE = 1
GENRE_TO_ID_PATH = "./files/genre_to_id.json"
ID_TO_GENRE_PATH = "./files/id_to_genre.json"


class VGGishClassifier(nn.Module):
    def __init__(self, num_classes):
        super(VGGishClassifier, self).__init__()
        self.vggish = models.vgg16(pretrained=True)
        self.vggish.classifier[6] = nn.Linear(4096, num_classes)

    def forward(self, x, return_embeddings=False):
        x = self.vggish.features(x)
        x = self.vggish.avgpool(x)
        x = torch.flatten(x, 1)
        embeddings = self.vggish.classifier[:-1](x)
        if return_embeddings:
            return embeddings, self.vggish.classifier[-1](embeddings)
        return self.vggish.classifier[-1](embeddings)


def compute_mel_spectrogram(audio, sr=SR, n_fft=1024, hop_length=512, n_mels=N_MELS):
    mel_spec = librosa.feature.melspectrogram(
        y=audio, sr=sr, n_fft=n_fft, hop_length=hop_length, n_mels=n_mels
    )
    mel_spec_db = librosa.power_to_db(mel_spec, ref=np.max)
    return mel_spec_db

def pad_spectrogram(mel_spec, max_frames=MAX_TIME_FRAMES):
    if mel_spec.shape[1] < max_frames:
        pad_width = max_frames - mel_spec.shape[1]
        mel_spec = np.pad(mel_spec, ((0, 0), (0, pad_width)), mode='constant')
    else:
        mel_spec = mel_spec[:, :max_frames]
    return mel_spec

def vectorize(model: VGGishClassifier, fs: int, data: np.ndarray):
    audio = resampy.resample(data, fs, SR)
    mel_spec = compute_mel_spectrogram(audio)
    mel_spec = pad_spectrogram(mel_spec)
    mel_spec = np.expand_dims(mel_spec, axis=0)
    mel_spec = np.repeat(mel_spec, 3, axis=0)
    mel_spec_tensor = torch.tensor(mel_spec, dtype=torch.float32).unsqueeze(0).to(device)
    with torch.no_grad():
        embeddings, outputs = model(mel_spec_tensor, return_embeddings=True)
        _, predicted = torch.max(outputs, 1)
    genre = id_to_genre[str(predicted.item())]
    embedding_vector = embeddings.cpu().numpy().flatten()
    return genre, embedding_vector


with open(ID_TO_GENRE_PATH, "r") as f:
    id_to_genre = json.load(f)
device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')