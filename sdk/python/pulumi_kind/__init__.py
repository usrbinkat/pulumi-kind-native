# coding=utf-8
# *** WARNING: this file was generated by pulumi-language-python. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

from . import _utilities
import typing
# Export this package's modules as members:
from .kind import *
from .provider import *
_utilities.register(
    resource_modules="""
[
 {
  "pkg": "kind",
  "mod": "index",
  "fqn": "pulumi_kind",
  "classes": {
   "kind:index:Kind": "Kind"
  }
 }
]
""",
    resource_packages="""
[
 {
  "pkg": "kind",
  "token": "pulumi:providers:kind",
  "fqn": "pulumi_kind",
  "class": "Provider"
 }
]
"""
)
