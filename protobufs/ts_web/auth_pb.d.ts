// package: 
// file: auth.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class ClientAuth extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClientAuth.AsObject;
  static toObject(includeInstance: boolean, msg: ClientAuth): ClientAuth.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ClientAuth, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClientAuth;
  static deserializeBinaryFromReader(message: ClientAuth, reader: jspb.BinaryReader): ClientAuth;
}

export namespace ClientAuth {
  export type AsObject = {
    username: string,
  }
}

export class AuthAck extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): void;

  getSessiontoken(): string;
  setSessiontoken(value: string): void;

  hasExpiry(): boolean;
  clearExpiry(): void;
  getExpiry(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setExpiry(value?: google_protobuf_timestamp_pb.Timestamp): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthAck.AsObject;
  static toObject(includeInstance: boolean, msg: AuthAck): AuthAck.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AuthAck, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthAck;
  static deserializeBinaryFromReader(message: AuthAck, reader: jspb.BinaryReader): AuthAck;
}

export namespace AuthAck {
  export type AsObject = {
    username: string,
    sessiontoken: string,
    expiry?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

