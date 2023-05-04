# Supports ARM
# firestore.tools doesn't support ARM: https://github.com/firebase/firebase-tools/issues/2341
FROM eclipse-temurin:17.0.7_7-jre-jammy

# latest firebase-tools (v11 and lator) require nodejs 14+.
# jammy's nodejs is 12, and install firebase-tools v10
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
    nodejs \
    npm \
  && apt-get -y clean \
  && rm -rf /var/lib/apt/lists/* \
  && npm install -g firebase-tools@10 \
  && firebase setup:emulators:firestore \
  && firebase setup:emulators:ui \
  && mkdir /data

COPY entrypoint.sh firebase.json /

EXPOSE 4000
EXPOSE 8080
VOLUME /data

CMD ["/bin/bash", "/entrypoint.sh"]
