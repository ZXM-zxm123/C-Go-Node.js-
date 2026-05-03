FROM gcc:12-alpine

RUN apk add --no-cache cmake make hiredis-dev

WORKDIR /app

COPY . .

RUN mkdir -p build && cd build && cmake .. && make

CMD ["./build/log_receiver"]