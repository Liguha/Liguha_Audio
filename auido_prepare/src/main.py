import grpc
import sys
import numpy as np
from concurrent import futures
from utils.vectorizer import vectorize_audio, make_segments
from grpc_core import audio_preparer_pb2, audio_preparer_pb2_grpc

    
class AudioPreparer(audio_preparer_pb2_grpc.AudioPreparer):
    def find_audio(self, request: audio_preparer_pb2.Audio, 
                   context: grpc.aio.ServicerContext) -> audio_preparer_pb2.AudioResponse:
        fs = request.sample_rate
        data = np.array(request.data).astype(float)
        segments = make_segments(fs, data)

        ok: bool = True
        err: str = "No errors"
        id: str = "Undefined"
        if len(segments) < 1:
            ok = False
            err = "Segment too short"
        else:
            embedding = vectorize_audio(*segments[0])
            # some magic with DB
            id = "aboba"
        return audio_preparer_pb2.AudioResponse(ok=ok, err=err, id=id)

    def add_audio(self, request: audio_preparer_pb2.Audio, 
                  context: grpc.aio.ServicerContext) -> audio_preparer_pb2.AudioResponse:
        fs = request.sample_rate
        data = np.array(request.data).astype(float)
        segments = make_segments(fs, data)
        # generate unique id with help of DB
        ok: bool = False
        err: str = "Audio is too short"
        id: str = "aboba_new"
        for seg in segments:
            embedding = vectorize_audio(*seg)
            # push to DB with id
            ok = True
            err = "No errors"
        return audio_preparer_pb2.AudioResponse(ok=ok, err=err, )


port = sys.argv[1]
server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
audio_preparer_pb2_grpc.add_AudioPreparerServicer_to_server(AudioPreparer(), server)
server.add_insecure_port("[::]:" + port)
server.start()
print("Server started, listening on " + port)
server.wait_for_termination()