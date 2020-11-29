/**
 * @fileoverview gRPC-Web generated client stub for 
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')
const proto = require('./message_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.MessagesClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.MessagesPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.RequestMessageHistory,
 *   !proto.MessageHistory>}
 */
const methodDescriptor_Messages_GetMessageHistory = new grpc.web.MethodDescriptor(
  '/Messages/GetMessageHistory',
  grpc.web.MethodType.UNARY,
  proto.RequestMessageHistory,
  proto.MessageHistory,
  /**
   * @param {!proto.RequestMessageHistory} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.MessageHistory.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.RequestMessageHistory,
 *   !proto.MessageHistory>}
 */
const methodInfo_Messages_GetMessageHistory = new grpc.web.AbstractClientBase.MethodInfo(
  proto.MessageHistory,
  /**
   * @param {!proto.RequestMessageHistory} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.MessageHistory.deserializeBinary
);


/**
 * @param {!proto.RequestMessageHistory} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.MessageHistory)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.MessageHistory>|undefined}
 *     The XHR Node Readable Stream
 */
proto.MessagesClient.prototype.getMessageHistory =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/Messages/GetMessageHistory',
      request,
      metadata || {},
      methodDescriptor_Messages_GetMessageHistory,
      callback);
};


/**
 * @param {!proto.RequestMessageHistory} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.MessageHistory>}
 *     Promise that resolves to the response
 */
proto.MessagesPromiseClient.prototype.getMessageHistory =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/Messages/GetMessageHistory',
      request,
      metadata || {},
      methodDescriptor_Messages_GetMessageHistory);
};


module.exports = proto;

