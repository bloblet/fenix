import glob
import hashlib
import json
import os
import pprint
import sys


# About 4 mb
BUF_SIZE = 4194304
# Languages and extensions to check.
languages = {
    'go': ['.pb.go', '_grpc.pb.go'],
    'web': ['_pb.js', '_grpc_web_pb.js']
}

# Ensure we are in the protobuf folder
if not os.getcwd().endswith("protobufs"):
    print("Sum ran from invalid directory; please run from protobuf directory.")
    sys.exit(1)


def get_sums(proto_file):
    proto_file = os.path.split(proto_file)[1]
    name = proto_file.split('.')[0]
    sums = {}
    sums['proto'] = hash_file(f'proto/{proto_file}')
    sums['grpc_gen'] = {}
    sums['proto_gen'] = {}
    for language, extensions in languages.items():
        proto_gen_path = f'{language}/{name}{extensions[0]}'
        grpc_gen_path = f'{language}/{name}{extensions[1]}'

        if os.path.exists(grpc_gen_path):
            sums['grpc_gen'][language] = hash_file(grpc_gen_path)

        sums['proto_gen'][language] = hash_file(proto_gen_path)
    return sums


def get_files():
    protobuf_files = glob.glob('proto/*.proto')
    protobuf_sum_files = glob.glob('proto/*.proto.sum')

    if len(protobuf_sum_files) == 0:
        return (protobuf_files, protobuf_sum_files)

    files_to_remove = []
    for file in protobuf_files:
        if f'{file}.sum' not in protobuf_sum_files:
            files_to_remove.append(file)

    for file in files_to_remove:
        protobuf_files.remove(file)

    return (protobuf_files, protobuf_sum_files)


def check_states(file):
    file = os.path.split(file)[1]
    sums = get_sums(file)

    name = file.split('.')[0]

    old_sums = {}
    with open(f'proto/{file}.sum', 'r') as f:
        old_sums = json.load(f)
    
    states = {}
    proto_has_changed = old_sums['proto'] != sums['proto']

    if proto_has_changed:
        print(f'{file} has changed!')
        # other files should have changed
        for language, sum in sums['proto_gen'].items():
            proto_file = f'{language}/{name}{languages[language][0]}'
            if language in old_sums['proto_gen']:
                states[proto_file] = old_sums['proto_gen'][language] != sum
            else:
                states[proto_file] = True

        for language, sum in sums['grpc_gen'].items():
            proto_file = f'{language}/{name}{languages[language][1]}'

            if language in old_sums['grpc_gen']:
                states[proto_file] = old_sums['grpc_gen'][language] != sum
            else:
                states[proto_file] = True
    else:
        # other files should not have changed
        for language, sum in sums['proto_gen'].items():
            proto_file = f'{language}/{name}{languages[language][0]}'

            if language in old_sums['proto_gen']:
                states[proto_file] = old_sums['proto_gen'][language] == sum
            else:
                states[proto_file] = True

        for language, sum in sums['grpc_gen'].items():
            proto_file = f'{language}/{name}{languages[language][1]}'

            if language in old_sums['grpc_gen']:
                states[proto_file] = old_sums['grpc_gen'][language] == sum
            else:
                states[proto_file] = True
    return (sums, states)


def hash_file(file: str):
    sha1 = hashlib.sha1()

    with open(file, 'rb') as f:
        while True:
            data = f.read(BUF_SIZE)
            if not data:
                break
            sha1.update(data)
    return sha1.hexdigest()


def main():
    dirty = []
    proto_sums = {}

    proto_files, sum_files = get_files()
    if len(sum_files) == 0:
        print("No sum files found, generating")
        p_files = glob.glob('proto/*.proto')
        for file in p_files:
            sums = get_sums(file)
            with open(f'{file}.sum', 'w') as f:
                json.dump(sums, f, sort_keys=True, indent=4)
        print('Done! No dirty files.')
        exit(0)

    for file in proto_files:
        sums, states = check_states(file)
        proto_sums[os.path.split(file)[1]] = sums

        for file, state in states.items():
            if not state:
                dirty.append(file)

    if len(dirty) != 0:
        print("FAIL: Dirty files found.")
        pprint.pprint(dirty)
        sys.exit(1)

    # write new sums
    for file, sums, in proto_sums.items():
        with open(f'proto/{file}.sum', 'w') as f:
            json.dump(sums, f, sort_keys=True, indent=4)


if __name__ == "__main__":
    main()
