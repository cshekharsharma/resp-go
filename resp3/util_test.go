package resp3

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

func TestReadLineCRLF(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		expectErr error
	}{
		{
			name:      "Normal line ending with CRLF",
			input:     "Hello, World!\r\n",
			expected:  "Hello, World!",
			expectErr: nil,
		},
		{
			name:      "Line ending with just LF",
			input:     "Hello, World!\n",
			expected:  "",
			expectErr: io.ErrUnexpectedEOF,
		},
		{
			name:      "Line without CRLF",
			input:     "Incomplete line",
			expected:  "",
			expectErr: io.EOF,
		},
		{
			name:      "Empty line with CRLF",
			input:     "\r\n",
			expected:  "",
			expectErr: nil,
		},
		{
			name:      "Line ending with CRLF in multiple reads",
			input:     "Hello, World!",
			expected:  "",
			expectErr: io.EOF,
		},
		{
			name:      "Line with intermediate CRLF",
			input:     "Hello\r\nWorld!\r\n",
			expected:  "Hello",
			expectErr: nil,
		},
		// Additional test for the second line
		{
			name:      "Second line after intermediate CRLF",
			input:     "World!\r\n",
			expected:  "World!",
			expectErr: nil,
		},
		{
			name:      "Line with only CR",
			input:     "Hello, World!\r",
			expected:  "",
			expectErr: io.EOF,
		},
		{
			name:      "Line with CRLF split across reads",
			input:     "Hello, World!\r",
			expected:  "",
			expectErr: io.EOF,
		},
		{
			name:      "Line with extra data after CRLF",
			input:     "Hello, World!\r\nExtra data",
			expected:  "Hello, World!",
			expectErr: nil,
		},
		{
			name:      "Input ending with EOF without CRLF",
			input:     "Hello, World!",
			expected:  "",
			expectErr: io.EOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(bytes.NewBufferString(tt.input))
			result, err := readLineCRLF(reader)
			if err != tt.expectErr {
				t.Errorf("Expected error '%v', got '%v'", tt.expectErr, err)
			}
			if result != tt.expected {
				t.Errorf("Expected result '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
