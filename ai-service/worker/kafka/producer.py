import json
import os
import time

from confluent_kafka import Producer

KAFKA_BOOTSTRAP = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")

producer = Producer(
    {
        "bootstrap.servers": KAFKA_BOOTSTRAP,
        "acks": "all",
    }
)


def send_aidone_work(note_id, title, transcript):
    event = {
        "event": "aiwork.done",
        "note_id": note_id,
        "title": title,
        "transcript": transcript,
        "created_at": time.strftime("%Y-%m-%dT%H:%M:%SZ"),
    }

    producer.produce(
        topic="aiwork.done",
        key=note_id.encode(),
        value=json.dumps(event).encode(),
    )
    producer.flush()
    print("📤 AI work sent")
