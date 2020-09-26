import os
import textwrap
import requests
import subprocess
import re
from datetime import datetime


def define_env(env):
    """
    This is the hook for defining variables, macros and filters

    - variables: the dictionary that contains the environment variables
    - macro: a decorator function, to declare a macro.
    """

    @env.macro
    def embed_code(filename, fragment=None, prefix=''):
        """
        Embed code, optionally specified by fragment name.
        """
        full_filename = os.path.join(env.project_dir, filename)
        extension = filename.rsplit('.', 1)[-1]
        with open(full_filename, 'r') as f:
            lines = f.readlines()

        str = None
        if fragment != None:
            i = 0
            found = None
            for line in lines:
                if line.strip() == '/// ['+fragment+']':
                    if found == None:
                        found = i+1
                    else:
                        str = ''.join(lines[found:i])
                        str = textwrap.dedent(str)
                        break
                i += 1
        else:
            str = ''.join(lines)

        return textwrap.indent('```'+extension+'\n'+str+'\n```', prefix)[len(prefix):]

    @env.macro
    def changelog():
        """
        Generate Changelog.
        """
        result = ""
        url = "https://api.github.com/repos/Fs02/rel/releases"
        releases = requests.get(url).json()

        for release in releases:
            body = ""
            if release["author"]["login"] == "github-actions[bot]":
                pattern = re.compile(r"^\w{7}\s.+\(\#\d+\)$")
                lines = release["body"].splitlines()
                for line in lines:
                    if pattern.match(line):
                        body += "\n- " + line[7:]
            else:
                body = release["body"]

            time = datetime.strptime(
                release["created_at"], '%Y-%m-%dT%H:%M:%SZ')
            result += "\n## **" + \
                release["name"] + "** - " + \
                time.strftime("%B %-d, %Y") + "\n\n" + body

        return result

    @env.macro
    def godoc(pkg):
        """
        Generate Godoc
        """
        result = subprocess.run(
            ['godoc2md', '-template=docs/reference/godoc.tpl', pkg], stdout=subprocess.PIPE)
        return result.stdout.decode("utf-8")
