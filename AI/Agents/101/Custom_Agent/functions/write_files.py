import os

from google.genai import types


def write_file(working_directory, file_path, content):
    if not os.path.isdir(working_directory):
        return f'Error: "{working_directory}" is not a directory'

    working_dir_abs = os.path.abspath(working_directory)
    abs_file_path = os.path.join(working_dir_abs, file_path)

    valid_target_dir = (
        os.path.commonpath([working_dir_abs, abs_file_path]) == working_dir_abs
    )
    if not valid_target_dir:
        return f'Error: Cannot write "{file_path}" as it is outside the permitted working directory'

    if os.path.isdir(abs_file_path):
        return f'Error: Cannot write to "{file_path}" as it is a directory'

    os.makedirs(os.path.dirname(abs_file_path), exist_ok=True)

    try:
        with open(abs_file_path, "w") as f:
            f.write(content)
        return (
            f'Successfully wrote to "{file_path}" ({len(content)} characters written)'
        )
    except Exception as e:
        return f"Error: {e}"


schema_write_file = types.FunctionDeclaration(
    name="write_file",
    description="Write content to a file",
    parameters=types.Schema(
        type=types.Type.OBJECT,
        properties={
            "file_path": types.Schema(
                type=types.Type.STRING,
                description="The path to the file",
            ),
            "content": types.Schema(
                type=types.Type.STRING,
                description="The content to write to the file",
            ),
        },
    ),
)
