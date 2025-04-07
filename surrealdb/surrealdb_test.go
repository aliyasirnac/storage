package surrealdb

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var testStore *Storage
var testConfig = ConfigDefault

func TestMain(m *testing.M) {
	var err error
	testStore, err = New(testConfig)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	if err := testStore.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to close store: %v\n", err)
	}
	os.Exit(code)
}

func Test_Surrealdb_Create(t *testing.T) {
	err := testStore.Set("create", []byte("test12345"), 0)
	require.NoError(t, err)
}

func Test_Surrealdb_CreateAndGet(t *testing.T) {
	err := testStore.Set("createandget", []byte("test1234"), 0)
	require.NoError(t, err)

	get, err := testStore.Get("createandget")
	require.NoError(t, err)
	require.NotEmpty(t, get)

}

func Test_Surrealdb_ListTable(t *testing.T) {
	bytes, err := testStore.List()
	require.NoError(t, err)
	require.NotEmpty(t, bytes)
}

func Test_Surrealdb_Get_WithNoErr(t *testing.T) {
	get, err := testStore.Get("create")
	require.NoError(t, err)
	require.NotEmpty(t, get)
}

func Test_Surrealdb_Delete(t *testing.T) {
	err := testStore.Set("delete", []byte("delete1234"), 0)
	require.NoError(t, err)

	err = testStore.Delete("delete")
	require.NoError(t, err)
}

func Test_Surrealdb_Flush(t *testing.T) {
	err := testStore.Reset()
	require.NoError(t, err)
}

func Test_Surrealdb_GetExpired(t *testing.T) {
	err := testStore.Set("temp", []byte("value"), 1*time.Second)
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	val, err := testStore.Get("temp")
	require.NoError(t, err)
	require.Nil(t, val)
}

func Test_Surrealdb_GetMissing(t *testing.T) {
	val, err := testStore.Get("non-existent-key")
	require.NoError(t, err)
	require.Nil(t, val)
}

func Test_Surrealdb_ListSkipsExpired(t *testing.T) {
	_ = testStore.Reset()

	// Geçerli kayıt
	_ = testStore.Set("valid", []byte("123"), 0)

	// Süresi geçen kayıt
	_ = testStore.Set("expired", []byte("456"), 1*time.Second)
	time.Sleep(2 * time.Second)

	// List çağrısı
	data, err := testStore.List()
	require.NoError(t, err)

	var result map[string][]byte
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	require.Contains(t, result, "valid")
	require.NotContains(t, result, "expired")
}
func BenchmarkSet(b *testing.B) {
	store, err := New(ConfigDefault)
	if err != nil {
		b.Fatalf("failed to init storage: %v", err)
	}
	defer store.Close()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench-key-%d-%d", i, time.Now().UnixNano())
		err := store.Set(fmt.Sprintf("bench-key-%d", key), []byte("value"), 0)
		if err != nil {
			b.Errorf("Set failed: %v", err)
		}
	}
}
