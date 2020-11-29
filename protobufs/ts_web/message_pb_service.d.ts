// package: 
// file: message.proto

import * as message_pb from "./message_pb";
import {grpc} from "@improbable-eng/grpc-web";

type MessagesHandleMessages = {
  readonly methodName: string;
  readonly service: typeof Messages;
  readonly requestStream: true;
  readonly responseStream: true;
  readonly requestType: typeof message_pb.CreateMessage;
  readonly responseType: typeof message_pb.Message;
};

type MessagesGetMessageHistory = {
  readonly methodName: string;
  readonly service: typeof Messages;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof message_pb.RequestMessageHistory;
  readonly responseType: typeof message_pb.MessageHistory;
};

export class Messages {
  static readonly serviceName: string;
  static readonly HandleMessages: MessagesHandleMessages;
  static readonly GetMessageHistory: MessagesGetMessageHistory;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class MessagesClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  handleMessages(metadata?: grpc.Metadata): BidirectionalStream<message_pb.CreateMessage, message_pb.Message>;
  getMessageHistory(
    requestMessage: message_pb.RequestMessageHistory,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: message_pb.MessageHistory|null) => void
  ): UnaryResponse;
  getMessageHistory(
    requestMessage: message_pb.RequestMessageHistory,
    callback: (error: ServiceError|null, responseMessage: message_pb.MessageHistory|null) => void
  ): UnaryResponse;
}

