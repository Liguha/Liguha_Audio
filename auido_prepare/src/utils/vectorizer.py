from scipy import signal
import numpy as np

seg_len: int = 10       # length of vectorized segment in seconds
overlap: float = 0.5    # overlap between segments in seg_len
lims: list[int] = [20, 88, 350, 1350, 2794, 20000]

def vectorize_audio(fs: int, audio: np.ndarray) -> np.ndarray:
    """Use sample rate and audio samples to create embedding."""
    freqs, times, fft = signal.stft(audio, fs, nperseg=2048, noverlap=0)
    fft = np.abs(fft).T
    dots = []
    for i in range(len(times)):
        global_power = 0.3 * (np.max(fft[i]) + np.mean(fft[i]))
        for j in range(1, len(lims)):
            lhs = np.searchsorted(freqs, lims[j - 1], "left")
            rhs = np.searchsorted(freqs, lims[j], "right")
            local_power = 0.5 * (np.max(fft[i][lhs:rhs]) + np.mean(fft[i][lhs:rhs]))
            idx = np.argmax(fft[i][lhs:rhs]) + lhs
            if fft[i][idx] > max(local_power, global_power):
                dots.append((times[i], freqs[idx]))
    dots.sort()
    # embedding = ...(dots) - NN usage
    return np.array([audio[0], audio[1], audio[2]])

def make_segments(fs: int, audio: np.ndarray) -> list[tuple[int, np.ndarray]]:
    n_samples = fs * seg_len
    shift = int(n_samples * overlap)
    n = len(audio) - n_samples
    return [(fs, audio[i:(i + n_samples)]) for i in range(0, n, shift)]