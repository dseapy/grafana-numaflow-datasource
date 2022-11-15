FROM golang:1.19.3 as builder

ARG mage_version=1.14.0

WORKDIR /workspace
RUN wget https://github.com/magefile/mage/releases/download/v${mage_version}/mage_${mage_version}_Linux-64bit.tar.gz && \
    tar -xzvf mage_${mage_version}_Linux-64bit.tar.gz && \
    apt update && \
    curl -sL https://deb.nodesource.com/setup_16.x -o /tmp/nodesource_setup.sh && \
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list && \
    bash /tmp/nodesource_setup.sh && \
    apt update && \
    apt install -y nodejs yarn

# may not need all these...
COPY .eslintrc .eslintrc
COPY .nvmrc .nvmrc
COPY .prettierrc.js .prettierrc.js
COPY jest.config.js jest.config.js
COPY jest-setup.js jest-setup.js
COPY tsconfig.json tsconfig.json
COPY CHANGELOG.md CHANGELOG.md
COPY cypress/ cypress/
COPY .config/ .config/
COPY package.json package.json
COPY yarn.lock yarn.lock
RUN yarn install

COPY go.mod go.mod
COPY go.sum go.sum
COPY Magefile.go Magefile.go
RUN go mod download

COPY pkg/ pkg/
COPY src/ src/

RUN yarn build && ./mage -v build:linux

FROM busybox:latest
WORKDIR /dist
COPY --from=builder /workspace/dist .