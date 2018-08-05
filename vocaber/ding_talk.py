#-*- coding: utf-8 -*-
import requests
import json


webhook = "https://oapi.dingtalk.com/robot/send?access_token=1a663c12d7a937ca828ae2108c2c22df867d7ff9f03e60be1f31be30a19f7032"
headers = {'Content-Type': 'application/json;charset=UTF-8'}


def get_url_res(url, headers={}, tries=2, timeout=19):
    '''get response of request'''
    times = 0
    res = None
    while(times < tries):
        times +=1
        try:
            res = requests.get(url=url, headers=headers, timeout=timeout)
            break
        except:
            continue

    return res 

url = "http://127.0.0.1:5000/yesterdays_item"
res = get_url_res(url)
json_obj = json.loads(res.text)
items = json_obj["items"]
words_list = []
for item in items:
    words_list.append(item["value"])
    words_list.append("\n")

words = "".join(words_list)
msg = {
    "msgtype" :"markdown",
    "markdown": {
        "title": "昨日单词回顾",
        "text": "### 昨日单词回顾：\n %s" % words
    }
}
print(msg)
res = requests.post(webhook, headers=headers, data=json.dumps(msg)) 

