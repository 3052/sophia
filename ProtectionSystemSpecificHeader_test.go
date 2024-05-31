package sofia

import (
   "bytes"
   "encoding/base64"
   "testing"
)

const cenc_pssh = "AAAAcHBzc2gAAAAA7e+LqXnWSs6jyCfc1R0h7QAAAFAIARIQmlNKHxLWjhojWfOHEP3bZRoFd3Vha2kiLTlhNTM0YTFmMTJkNjhlMWEyMzU5ZjM4NzEwZmRkYjY1LW1jLTAtMTQ3LTAtMEjj3JWbBg=="

func TestPssh(t *testing.T) {
   data, err := base64.StdEncoding.DecodeString(cenc_pssh)
   if err != nil {
      t.Fatal(err)
   }
   var value File
   err = value.Read(bytes.NewReader(data))
   if err != nil {
      t.Fatal(err)
   }
}
