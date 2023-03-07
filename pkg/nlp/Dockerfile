FROM python:3.9-slim

COPY ./app /app/src
COPY ./requirements.txt /app
COPY ./config.toml /app
COPY ./model.pckl /app

WORKDIR /app

RUN --mount=type=cache,target=/root/.cache/pip pip3 install -r requirements.txt

# RUN pip3 install -r requirements.txt

CMD ["python3", "src/main.py"]