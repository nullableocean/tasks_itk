package filewords

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
)

func createTestFile(t *testing.T, filename, data string) (tmp *os.File, removeTmp func()) {
	t.Helper()

	tmp, err := os.CreateTemp("", filename)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tmp.Write([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	err = tmp.Close()
	if err != nil {
		t.Fatal(err)
	}

	remove := func() {
		os.Remove(tmp.Name())
	}

	return tmp, remove
}

func Test_WordCount(t *testing.T) {
	testCases := []struct {
		name        string
		data        string
		expectCount int
	}{
		{
			name:        "empty",
			data:        "",
			expectCount: 0,
		},
		{
			name:        "only spaces",
			data:        "   ",
			expectCount: 0,
		},
		{
			name:        "one word",
			data:        "word",
			expectCount: 1,
		},
		{
			name:        "one word with spaces",
			data:        " word ",
			expectCount: 1,
		},
		{
			name:        "many words",
			data:        strings.Repeat("word ", 1000),
			expectCount: 1000,
		},
	}

	for _, tcase := range testCases {
		t.Run(tcase.name, func(t *testing.T) {
			testFile, removeFn := createTestFile(t, "testname.txt", tcase.data)
			defer removeFn()

			count, err := countWordInFile(testFile.Name())
			if err != nil {
				t.Fatal(err)
			}

			if tcase.expectCount != count {
				t.Errorf("expected: %d, got: %d", tcase.expectCount, count)
			}
		})
	}
}

type testFileCase struct {
	testFile  *os.File
	rmFile    func()
	expect    int
	gotResult bool
}

func Test_ParallelFileSuccessHandle(t *testing.T) {
	filesNum := 10

	testCases := make(map[string]*testFileCase, filesNum)
	mu := sync.Mutex{}

	paths := make([]string, 0, filesNum)
	for i := range filesNum {
		filename := fmt.Sprintf("testname_%d.txt", i)
		newTmpFile, rmFn := createTestFile(t, filename, strings.Repeat(" word ", 100*i))

		testCases[newTmpFile.Name()] = &testFileCase{
			testFile: newTmpFile,
			rmFile:   rmFn,
			expect:   100 * i,
		}

		paths = append(paths, newTmpFile.Name())
	}

	resChannels := CountWordsInFiles(paths)
	wg := sync.WaitGroup{}
	wg.Add(len(resChannels))
	for _, resCh := range resChannels {
		go func() {
			defer wg.Done()

			res, ok := <-resCh
			if !ok {
				t.Errorf("channel was closed")
				return
			}

			if res.Err != nil {
				t.Errorf("error handle file: %v", res.Err)
				return
			}

			mu.Lock()
			defer mu.Unlock()

			fileCase, ex := testCases[res.File]
			if !ex {
				t.Errorf("incorret filename: %s not found in created cases", res.File)
				return
			}

			if fileCase.gotResult {
				t.Errorf("already result got. file: %s", res.File)
				return
			}

			if fileCase.expect != res.WordCount {
				t.Errorf("count not equal with expected. exp: %d, got: %d", fileCase.expect, res.WordCount)
			}

			fileCase.gotResult = true

			// fmt.Printf("--  FILE: %s, WORDS: %d  --\n", res.File, res.WordCount)
		}()
	}

	wg.Wait()
	for _, v := range testCases {
		v.rmFile()
	}
}
