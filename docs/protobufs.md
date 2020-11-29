# How to download Fenix protobuf files:
To avoid having to download the entire fenix repo for to only use the protobuf folder, use a sparse-checkout.
Only works for git 1.7.0 (Feb 2012) or higher

```
mkdir fenix_protobufs
cd fenix_protobufs
git init
git remote add origin -f https://github.com/bloblet/fenix.git

git config core.sparseCheckout true
echo "protobufs/" >> .git/info/sparse-checkout

git pull origin master
```

# Building protobuf files on windows:
My best advice for people using windows and building protobuf files is to use WSL, and just run `bash --login gen.bash`.  I haven't had any success actually building them on linux, but that is probably just me not knowing how.  PRs are welcome.

# Recommendations for protobuf workflow
Since the generated protobuf files are checked into git, if you make a contribution to them, they need to be updated.

Fenix has a automated check for a protobuf file being generated on GitHub, but if you want to check yourself, make an executable file in the `.git/hooks` folder called `pre-push`

```bash
#!/bin/env bash
cd protobufs
python3 sum.py
```

This is the same check ran on github.

