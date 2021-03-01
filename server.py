from aiohttp import web
import socketio
import os
import numpy as np
# creates a new Async Socket IO Server
sio = socketio.AsyncServer(cors_allowed_origins='*')

# Creates a new Aiohttp Web Application
app = web.Application()

# Binds our Socket.IO server to our Web App
# instance
sio.attach(app)
# we can define aiohttp endpoints just as we normally
# would with no change


# async def index(request):
#     with open('client/build/index.html') as f:
#         return web.Response(text=f.read(), content_type='text/html')

# If we wanted to create a new websocket endpoint,
# use this decorator, passing in the name of the
# event we wish to listen out for

CHAT = []
MATRIX = np.zeros((20,20))

@sio.on('connect')
def on_connect(a, b):
    print('connection established', a, " ", b['REMOTE_ADDR'])


@sio.on('message')
async def print_message(sid, message):
    # When we receive a new event of type
    # 'message' through a socket.io connection
    # we print the socket ID and the message
    print("Socket ID: ", sid)
    CHAT.append(message)
    print("CHAT:", CHAT)
    await sio.emit('chatUpdate', {'chat': CHAT})

@sio.on('boxplaced')
async def box_placed(sid, coords):
    print("Socket ID: ", sid)
    
    x = coords['x']
    y = coords['y']
    MATRIX[x,y] = 1

    print("Box placed on: (" , x, ',', y,')')

    await sio.emit('boxupdate', {'x': x, 'y': y})

# We bind our aiohttp endpoint to our app
# router
# app.router.add_get('/', index)

# We kick off our server
if __name__ == '__main__':
    web.run_app(app, host='0.0.0.0', port=7777)
