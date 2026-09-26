# h8

`h8` is a small CLI tool for calculating and visualizing Shannon entropy across files and raw disk images.

It is designed for quick inspection of data regions that may be encrypted, compressed, sparse, or otherwise unusual.

## Features

- Shannon entropy analysis
- Byte and bit entropy modes
- Configurable entropy window size
- Logical Block Address (LBA) output
- Terminal-based entropy graph
- File and raw disk image support
- Recursive directory scanning

## Example

```bash
h8 disk.raw
```

Example output:

```text
[i] -> Settings:
      \____ Mode: byte
      \____ Logical Sector Size: 512 bytes
      \____ Entropy Window: 1 MiB

LBA: 0    - 2047  | H(X): 7.9998 bits/byte
LBA: 2048 - 4095  | H(X): 7.9998 bits/byte
LBA: 4096 - 6143  | H(X): 2.0000 bits/byte
```

## Usage

```bash
h8 [path] [flags]
```

Useful flags include:

```text
--byte              use byte entropy
--bit               use bit entropy
-s, --size          entropy window size in bytes
--sector-size       logical sector size in bytes
```

## Why entropy?

High Shannon entropy can indicate encrypted, compressed, or randomized data, while low entropy often corresponds to structured, repetitive, or empty regions.

Entropy alone cannot determine whether data is encrypted.

## Use cases

`h8` can be useful for:

- ransomware recovery analysis
- digital forensics
- disk image inspection
- locating unusual regions in large files
- comparing entropy across storage ranges

## Disclaimer

Entropy is a statistical indicator, not proof of encryption or compromise. Results should be interpreted alongside filesystem, forensic, and recovery context.

<br>

---
