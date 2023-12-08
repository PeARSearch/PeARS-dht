import grpc
from concurrent import futures
from threading import Thread

import sys, os
sys.path.append(os.path.dirname(os.path.realpath(__file__)) + '/proto/messages')

from pears_dht.proto.messages.dht_pb2_grpc import DhtMessageServicer, add_DhtMessageServicer_to_server
 

def serve_grpc(address: str):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    dht_message_servicer = DhtMessageServicer()
    add_DhtMessageServicer_to_server(dht_message_servicer, server)
    server.add_insecure_port(address)
    server.start()
    print(f"grpc server listening on {address}")
    server.wait_for_termination()

if __name__ == "__main__":
    grpc_address = "localhost:8080"
    dht_port = 50051

    # Start grpc server in a separate thread
    grpc_thread = Thread(target=serve_grpc, args=(grpc_address,))
    grpc_thread.start()

    # Wait for both threads to finish
    grpc_thread.join()