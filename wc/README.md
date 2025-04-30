## Build Code 
```bash 
make build
```
This will build binary package at `/binaries/wc`. 

## Usage 

### Count Lines 
```bash 
./binaries/wc -l testFiles/test1.txt testFiles/test2.txt
```
**Sample Output :**
```bash
       0 testFiles/test2.txt
      10 testFiles/test1.txt
      10 total
```

###  Without Any Flags (Counts lines, words, and bytes)
```bash 
./binaries/wc testFiles/test1.txt testF
iles/test2.txt
```
**Sample Output :**
```bash
       0       0       0 testFiles/test2.txt
      10      78     445 testFiles/test1.txt
      10      78     445 total
```
###  Read from STDIN

```bash 
./binaries/wc
```
***Past the following input:***
```txt
Hello world
This is a test.
Go is awesome!
```
***Press Ctrl+D (Unix/macOS) or Ctrl+Z + Enter (Windows) to end input.***
**Expected Output:**
```bash
       3       9      43
```
### Help Flag 
To display the help message and see all available options, use the `-h` flag:

```bash
./binaries/wc -h
```
**Sample Output:**
```bash 
Usage of ./binaries/wc:
  -c    count bytes
  -l    count lines
  -w    count words
```