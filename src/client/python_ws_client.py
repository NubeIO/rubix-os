import websocket  #pip install websocket-client
import os.path
import configparser
import argparse

home = os.path.expanduser('~')
parser = argparse.ArgumentParser("path")
parser.add_argument("path", help="path to conf file", type=str)
args = parser.parse_args()
config_path =  args.path #  ~/code/go/flow-framework/client/ws.conf


if not os.path.isfile(config_path):
    from shutil import copyfile
    copyfile('ws.conf', config_path)

config = configparser.ConfigParser()
config.read(config_path)
domain = config.get('server', 'domain', fallback=None)
token = config.get('server', 'token')
ssl = "true" == config.get('server', 'ssl', fallback='false').lower()


def on_message(ws, message):
    print(message)


def on_error(ws, error):
    print(error)


def on_close(ws):
    print("### closed ###")


if __name__ == "__main__":
    websocket.enableTrace(True)
    if ssl:
        protocol = 'ws'
    else:
        protocol = 'ws'

    ws = websocket.WebSocketApp("{}://{}/stream?token={}".format(protocol, domain, token),
                                on_message=on_message,
                                on_error=on_error,
                                on_close=on_close)
    ws.run_forever()
