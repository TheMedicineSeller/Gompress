<h3>A text file compression tool developed from scratch using golang.</h3>

Two executables are included with this project: 
- gompress which runs the compression algorithm on the text file provided via a cmd line option and
- degompress that takes in the compressed file path (.enc file) and recovers it as filename_extracted.txt.

The algorithm uses Huffman trees to generate a lower length bit representation for each existing token in the input file (vocabulary). In this case I have tried to specifically handle tokens as 32 bit unicode characters instead of ASCII chars thanks to runes inbuilt in golang.
