import json
import os

from confluent_kafka import Consumer
from dotenv import load_dotenv
from groq import Groq

from worker.prompt import parse_llm_json, system_prompt, user_prompt

from .producer import send_aidone_work

# Load .env ONLY if present (safe)
load_dotenv()

KAFKA_BOOTSTRAP = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9092")
GROQ_API_KEY = os.getenv("GROQ_API_KEY")

groq_client = Groq(api_key=GROQ_API_KEY)


def start_consumer():
    conf = {
        "bootstrap.servers": KAFKA_BOOTSTRAP,
        "group.id": "ai-workers",
        "auto.offset.reset": "earliest",
        "enable.auto.commit": False,
    }

    consumer = Consumer(conf)
    consumer.subscribe(["note.created"])
    print("🤖 Kafka AI consumer started")

    try:
        while True:
            msg = consumer.poll(1.0)
            if msg is None:
                continue
            if msg.error():
                print("❌ Kafka error:", msg.error())
                continue

            value = json.loads(msg.value().decode())
            print("📩 Received event:", value["note_id"])

            transcript = groq_client.audio.transcriptions.create(
                url=value["audio_url"],
                model="whisper-large-v3-turbo",
                response_format="verbose_json",
                language="en",
            )

            completion = groq_client.chat.completions.create(
                model="llama-3.1-8b-instant",
                messages=[
                    {"role": "system", "content": system_prompt},
                    {"role": "user", "content": user_prompt(transcript.text)},
                ],
                temperature=0.2,
                max_completion_tokens=1024,
            )

            llm_text = completion.choices[0].message.content
            if not llm_text:
                raise ValueError("LLM returned empty response")

            note_json = parse_llm_json(llm_text)

            send_aidone_work(
                value["note_id"],
                note_json["title"],
                note_json["content"],
            )

            consumer.commit(msg)
            print("✅ AI work committed")

    except KeyboardInterrupt:
        print("🛑 AI worker stopped")
    finally:
        consumer.close()
