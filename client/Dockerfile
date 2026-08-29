FROM ghcr.io/cirruslabs/flutter:stable AS builder

WORKDIR /app

USER root
RUN apt-get update && apt-get install -y libglu1-mesa

COPY pubspec.yaml pubspec.lock ./
RUN --mount=type=cache,target=/root/.pub-cache \
    flutter pub get

COPY . .
# Add a dummy build argument to break cache for the build step
ARG BUILD_VERSION=unknown
RUN --mount=type=cache,target=/root/.pub-cache \
    --mount=type=cache,target=/app/.dart_tool \
    --mount=type=cache,target=/app/build \
    flutter build web --release && \
    mkdir -p /app/build_output && \
    cp -r build/web/. /app/build_output/

FROM nginx:alpine

COPY --from=builder /app/build_output /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
