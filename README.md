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
  \__   Mode: byte
  \__   Logical Sector Size: 512-bytes
  \__   Entropy Window: 1 MiB (1048576-bytes)
 8.00 ┼─────╮
 7.20 ┤     ╰╮
 6.40 ┤      ╰╮
 5.60 ┤       │            ╭─╮
 4.80 ┤       ╰╮          ╭╯ ╰─╮
 4.00 ┤        ╰╮         │    ╰╮          ╭─╮
 3.20 ┤         │        ╭╯     ╰─╮       ╭╯ ╰─╮     ╭─
 2.40 ┤         ╰╮      ╭╯        ╰─╮   ╭─╯    ╰─────╯
 1.60 ┤          ╰─╮    │           ╰─╮╭╯
 0.80 ┤            ╰─╮ ╭╯             ╰╯
 0.00 ┤              ╰─╯
      └┬────┬─────┬────┬────┬─────┬────┬────┬─────┬────┬
       1    2     3    4    5     6    7    8     9   10
       Shannon Entropy Graph for test/disk.raw (X-Units in MiB)
[i] -> Total size of test/disk.raw: 10 MiB
   \__ H(X): 7.9998 bits/byte |> Sector: 0 - 2047
   \__ H(X): 7.9998 bits/byte |> Sector: 2048 - 4095
   \__ H(X): 2.0000 bits/byte |> Sector: 4096 - 6143
   \__ H(X): 0.0000 bits/byte |> Sector: 6144 - 8191
   \__ H(X): 6.0000 bits/byte |> Sector: 8192 - 10239
   \__ H(X): 3.0000 bits/byte |> Sector: 10240 - 12287
   \__ H(X): 1.0000 bits/byte |> Sector: 12288 - 14335
   \__ H(X): 4.0000 bits/byte |> Sector: 14336 - 16383
   \__ H(X): 2.0000 bits/byte |> Sector: 16384 - 18431
   \__ H(X): 3.0000 bits/byte |> Sector: 18432 - 20479
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
