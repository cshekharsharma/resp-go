package resp3

import (
	"bufio"
	"io"
)

// readLineCRLF reads a line from the provided bufio.Reader until a '\n' character is encountered,
// ensuring that the line ends with '\r\n'. It returns the line content without the trailing '\r\n'.
//
// This function is specifically designed to handle RESP3 protocol lines, which are expected to end with CRLF (\r\n).
// If the line does not end with '\r\n', it indicates that the message is incomplete, and the function returns
// io.ErrUnexpectedEOF to signal that more data is expected.
//
// Parameters:
//   - reader *bufio.Reader: The buffered reader from which to read the line.
//
// Returns:
//   - string: The line read from the reader, excluding the trailing '\r\n'.
//   - error: An error if the read operation fails or if the line does not end with '\r\n'.
//
// Example usage:
//
//	line, err := readLineCRLF(reader)
//	if err != nil {
//	    if err == io.ErrUnexpectedEOF {
//	        // Handle incomplete data
//	    } else {
//	        // Handle other errors
//	    }
//	}
func readLineCRLF(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return "", err
		}
		return "", err
	}
	if len(line) < 2 || line[len(line)-2] != '\r' {
		return "", io.ErrUnexpectedEOF
	}
	return line[:len(line)-2], nil
}
