# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: pears_dht/proto/messages/dht.proto
# Protobuf Python Version: 4.25.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\"pears_dht/proto/messages/dht.proto\x12\x0eproto.messages\"(\n\nPutRequest\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\x0c\"\x1e\n\x0bPutResponse\x12\x0f\n\x07success\x18\x01 \x01(\x08\"\x19\n\nGetRequest\x12\x0b\n\x03key\x18\x01 \x01(\t\"\x1c\n\x0bGetResponse\x12\r\n\x05value\x18\x01 \x01(\x0c\x32\x8c\x01\n\nDhtMessage\x12>\n\x03Put\x12\x1a.proto.messages.PutRequest\x1a\x1b.proto.messages.PutResponse\x12>\n\x03Get\x12\x1a.proto.messages.GetRequest\x1a\x1b.proto.messages.GetResponseb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'pears_dht.proto.messages.dht_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  DESCRIPTOR._options = None
  _globals['_PUTREQUEST']._serialized_start=54
  _globals['_PUTREQUEST']._serialized_end=94
  _globals['_PUTRESPONSE']._serialized_start=96
  _globals['_PUTRESPONSE']._serialized_end=126
  _globals['_GETREQUEST']._serialized_start=128
  _globals['_GETREQUEST']._serialized_end=153
  _globals['_GETRESPONSE']._serialized_start=155
  _globals['_GETRESPONSE']._serialized_end=183
  _globals['_DHTMESSAGE']._serialized_start=186
  _globals['_DHTMESSAGE']._serialized_end=326
# @@protoc_insertion_point(module_scope)
