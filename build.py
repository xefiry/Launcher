"""
Script to build the Launcher and zip it with ressources/dependencies.

For the moment, only works for Windows without CGO.
"""

from zipfile import ZipFile
import subprocess
import os

NAME = "Launcher"
ARCH = "win-x64"
EXE = "launcher.exe"

# list of files and dirs to include in output archive
FILES = [EXE, "config.toml", "raylib.dll"]
DIRS = ["Fonts"]


def run(args: list[str]) -> str:
    """Run command and return the output.

    Args:
        args (list[str]): Command to execute with args.

    Raises:
        RuntimeError: If the process exit code is not 0.

    Returns:
        str: Output of the command.
    """

    result = subprocess.run(args, stdout=subprocess.PIPE)

    if result.returncode != 0:
        raise RuntimeError

    return result.stdout.decode("utf-8").strip()


def get_file_list() -> list[str]:
    """Generate the list of files to archive. Uses FILES & DIRS constants.

    Raises:
        RuntimeError: if on of the files in FILES does not exists

    Returns:
        list[str]: the list of files
    """

    output = []

    # add all the files
    for file in FILES:
        if os.path.exists(file):
            output.append(file)
        else:
            print(f"could not find {file}")
            raise RuntimeError

    # and all the directories with their content
    for dir in DIRS:
        for root, _, files in os.walk(dir):
            for file in files:
                output.append(f"{root}/{file}")

    return output


# build
run(["go", "build", "-ldflags", "-H=windowsgui -w -s"])

# get version number from tag
tag = run(["git", "describe", "--tags"])

# create the archive
with ZipFile(f"{NAME}_{tag}_{ARCH}.zip", "w") as myzip:
    for file in get_file_list():
        print(f"{file} -> {NAME}/{file}")
        myzip.write(file, f"{NAME}/{file}")

# clean executable
#run(["go", "clean"])
