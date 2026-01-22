import argparse
import os

from dotenv import load_dotenv
from google import genai
from google.genai import types

from call_function import available_functions, call_function
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
    messages = [
        types.Content(
            role="user",
            parts=[
                types.Part(
                    text=f"{args.user_prompt}\n\nContext: You're working in a calculator project with main.py and a pkg/ directory."
                )
            ],
        )
    ]

    client = genai.Client(api_key=api_key)

    for _ in range(15):
        response = client.models.generate_content(
            model="gemini-2.5-flash",
            contents=messages,
            config=types.GenerateContentConfig(
                tools=[available_functions],
                system_instruction=system_prompt,
            ),
        )

        if response.text:
            print("Response:")
            print(response.text)
            messages.append(
                types.Content(role="assistant", parts=[types.Part(text=response.text)])
            )
            break

        if response.function_calls:
            function_results = []

            for function_call in response.function_calls:
                function_call_result = call_function(function_call, verbose=True)
                first_part = (
                    function_call_result.parts[0]
                    if function_call_result.parts
                    else None
                )

                function_response = first_part.function_response if first_part else None
                if function_response is None:
                    raise RuntimeError("Tool response missing function_response")

                print(f"-> {function_response.response}")
                function_results.append(first_part)

            messages.append(types.Content(role="tool", parts=function_results))


if __name__ == "__main__":
    main()
