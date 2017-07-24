source $HOME/.gcloud_profile


for i in 1 2 3 4 5 6
do
  echo $(echo "Y" | gcloud beta emulators pubsub --help) > /dev/null
  echo $(echo "Y" | gcloud beta emulators pubsub env-init) > /dev/null
done

echo "initialize is done"

pubsubENV=$(gcloud beta emulators pubsub env-init)
echo $pubsubENV
eval $pubsubENV

if ps aux | grep -v grep | grep -q gcloud
then
    echo "process already running"
    exit 1

else
    gcloud beta emulators pubsub start --quiet --host-port=$PUBSUB_EMULATOR_HOST &

    # Bad!
    sleep 2

fi

