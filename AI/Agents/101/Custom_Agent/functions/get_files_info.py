import os

from google.genai import types


def get_files_info(working_directory, directory="."):
    if not os.path.isdir(working_directory):
        return f'Error: "{working_directory}" is not a directory'
    if not os.path.isdir(directory):
        return f'Error: "{directory}" is not a directory'

    working_dir_abs = os.path.abspath(working_directory)
    target_dir = os.path.normpath(os.path.join(working_dir_abs, directory))

    valid_target_dir = (
        os.path.commonpath([working_dir_abs, target_dir]) == working_dir_abs
    )
    if not valid_target_dir:
        return f'Error: Cannot list "{directory}" as it is outside the permitted working directory'

    try:
        print(
            f"Result for {'current' if working_directory == '.' else working_directory} directory:"
        )
        for f in os.listdir(target_dir):
            f_path = os.path.join(target_dir, f)
            file_sz = os.path.getsize(f_path)
            is_dir = os.path.isdir(f_path)
            return f"- {f}: file_size={file_sz} bytes, is_dir={is_dir}"
    except Exception as e:
        return "Error: ", e


schema_get_files_info = types.FunctionDeclaration(
    name="get_files_info",
    description="Lists files in a specified directory relative to the working directory, providing file size and directory status",
    parameters=types.Schema(
        type=types.Type.OBJECT,
        properties={
            "directory": types.Schema(
                type=types.Type.STRING,
                description="Directory path to list files from, relative to the working directory (default is the working directory itself)",
            ),
        },
    ),
)
