FROM ghcr.io/cirruslabs/flutter:stable AS builder

WORKDIR /app

USER root
RUN apt-get update && apt-get install -y libglu1-mesa

COPY pubspec.yaml pubspec.lock ./
RUN flutter pub get

COPY . .
# Add a dummy build argument to break cache for the build step
ARG BUILD_VERSION=unknown
RUN flutter build web --release

FROM nginx:alpine

COPY --from=builder /app/build/web /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
