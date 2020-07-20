import asyncio
import time

async def tcp_echo_client(message, loop):
    while True:
        reader, writer = await asyncio.open_connection('127.0.0.1', 5000,
                                                   loop=loop)
        print('Send: %r' % message)
        writer.write(message.encode())

        data = await reader.read(100)
        print('Received: %r' % data.decode())
        time.sleep(3)

        print('WAIT Close the socket')
        writer.close()

message = 'Hello World!'
loop = asyncio.get_event_loop()
loop.run_until_complete(tcp_echo_client(message, loop))
loop.close()