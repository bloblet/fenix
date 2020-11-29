// package: 
// file: user.proto

import * as jspb from "google-protobuf";

export class Nil extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Nil.AsObject;
  static toObject(includeInstance: boolean, msg: Nil): Nil.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Nil, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Nil;
  static deserializeBinaryFromReader(message: Nil, reader: jspb.BinaryReader): Nil;
}

export namespace Nil {
  export type AsObject = {
  }
}

export class User extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getUsername(): string;
  setUsername(value: string): void;

  getDiscriminator(): string;
  setDiscriminator(value: string): void;

  hasActivity(): boolean;
  clearActivity(): void;
  getActivity(): Activity | undefined;
  setActivity(value?: Activity): void;

  clearServersList(): void;
  getServersList(): Array<string>;
  setServersList(value: Array<string>): void;
  addServers(value: string, index?: number): string;

  clearFriendsList(): void;
  getFriendsList(): Array<string>;
  setFriendsList(value: Array<string>): void;
  addFriends(value: string, index?: number): string;

  getToken(): string;
  setToken(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getSalt(): Uint8Array | string;
  getSalt_asU8(): Uint8Array;
  getSalt_asB64(): string;
  setSalt(value: Uint8Array | string): void;

  getPassword(): Uint8Array | string;
  getPassword_asU8(): Uint8Array;
  getPassword_asB64(): string;
  setPassword(value: Uint8Array | string): void;

  hasSettings(): boolean;
  clearSettings(): void;
  getSettings(): UserSettings | undefined;
  setSettings(value?: UserSettings): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): User.AsObject;
  static toObject(includeInstance: boolean, msg: User): User.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: User, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): User;
  static deserializeBinaryFromReader(message: User, reader: jspb.BinaryReader): User;
}

export namespace User {
  export type AsObject = {
    id: string,
    username: string,
    discriminator: string,
    activity?: Activity.AsObject,
    serversList: Array<string>,
    friendsList: Array<string>,
    token: string,
    email: string,
    salt: Uint8Array | string,
    password: Uint8Array | string,
    settings?: UserSettings.AsObject,
  }
}

export class UserSettings extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserSettings.AsObject;
  static toObject(includeInstance: boolean, msg: UserSettings): UserSettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UserSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserSettings;
  static deserializeBinaryFromReader(message: UserSettings, reader: jspb.BinaryReader): UserSettings;
}

export namespace UserSettings {
  export type AsObject = {
  }
}

export class Activity extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Activity.AsObject;
  static toObject(includeInstance: boolean, msg: Activity): Activity.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Activity, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Activity;
  static deserializeBinaryFromReader(message: Activity, reader: jspb.BinaryReader): Activity;
}

export namespace Activity {
  export type AsObject = {
  }
}

export class Authenticate extends jspb.Message {
  getToken(): string;
  setToken(value: string): void;

  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Authenticate.AsObject;
  static toObject(includeInstance: boolean, msg: Authenticate): Authenticate.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Authenticate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Authenticate;
  static deserializeBinaryFromReader(message: Authenticate, reader: jspb.BinaryReader): Authenticate;
}

export namespace Authenticate {
  export type AsObject = {
    token: string,
    id: string,
  }
}

export class CreateUser extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateUser.AsObject;
  static toObject(includeInstance: boolean, msg: CreateUser): CreateUser.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CreateUser, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateUser;
  static deserializeBinaryFromReader(message: CreateUser, reader: jspb.BinaryReader): CreateUser;
}

export namespace CreateUser {
  export type AsObject = {
    username: string,
    password: string,
    email: string,
  }
}

