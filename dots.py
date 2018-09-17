#!/usr/bin/env python3
""" Util for backup / deploy dotfiles.
"""

import os
import re
import sys
import shutil
import platform
import datetime
from subprocess import run, PIPE
import yaml

HOME_DIR = os.getenv("HOME")
HOME_PIP = "[_HOME_]"
THIS_DIR = os.path.dirname(os.path.realpath(__file__))
SYS_NAME = platform.system().lower()
COLORS = {
    "x": "\x1b[0m",
    "w": "\x1b[37m%s\x1b[0m",
    "_": "\x1b[90m%s\x1b[0m",
    "r": "\x1b[31m%s\x1b[0m",
    "y": "\x1b[33m%s\x1b[0m",
    "g": "\x1b[32m%s\x1b[0m",
    "b": "\x1b[34m%s\x1b[0m",
    "m": "\x1b[35m%s\x1b[0m",
    "c": "\x1b[36m%s\x1b[0m",
}

with open(os.path.join(THIS_DIR, "dots.yml")) as fs:
    try:
        CONF = list(yaml.safe_load_all(fs))[0]
    except yaml.YAMLError as err:
        print(COLORS["r"] % " ✘ Fucked up:")
        print(err)
        sys.exit()
PATHS = CONF['sources']
ARCHIVE = CONF['archive_dir']
DOTS_DIR = os.path.join(THIS_DIR, "dotfiles")


def backup():
    """ Parse listed paths and copy it's content to this repo
        with storing previous version.
    """

    if not os.path.exists(DOTS_DIR):
        os.makedirs(DOTS_DIR)

    # Copying files
    clone_targets()

    # Git commit
    commit_changes()

    # Archive
    archive()


def clone_targets():
    """ Clone files listed in dots.yml in ./dotfiles
    """
    print(COLORS["b"] % "\n → Copying files...")
    for name, path in PATHS.items():
        print("   '%s' to '%s'" % (path, os.path.join("dotfiles", name)))
        target_path = path.replace("~", HOME_DIR)
        backup_path = os.path.join(DOTS_DIR, name)

        shutil.rmtree(backup_path, True)

        if os.path.isdir(target_path):
            clone_dir(target_path, backup_path)
        elif os.path.isfile(target_path):
            try:
                safe_clone_file(target_path, backup_path)
            except OSError:
                continue

    print(COLORS["g"] % " ✓ Done")


def clone_dir(src_path, dest_dir):
    """ Clone dir recursively.
    """
    if not os.path.exists(dest_dir):
        os.makedirs(dest_dir)
    for (dirpath, dirnames, filenames) in os.walk(src_path):
        rel_path = os.path.relpath(dirpath, src_path)
        # Create dirs
        for dirname in dirnames:
            dirname = os.path.join(dest_dir, rel_path, dirname)
            dirname = os.path.abspath(dirname)
            if not os.path.exists(dirname):
                os.makedirs(dirname)
        # Clone files
        for filename in filenames:
            src = os.path.join(dirpath, filename)
            dest = os.path.join(dest_dir, rel_path, filename)
            dest = os.path.abspath(dest)
            try:
                safe_clone_file(src, dest)
            except OSError:
                continue


def safe_clone_file(src_path, dest_path):
    """ If file is text, read it, replace all home
        path occurances and write to dest path.
        If not text - just copy.
    """
    with open(src_path) as file:
        try:
            normalized = file.read().replace(HOME_DIR, HOME_PIP)
            out = open(dest_path, "w")
            out.write(normalized)
            out.close()
        except ValueError:
            shutil.copy2(src_path, dest_path)
        except Exception as err:
            print(COLORS["r"] % " ✘ Fucked up with ", end="")
            print(src_path)
            raise err


def commit_changes():
    """ Git commit.
    """
    print(COLORS["b"] % "\n → Git commit...")
    git_status = run(["git", "status", "--porcelain"],
                     stdout=PIPE).stdout.decode("utf-8")
    for dif in git_status.split("\n"):
        dif_match = re.match(r"^.(.+)\s(.*)", dif)
        if dif_match is None:
            continue
        if dif_match.group(1) == "M":
            print(COLORS["b"] % ("   ~ " + dif_match.group(2)))
        elif dif_match.group(1) == "D":
            print(COLORS["r"] % ("   x " + dif_match.group(2)))
        else:
            print(COLORS["_"] % ("   ? " + dif_match.group(2)))
    run(["git", "add", "."], stdout=PIPE)
    run(["git", "commit", "-m", "'Backup'"], stdout=PIPE)
    print(COLORS["g"] % " ✓ Done")


def archive():
    """ Archive dotfiles
    """
    if not ARCHIVE:
        return

    print(COLORS["b"] % "\n → Create archive...")
    time = datetime.datetime.now().strftime("%Y.%m.%d-%H.%M")
    target_dir = ARCHIVE.replace("~", HOME_DIR)
    if not os.path.exists(target_dir):
        os.makedirs(target_dir)
    filename = "dotfiles-" + time
    dest_path = os.path.join(target_dir, filename)
    dest_path = os.path.abspath(dest_path)
    print(COLORS["_"] % "   " + dest_path + ".zip")
    shutil.make_archive(dest_path, "zip", DOTS_DIR)
    print(COLORS["g"] % " ✓ Done")


def restore():
    """ Copy content of this repo to destinition paths.
    """
    pass


def help_msg():
    """ Prints help message.
    """
    print("""
 Commands:
   dots.py backup - Parse listed paths and copy it's content to this repo
                    with storing previous version.
   dots.py restore username - Copy content of this repo to destination paths.
""")


if "backup" in sys.argv:
    backup()
elif "restore" in sys.argv:
    restore()
else:
    help_msg()
