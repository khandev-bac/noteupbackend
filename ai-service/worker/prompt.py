import json
import re

system_prompt = """You are a professional note-writing assistant.

Your task is to transform raw speech transcripts into clear, accurate, and well-structured written notes.

STRICT RULES:
- Output MUST be valid JSON only.
- Do NOT include any text before or after the JSON.
- Follow the JSON schema provided by the user exactly.
- Do NOT add, infer, or assume any information that is not explicitly present in the transcript.
- Preserve the original meaning and intent of the speaker at all times.

WRITING GUIDELINES:
- Remove filler words, repetitions, and speech disfluencies.
- Convert fragmented or incomplete spoken sentences into complete, clear written sentences.
- Improve clarity, organization, and logical flow.
- Use polished, professional, and neutral language.
- Organize content using clear sections, paragraphs, and lists when appropriate.

ERROR HANDLING:
- If the transcript is unclear, incomplete, or ambiguous, rewrite it conservatively without guessing.
- Never fabricate facts, examples, or explanations.

Your response must contain ONLY the final JSON output.

"""


def user_prompt(text: str) -> str:
    return f"""
Convert the following transcript into clear, structured, and well-written study notes.

Return the output as STRICT JSON using exactly the following schema:
{{
  "title": "string",
  "content": "string"
}}

CONTENT REQUIREMENTS (VERY IMPORTANT):

1. Writing Style & Quality
- Write in clear, formal, and academic language suitable for studying.
- Use complete, grammatically correct sentences.
- Remove filler words, repetitions, and spoken-language artifacts.
- Expand unclear or incomplete spoken sentences into full explanations.
- Maintain a smooth, natural flow between ideas.

2. Structure
- Organize the content using clear section headings.
- Each section should contain well-developed paragraphs.
- Do NOT use bullet points, numbering, symbols, or ticks.
- Do NOT summarize or shorten the content.

3. Depth & Clarity
- Explain concepts thoroughly with sufficient detail.
- Add missing context where ideas are implied but not fully explained.
- Ensure the notes are comprehensive and self-contained.

4. Tables
- Do NOT use tables.

FORMATTING RULES:
- Do NOT use markdown.
- Do NOT include bullet points, lists, symbols, or summary sections.
- Do NOT include any text outside the JSON.
- The JSON must be valid and directly parsable.
- The "content" field must contain plain text only.

Transcript:
\"\"\" {text} \"\"\"
"""


def parse_llm_json(text: str | None) -> dict:
    if text is None:
        raise ValueError("LLM returned empty content")

    text = text.strip()

    # Remove markdown fences if present
    if text.startswith("```"):
        text = text.strip("`")
        if text.lower().startswith("json"):
            text = text[4:].strip()

    # 🔥 CRITICAL: remove invalid control characters
    text = re.sub(r"[\x00-\x1f\x7f]", "", text)

    return json.loads(text)
