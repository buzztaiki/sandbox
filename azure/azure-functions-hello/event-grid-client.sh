#!/bin/bash

body() {
  cat <<EOF
{
  "id": "Incoming_20200918002745d29ebbea-3341-4466-9690-0a03af35228e",
  "topic": "/subscriptions/aaaa0a0a-bb1b-cc2c-dd3d-eeeeee4e4e4e/resourcegroups/test/providers/microsoft.communication/communicationservices/test-resource",
  "subject": "/phonenumber/15555555555",
  "data": {
    "MessageId": "Incoming_20200918002745d29ebbea-3341-4466-9690-0a03af35228e",
    "From": "15555555555",
    "To": "15555555555",
    "Message": "Great to connect with Azure Communication Services events",
    "ReceivedTimestamp": "2020-09-18T00:27:45.32Z"
  },
  "eventType": "Microsoft.Communication.SMSReceived",
  "dataVersion": "1.0",
  "metadataVersion": "1",
  "eventTime": "2020-09-18T00:27:47Z"
}
EOF
}

curl -v -H 'aeg-event-type: Notification' 'http://localhost:7071/runtime/webhooks/EventGrid?functionName=EventGridExample' --json @- < <(body)
