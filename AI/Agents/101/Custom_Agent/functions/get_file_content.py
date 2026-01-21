import os

from google.genai import types

from config import MAX_CHARS


def get_file_content(working_directory, file_path):
    if not os.path.isdir(working_directory):
        print(f'Error: "{working_directory}" is not a directory')
        return

    working_dir_abs = os.path.abspath(working_directory)
    abs_file_path = os.path.join(working_dir_abs, file_path)

    valid_target_dir = (
        os.path.commonpath([working_dir_abs, abs_file_path]) == working_dir_abs
    )
    if not valid_target_dir:
        print(
            f'Error: Cannot read "{file_path}" as it is outside the permitted working directory'
        )
        return

    if not os.path.isfile(abs_file_path):
        print(f'Error: File not found or is not a regular file: "{file_path}"')
        return

    try:
        content = ""
        with open(abs_file_path, "r") as f:
            content += f.read(MAX_CHARS)
            if f.read(1):
                content += (
                    f'[...File "{file_path}" truncated at {MAX_CHARS} characters]'
                )
        print(content)
    except Exception as e:
        print("Error: ", e)


schema_get_file_content = types.FunctionDeclaration(
    name="get_file_content",
    description="Reads the content of a file within the permitted working directory",
    parameters=types.Schema(
        type=types.Type.OBJECT,
        properties={
            "file_path": types.Schema(
                type=types.Type.STRING,
                description="Path to the file within the permitted working directory",
            ),
        },
    ),
)
