import os


def get_files_info(working_directory, directory="."):
    if not os.path.isdir(working_directory):
        print(f'Error: "{working_directory}" is not a directory')
    if not os.path.isdir(directory):
        print(f'Error: "{directory}" is not a directory')

    working_dir_abs = os.path.abspath(working_directory)
    target_dir = os.path.normpath(os.path.join(working_dir_abs, directory))
    valid_target_dir = (
        os.path.commonpath([working_dir_abs, target_dir]) == working_dir_abs
    )
    if not valid_target_dir:
        print(
            f'Error: Cannot list "{directory}" as it is outside the permitted working directory'
        )
        return
    try:
        print(
            f"Result for {'current' if working_directory == '.' else working_directory} directory:"
        )
        for f in os.listdir(target_dir):
            f_path = os.path.join(target_dir, f)
            file_sz = os.path.getsize(f_path)
            is_dir = os.path.isdir(f_path)
            print(f"- {f}: file_size={file_sz} bytes, is_dir={is_dir}")
    except Exception as e:
        print("Error: ", e)
