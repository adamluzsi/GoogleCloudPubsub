eval `gcloud beta emulators pubsub env-init`

if ps aux | grep -v grep | grep -q gcloud
then
    gcloud beta emulators pubsub --quiet --host-port=$PUBSUB_EMULATOR_HOST start

else
    echo "process already running"
    exit 1

fi

