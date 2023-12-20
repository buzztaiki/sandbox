FROM debian:bookworm as builder

RUN apt-get update && apt-get install -y gcc
COPY coredump.c /app/coredump.c
RUN gcc -o /app/coredump /app/coredump.c

FROM debian:bookworm-slim

COPY --from=builder /app/coredump /app/coredump
CMD ["/app/coredump"]
