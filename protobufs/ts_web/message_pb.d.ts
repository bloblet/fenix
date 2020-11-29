// package: 
// file: message.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class CreateMessage extends jspb.Message {
  getContent(): string;
  setContent(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateMessage.AsObject;
  static toObject(includeInstance: boolean, msg: CreateMessage): CreateMessage.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateMessage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateMessage;
  static deserializeBinaryFromReader(message: CreateMessage, reader: jspb.BinaryReader): CreateMessage;
}

export namespace CreateMessage {
  export type AsObject = {
    content: string,
  }
}

export class Message extends jspb.Message {
  getMessageid(): string;
  setMessageid(value: string): void;

  getUserid(): string;
  setUserid(value: string): void;

  hasSentat(): boolean;
  clearSentat(): void;
  getSentat(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setSentat(value?: google_protobuf_timestamp_pb.Timestamp): void;

  getContent(): string;
  setContent(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Message.AsObject;
  static toObject(includeInstance: boolean, msg: Message): Message.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Message, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Message;
  static deserializeBinaryFromReader(message: Message, reader: jspb.BinaryReader): Message;
}

export namespace Message {
  export type AsObject = {
    messageid: string,
    userid: string,
    sentat?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    content: string,
  }
}

export class MessageHistory extends jspb.Message {
  clearMessagesList(): void;
  getMessagesList(): Array<Message>;
  setMessagesList(value: Array<Message>): void;
  addMessages(value?: Message, index?: number): Message;

  getNumberofmessages(): number;
  setNumberofmessages(value: number): void;

  getPages(): number;
  setPages(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MessageHistory.AsObject;
  static toObject(includeInstance: boolean, msg: MessageHistory): MessageHistory.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: MessageHistory, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MessageHistory;
  static deserializeBinaryFromReader(message: MessageHistory, reader: jspb.BinaryReader): MessageHistory;
}

export namespace MessageHistory {
  export type AsObject = {
    messagesList: Array<Message.AsObject>,
    numberofmessages: number,
    pages: number,
  }
}

export class RequestMessageHistory extends jspb.Message {
  hasLastmessagetime(): boolean;
  clearLastmessagetime(): void;
  getLastmessagetime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLastmessagetime(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RequestMessageHistory.AsObject;
  static toObject(includeInstance: boolean, msg: RequestMessageHistory): RequestMessageHistory.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RequestMessageHistory, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RequestMessageHistory;
  static deserializeBinaryFromReader(message: RequestMessageHistory, reader: jspb.BinaryReader): RequestMessageHistory;
}

export namespace RequestMessageHistory {
  export type AsObject = {
    lastmessagetime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

