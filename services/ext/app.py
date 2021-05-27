import time
import os

import redis
from flask import Flask

app = Flask(__name__)
cache = redis.Redis(host='0.0.0.0', port=6379)
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
    #count = get_hit_count()
    #return 'Hello World! you have seen {} times\n'.format(count)
    return 'Hello world!'

if __name__ == "__main__":
    print("Start listening in port: {port}".format(port=os.environ["API_PORT"]))
    app.run(host='0.0.0.0', port=os.environ["API_PORT"])
