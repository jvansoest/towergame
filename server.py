from aiohttp import web
import socketio
import os
import numpy as np
from datetime import datetime


class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKCYAN = '\033[96m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'


sio = socketio.AsyncServer(cors_allowed_origins='*')
app = web.Application()
# Binds our Socket.IO server to our Web App instance
sio.attach(app)


STATE = {"messages": [],
         "matrix": np.zeros((20, 20))}


@sio.on('connect')
async def on_connect(a, b):
    print(bcolors.WARNING + 'connection established' + bcolors.ENDC,
          a, " to ", b['REMOTE_ADDR'])

    #STATE['messages'] = []
    #STATE['matrix'] = np.zeros((20, 20))

    await sio.emit('chatUpdate', {'chat': STATE['messages']})


@sio.on('disconnect')
def disconnect(sid):
    print(bcolors.WARNING + 'connection disconnected' + bcolors.ENDC, sid)


def get_time_string():
    now = datetime.now()
    current_time = now.strftime("%H:%M:%S")
    return f'[{current_time}] '


@sio.on('message')
async def print_message(sid, message):

    current_time = get_time_string()
    STATE['messages'].append(current_time+message)

    print(bcolors.WARNING +
          current_time +
          bcolors.ENDC +
          bcolors.BOLD +
          "Socket ID: " +
          bcolors.ENDC,
          sid,
          bcolors.BOLD +
          " received message: " +
          bcolors.ENDC +
          "[" +
          message[:20] + "...]")

    await sio.emit('chatUpdate', {'chat': STATE['messages']})


@sio.on('boxplaced')
async def box_placed(sid, coords):
    print("Socket ID: ", sid)
    x = coords['x']
    y = coords['y']
    print("Box placed on: (", x, ',', y, ')')

    STATE['matrix'][x, y] = 1

    await sio.emit('boxupdate', {'x': x, 'y': y, 'matrix': STATE['matrix'].tolist()})


if __name__ == '__main__':
    web.run_app(app, host='0.0.0.0', port=7777)
