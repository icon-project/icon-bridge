load("@rules_python//python:pip.bzl", "pip_install")

def dependencies():
    pip_install(
        name = "icon_util_python_deps",
        requirements = "//utils:requirements.txt",
    )