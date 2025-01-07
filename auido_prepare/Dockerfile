FROM python:3.10

WORKDIR /usr/src/app
COPY ./src .

# RUN pip install --upgrade pip
RUN pip install -r requirements.txt
# RUN python -m grpc_tools.protoc --python_out=. --grpc_python_out=. --pyi_out=. --proto_path=. ./grpc_core/audio_preparer.proto

EXPOSE 7001
CMD [ "python", "./main.py", "7001"]