#! /home/nautilis/local/env_py3/bin/python
from googletrans.gtoken import TokenAcquirer
import sys
acquirer = TokenAcquirer()
phrase = sys.argv[1]
tk = acquirer.do(phrase)
print(tk)
#print("2018.2018")
