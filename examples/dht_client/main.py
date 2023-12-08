import argparse
import grpc
from pears_dht.proto.messages import dht_pb2
from pears_dht.proto.messages import dht_pb2_grpc

def put(stub, key, value):
    put_request = dht_pb2.PutRequest(key=key, value=value.encode())
    put_response = stub.Put(put_request)
    print("PutResponse:", put_response.success)

def get(stub, key):
    get_request = dht_pb2.GetRequest(key=key)
    get_response = stub.Get(get_request)
    print("GetResponse:", get_response.value)

def run():
    parser = argparse.ArgumentParser(description='DHT Client')
    subparsers = parser.add_subparsers(dest='command')

    # Put command
    put_parser = subparsers.add_parser('put')
    put_parser.add_argument('--key', required=True, help='Key to put in DHT')
    put_parser.add_argument('--value', required=True, help='Value to put in DHT')

    # Get command
    get_parser = subparsers.add_parser('get')
    get_parser.add_argument('--key', required=True, help='Key to get from DHT')

    args = parser.parse_args()

    # Assuming your gRPC server is running on localhost at port 50051
    channel = grpc.insecure_channel('localhost:50051')
    stub = dht_pb2_grpc.DhtMessageStub(channel)

    if args.command == 'put':
        put(stub, args.key, args.value)
    elif args.command == 'get':
        get(stub, args.key)

if __name__ == '__main__':
    run()