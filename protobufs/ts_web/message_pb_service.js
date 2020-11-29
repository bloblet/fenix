// package: 
// file: message.proto

var message_pb = require("./message_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var Messages = (function () {
  function Messages() {}
  Messages.serviceName = "Messages";
  return Messages;
}());

Messages.HandleMessages = {
  methodName: "HandleMessages",
  service: Messages,
  requestStream: true,
  responseStream: true,
  requestType: message_pb.CreateMessage,
  responseType: message_pb.Message
};

Messages.GetMessageHistory = {
  methodName: "GetMessageHistory",
  service: Messages,
  requestStream: false,
  responseStream: false,
  requestType: message_pb.RequestMessageHistory,
  responseType: message_pb.MessageHistory
};

exports.Messages = Messages;

function MessagesClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

MessagesClient.prototype.handleMessages = function handleMessages(metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.client(Messages.HandleMessages, {
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport
  });
  client.onEnd(function (status, statusMessage, trailers) {
    listeners.status.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners.end.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners = null;
  });
  client.onMessage(function (message) {
    listeners.data.forEach(function (handler) {
      handler(message);
    })
  });
  client.start(metadata);
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    write: function (requestMessage) {
      client.send(requestMessage);
      return this;
    },
    end: function () {
      client.finishSend();
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

MessagesClient.prototype.getMessageHistory = function getMessageHistory(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(Messages.GetMessageHistory, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

exports.MessagesClient = MessagesClient;

