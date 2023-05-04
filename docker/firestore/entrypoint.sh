#!/bin/bash

: ${CLOUDSDK_CORE_PROJECT:=project}

FIRESTORE_PID=

on_exit() {
  firebase emulators:export --project "${CLOUDSDK_CORE_PROJECT}" --only firestore -f /data/firestore
  kill -TERM "${FIRESTORE_PID}"
}

trap "on_exit" SIGINT SIGTERM SIGHUP

OPTS=
if [[ -d "/data/firestore" ]]; then
  OPTS="--import /data/firestore"
fi

firebase emulators:start --project "${CLOUDSDK_CORE_PROJECT}" --only firestore ${OPTS} &
FIRESTORE_PID="$!"
wait "${FIRESTORE_PID}"
