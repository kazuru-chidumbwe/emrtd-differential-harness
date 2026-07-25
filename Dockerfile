# Bounded reproduction image for the eMRTD differential harness.
# Goal: one buildable artifact reviewers can run without a hand-tuned host.
#
#   docker build -t emrtd-harness .
#   docker run --rm emrtd-harness                 # make smoke (TC-AC-01)
#   docker run --rm emrtd-harness bash scripts/run_offline_pa.sh
#
# Requires network at build time (gmrtd clone + Maven Central).
# Synthetic profiles only — no NFC / physical documents.

FROM eclipse-temurin:17-jdk-jammy

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    git \
    make \
    maven \
    python3 \
    && rm -rf /var/lib/apt/lists/*

ARG GO_VERSION=1.25.0
RUN curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" \
    | tar -C /usr/local -xz
ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"
ENV GOTOOLCHAIN=auto
ENV GOPROXY=https://proxy.golang.org,direct

# go.mod replace path expects sibling ../_vendor/gmrtd relative to this tree.
WORKDIR /workspace/emrtd-differential-harness
COPY . .

RUN bash scripts/bootstrap-vendor.sh \
    && bash scripts/install-jmrtd-local.sh \
    && (cd drivers/jmrtd && mvn -q -DskipTests package) \
    && go test ./classifier/... ./middleware/...

CMD ["make", "smoke"]
