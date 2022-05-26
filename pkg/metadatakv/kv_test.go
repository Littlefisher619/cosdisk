package kv

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"
)

func TestSerializeKey(t *testing.T) {
	testCases := []struct {
		name string
		key  Key
		want string
	}{
		{
			name: "normal",
			key: Key{
				UID:  "abc",
				Path: "file2.txt",
			},
			want: "abc\x00file2.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.key.Serialize()
			require.NoError(t, err)
			require.Equal(t, tc.want, string(actual))
		})
	}
}

func TestUnSerializeKey(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  *Key
	}{
		{
			name: "normal",
			want: &Key{
				UID:  "abc",
				Path: "file2.txt",
			},
			input: "abc\x00file2.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := &Key{}
			err := actual.UnSerialize([]byte(tc.input))
			require.NoError(t, err)
			require.Equal(t, tc.want, actual)
		})
	}
}

func TestCleanupKey(t *testing.T) {
	testCases := []struct {
		name    string
		input   Key
		want    Key
		wantErr bool
	}{
		{
			name: "normal",
			input: Key{
				UID:  "abcbc",
				Path: "/file2.txt",
			},
			want: Key{
				UID:  "abcbc",
				Path: "/file2.txt",
			},
		},
		{
			name: "multiple slash",
			input: Key{
				UID:  "abcbc",
				Path: "///file2.txt",
			},
			want: Key{
				UID:  "abcbc",
				Path: "/file2.txt",
			},
		},
		{
			name: "dot slash",
			input: Key{
				UID:  "abcbc",
				Path: "/dir/././",
			},
			want: Key{
				UID:  "abcbc",
				Path: "/dir",
			},
		},
		{
			name: "end with slash",
			input: Key{
				UID:  "abcbc",
				Path: "/dir/",
			},
			want: Key{
				UID:  "abcbc",
				Path: "/dir",
			},
		},
		{
			name: "back slash",
			input: Key{
				UID:  "abcbc",
				Path: "/\\//",
			},
			want: Key{
				UID:  "abcbc",
				Path: "/\\",
			},
		},
		{
			name: "malformed username",
			input: Key{
				UID:  "abc\x00bc",
				Path: "file2.txt",
			},
			want: Key{
				UID:  "abc\x00bc",
				Path: "file2.txt",
			},
			wantErr: true,
		},
		{
			name: "malformed path",
			input: Key{
				UID:  "abcbc",
				Path: "file2.\x00txt",
			},
			want: Key{
				UID:  "abcbc",
				Path: "file2.\x00txt",
			},
			wantErr: true,
		},
		{
			name: "path without /",
			input: Key{
				UID:  "abc",
				Path: "file2.txt",
			},
			want: Key{
				UID:  "abc",
				Path: "file2.txt",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Cleanup()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			actual := tc.input
			require.Equal(t, tc.want, actual)
		})
	}
}

type value struct {
	time.Time
}

func (v *value) UnSerialize(b []byte) error {
	newV := value{}
	err := msgpack.Unmarshal(b, &newV)
	if err != nil {
		return err
	}
	*v = newV
	return nil
}

func (v *value) Serialize() ([]byte, error) {
	return msgpack.Marshal(v)
}

func TestSerializeValue(t *testing.T) {

	testCases := []struct {
		name string
		val  value
	}{
		{
			name: "normal",
			val: value{
				Time: time.Now(),
			},
			// val: Value{
			// 	Type:  FileValueType,
			// 	FileTypeValue: &FileValue{
			// 		FileInfo: model.FileInfo{
			// 			FName:    "test",
			// 			FModTime: time.Now(),
			// 			FIsDir:   false,
			// 			FMode:    os.ModeType,
			// 		},
			// 		FileObject: model.FileObject{
			// 			Bucket: "newcontent",
			// 			Key:    "newcontent",
			// 		},
			// 	},

			// },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bytes, err := tc.val.Serialize()
			require.NoError(t, err)
			actual := &value{}
			err = actual.UnSerialize(bytes)
			require.NoError(t, err)
			bytes2, err2 := actual.Serialize()
			require.NoError(t, err2)
			require.Equal(t, bytes, bytes2)
			fmt.Println(*actual)
			fmt.Println(tc.val)

			require.Equal(t, &tc.val, actual)
		})
	}
}
