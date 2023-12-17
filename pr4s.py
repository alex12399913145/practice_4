from flask import Flask, request

from datetime import datetime
import random

import socket
import json

app = Flask(__name__)

state = []

@app.route('/report', methods=['POST'])
def report():
    return json.dumps(state)

def connection(input_value):
    client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client.connect(('127.0.0.1', 6379))
    request_data = f'post4{input_value}\n'
    client.send(request_data.encode('utf-8'))
    response_data = client.recv(1024).decode('utf-8')
    client.send(request_data.encode('utf-8'))
    response_data = client.recv(1024).decode('utf-8')
    client.close()
    print(response_data)

@app.route('/', methods=['POST'])
def get_statistics():
    dimensions = request.get_json()

    time = datetime.now().time()

    url_found = False
    for url in state:
        if url['URL'] == dimensions["URL"]:
            time_state = url['TimeInterval'][:-5] + str(time)[:5]
            url['TimeInterval'] = time_state
            url['Count'] += 1
            url_found = True
            break

    if not url_found:
        Id_state = random.randint(1, 100)
        dimensions['Count'] = 1
        time_state = dimensions['TimeInterval'] + "-" + str(time)[:5]
        dimensions['TimeInterval'] = time_state
        dimensions['ID'] = Id_state
        state.append(dimensions)


    input_value = ";".join([dimensions['TimeInterval'], dimensions['SourceIP'], dimensions['URL']])

    connection(input_value)

    print(state)

    return json.dumps({"message": "Statistics updated." if time_state else "Statistics added."})

if __name__ == '__main__':
    app.run(host="127.0.0.1", port=5000)