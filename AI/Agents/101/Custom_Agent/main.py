import argparse
import os

from dotenv import load_dotenv
from google import genai
from google.genai import types

from call_function import available_functions
from prompts import system_prompt


def main():
    load_dotenv()
    api_key = os.environ.get("GEMINI_API_KEY")
    if api_key is None:
        raise RuntimeError("Gemini api key is missing!")

    parser = argparse.ArgumentParser(description="VOID CHATBOT")
    parser.add_argument("user_prompt", type=str, help="User prompt")
    parser.add_argument("--verbose", action="store_true", help="Enable verbose output")
    args = parser.parse_args()

    # "parts" argument is a list because Google Gemini API is multimodal,
    # which means a single message can contain multiple types of content at once.
    messages = [types.Content(role="user", parts=[types.Part(text=args.user_prompt)])]

    client = genai.Client(api_key=api_key)
    response = client.models.generate_content(
        model="gemini-2.5-flash",
        contents=messages,
        config=types.GenerateContentConfig(
            tools=[available_functions], system_instruction=system_prompt
        ),
    )

    if response.usage_metadata is None:
        raise RuntimeError("Failed API request")

    prompt_tokens = response.usage_metadata.prompt_token_count
    response_tokens = response.usage_metadata.candidates_token_count

    if args.verbose:
        print("User prompt: ", args.user_prompt)
        print("Prompt tokens: ", prompt_tokens)
        print("Response tokens: ", response_tokens)

    print("Response: ")
    print(response.text)

    if response.function_calls is not None:
        for function_call in response.function_calls:
            print(f"Calling function: {function_call.name}({function_call.args})")


if __name__ == "__main__":
    main()
