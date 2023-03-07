FROM python:3.9-slim

COPY ./app /app/src
COPY ./requirements.txt /app
COPY ./config.toml /app

WORKDIR /app

RUN pip3 install -r requirements.txt

CMD ["python3", "src/main.py"]