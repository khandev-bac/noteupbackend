# ============================
# Python AI Service (glibc)
# ============================
FROM python:3.11-bookworm

# Prevent Python from buffering logs
ENV PYTHONUNBUFFERED=1
ENV PYTHONDONTWRITEBYTECODE=1

# Install system dependencies required by confluent-kafka
RUN apt-get update && apt-get install -y \
    gcc \
    librdkafka-dev \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy requirements first (better Docker caching)
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy application source code
COPY . .

# Start the Kafka AI worker
CMD ["python", "main.py"]
