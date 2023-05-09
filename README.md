# Image-Dedup-Go

* Implements image hashing to detect duplicates
* Supports:
  * Average Hash
  * Perceptual Hash
  * Difference Hash
* Images with matching hashes are then compared via SSIM to see if they're a close match

## Building

Make sure you have `go` installed. Run:

```shell
make
```

## Running
Then you can run `./imgdup`:

```shell
./imgdup --directory /path/to/images/ --hashMethod Difference
```