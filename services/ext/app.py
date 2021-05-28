#!/usr/bin/env python3

import time
import os

import redis
from flask import Flask

app = Flask(__name__)

# redis_url = os.getenv('REDISTOGO_URL')
# print('redis_url:', redis_url)

cache = redis.Redis(host='0.0.0.0', port=os.environ["REDIS_PORT"])
app.config["DEBUG"] = True

def get_hit_count():
    retries = 5
    while True:
        try:
            return cache.incr('hits')
        except redis.exceptions.ConnectionError as exc:
            if retries == 0:
                raise exc
            retries -= 1
            time.sleep(0.5)

@app.route('/')
def hello():
    count = get_hit_count()
    return "Hello"

if __name__ == "__main__":
    print("Start listening in port: {port}".format(port=os.environ["API_PORT"]))
    app.run(host='0.0.0.0', port=os.environ["API_PORT"])
