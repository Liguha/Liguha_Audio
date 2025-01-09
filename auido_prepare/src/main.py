import boto3
import torch
import psycopg2
import json
from flask import Flask, request, jsonify

from utils.audio_helper import bytes2audio, make_segments
from utils.vectorizer import VGGishClassifier, vectorize

db_creds = json.load(open("./files/db_creds.json"))
conn = psycopg2.connect(dbname=db_creds["db_name"], 
                        user=db_creds["user"], 
                        password=db_creds["password"], 
                        host=db_creds["host"], 
                        port=db_creds["port"])

session = boto3.session.Session(profile_name='default')
s3 = session.client(
   service_name='s3',
   endpoint_url='https://s3.cloud.ru'
)
bucket_name = "audiobucket"

MODEL_PATH = "./files/vggish_classifier_full.pth"
device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
model: VGGishClassifier = torch.load(MODEL_PATH, map_location=device)
model = model.to(device)


app = Flask("audio_processing")

@app.route('/find_by_segment', methods=['POST'])
def find_by_segment():
    s3_id = request.get_json()["s3_id"]
    response = s3.get_object(Bucket=bucket_name, Key=s3_id)
    data = response['Body'].read()
    fs, channels = bytes2audio(data)
    segs = make_segments(channels[0])
    if len(segs) > 0:
        _, vec = vectorize(model, fs, segs[0])
        vec = list(vec)
        cur = conn.cursor()
        cur.execute(f"SELECT song_id, embedding <-> '{vec}' AS distance\
                    FROM music.vectors\
                    ORDER BY distance\
                    LIMIT 1;")
        id, _ = cur.fetchall()[0]
        return jsonify({"status": "ok", "id": id})
    return jsonify({"status": "Error. Record too short"})

@app.route('/add_audio', methods=['POST'])
def add_audio():
    cur = conn.cursor()
    id = request.get_json()["id"]
    response = s3.get_object(Bucket=bucket_name, Key=f"Song_{id}")
    data = response['Body'].read()
    fs, channels = bytes2audio(data)
    x = channels[0]
    # tag, _ = vectorize(model, fs, x)
    segs = make_segments(fs, x)
    for seg in segs:
        _, vec = vectorize(model, fs, seg)
        vec = list(vec)
        cur.execute(f"INSERT INTO music.vectors (embedding, song_id) VALUES ('{vec}', {id})")
    conn.commit()

app.run(port=7001)