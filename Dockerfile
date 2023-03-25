FROM cosmtrek/air
WORKDIR /src
COPY . .
CMD ["air", "-c", ".air.toml"]
EXPOSE 80