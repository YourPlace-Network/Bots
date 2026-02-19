FROM public.ecr.aws/ubuntu/ubuntu:24.04 AS builder

RUN apt-get update && apt-get install -y \
    wget ca-certificates make \
    && rm -rf /var/lib/apt/lists/*

RUN wget https://go.dev/dl/go1.25.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.25.4.linux-amd64.tar.gz && \
    rm go1.25.4.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /build

COPY . .

RUN make build

FROM public.ecr.aws/ubuntu/ubuntu:24.04

RUN apt-get update && apt-get install -y \
    ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /build/newsbot .
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

VOLUME ["/app/data"]

ENTRYPOINT ["./entrypoint.sh"]
