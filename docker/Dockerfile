FROM amd64/debian:stable-slim
ENV FLOW_SERVER_PORT="80"
WORKDIR /app
RUN export DEBIAN_FRONTEND=noninteractive && apt-get update && apt-get install -yq \
  tzdata \
  curl \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*
ADD flow-app /app/
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s CMD curl --fail http://localhost:$FLOW_SERVER_PORT/health || exit 1
EXPOSE 80
ENTRYPOINT ["./flow-app"]
