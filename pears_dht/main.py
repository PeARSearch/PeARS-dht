import grpc
from concurrent import futures
from threading import Thread
import opendht as dht
import argparse


import os.path, sys
sys.path.insert(0, os.path.dirname(os.path.dirname(__file__)))

from pears_dht.proto.messages.dht_pb2_grpc import DhtMessageServicer, add_DhtMessageServicer_to_server
from pears_dht.proto.messages.dht_pb2 import PutResponse, GetResponse, PutRequest, GetRequest


class DhtMessageService(DhtMessageServicer):
    isLeaf = True
    node = dht.DhtRunner()

    def __init__(self, port: int, bootstrap: str):
        print("Starting DHT node in port", port)
        self.node.run(port=port)
        if bootstrap != "":
            b_url = urlparse('//'+bootstrap)
            self.node.bootstrap(b_url.hostname, str(b_url.port) if b_url.port else '4222')

    def Put(self, request: PutRequest, context):
        # Implement your logic for Put here
        self.node.put(dht.InfoHash.get(request.key), dht.Value(request.value))
        return PutResponse(success=True)

    def Get(self, request: GetRequest, context):
        # Implement your logic for Get here
        results = self.node.get(dht.InfoHash.get(request.key))
        for r in results:
            print(r.data)
        return GetResponse(value=[r.data for r in results])
 

def serve_grpc(address: str, dht_port: str, bootstrap_node_ip: str):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    dht_message_servicer = DhtMessageServicer()
    add_DhtMessageServicer_to_server(DhtMessageService(int(dht_port), bootstrap_node_ip), server)
    server.add_insecure_port(address)
    server.start()
    print(f"grpc server listening on {address}")
    server.wait_for_termination()

def run():
    parser = argparse.ArgumentParser(description='DHT Node')
    subparsers = parser.add_subparsers(dest='command')

    parser.add_argument('--dhtport', default="4222", required=False, help='Port to start the DHT in (Other nodes in the network will use this to connect to this node)')
    parser.add_argument('--serverport', default="8080", required=False, help='Port to start the API in (clients will be able to feed data to the DHT on this port)')
    parser.add_argument('--bootstrapip', default="", required=False, help='IP of the node to connect to when bootstrapping, in the form 10.45.6.6:4222')

    args = parser.parse_args()

    return args


if __name__ == "__main__":
    args = run()

    grpc_address = "localhost:" + args.serverport

    # Start grpc server in a separate thread
    grpc_thread = Thread(target=serve_grpc, args=(grpc_address,args.dhtport,args.bootstrapip,))
    grpc_thread.start()

    # Wait for both threads to finish
    grpc_thread.join()