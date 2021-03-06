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
const proto = require('./auth_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.AuthClient =
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
proto.AuthPromiseClient =
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
 *   !proto.ClientAuth,
 *   !proto.AuthAck>}
 */
const methodDescriptor_Auth_Login = new grpc.web.MethodDescriptor(
  '/Auth/Login',
  grpc.web.MethodType.UNARY,
  proto.ClientAuth,
  proto.AuthAck,
  /**
   * @param {!proto.ClientAuth} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.AuthAck.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ClientAuth,
 *   !proto.AuthAck>}
 */
const methodInfo_Auth_Login = new grpc.web.AbstractClientBase.MethodInfo(
  proto.AuthAck,
  /**
   * @param {!proto.ClientAuth} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.AuthAck.deserializeBinary
);


/**
 * @param {!proto.ClientAuth} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.AuthAck)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.AuthAck>|undefined}
 *     The XHR Node Readable Stream
 */
proto.AuthClient.prototype.login =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/Auth/Login',
      request,
      metadata || {},
      methodDescriptor_Auth_Login,
      callback);
};


/**
 * @param {!proto.ClientAuth} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.AuthAck>}
 *     Promise that resolves to the response
 */
proto.AuthPromiseClient.prototype.login =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/Auth/Login',
      request,
      metadata || {},
      methodDescriptor_Auth_Login);
};


module.exports = proto;

