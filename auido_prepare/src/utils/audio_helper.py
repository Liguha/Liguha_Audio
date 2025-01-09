import io
import numpy as np
from pydub import AudioSegment

SEG_LEN = 10
OVERLAP = 0.5

def bytes2audio(x: bytes) -> tuple[int, list[np.ndarray]]:
    """Returns sample rate and audio channels."""
    s = io.BytesIO(x)
    audio = AudioSegment.from_file(s)
    fs = audio.frame_rate
    samples = audio.get_array_of_samples()
    n_ch = audio.channels
    width = audio.sample_width * 8
    channels = [np.array(samples[i::n_ch]).astype(float) / (2 ** (width - 1)) for i in range(n_ch)]
    return fs, channels

def make_segments(fs: int, audio: np.ndarray) -> list[tuple[int, np.ndarray]]:
    n_samples = fs * SEG_LEN
    shift = int(n_samples * OVERLAP)
    n = len(audio) - n_samples
    return [(fs, audio[i:(i + n_samples)]) for i in range(0, n, shift)]