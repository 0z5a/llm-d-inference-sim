/*
Copyright 2026 The llm-d-inference-sim Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tokenizer

import (
	"fmt"
	"sync"
	"testing"

	"github.com/llm-d/llm-d-inference-sim/pkg/api"
	"github.com/stretchr/testify/require"
)

func TestSimpleTokenizerRememberTokenized(t *testing.T) {
	t.Run("records paired token strings", func(t *testing.T) {
		st := NewSimpleTokenizer()
		RememberTokenized(st, &api.Tokenized{
			Tokens:  []uint32{900001, 900002},
			Strings: []string{"custom ", "answer"},
		})

		output, err := st.Detokenize([]uint32{900001, 900002})
		require.NoError(t, err)
		require.Equal(t, "custom answer", output)
	})

	t.Run("ignores unpaired token ids", func(t *testing.T) {
		st := NewSimpleTokenizer()
		RememberTokenized(st, &api.Tokenized{
			Tokens:  []uint32{900001, 900002},
			Strings: []string{"known"},
		})

		output, err := st.Detokenize([]uint32{900001, 900002})
		require.NoError(t, err)
		require.Equal(t, "known<unk_900002>", output)
	})

	t.Run("accepts a nil tokenized response", func(t *testing.T) {
		st := NewSimpleTokenizer()
		require.NotPanics(t, func() {
			RememberTokenized(st, nil)
		})
	})

	t.Run("supports concurrent registration and detokenization", func(t *testing.T) {
		const workers = 16
		const iterations = 64

		st := NewSimpleTokenizer()
		errs := make(chan error, workers)
		var wg sync.WaitGroup
		for worker := range workers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for iteration := range iterations {
					id := uint32(910000 + worker*iterations + iteration)
					want := fmt.Sprintf("token-%d-%d", worker, iteration)
					RememberTokenized(st, &api.Tokenized{
						Tokens:  []uint32{id},
						Strings: []string{want},
					})

					got, err := st.Detokenize([]uint32{id})
					if err != nil {
						errs <- err
						return
					}
					if got != want {
						errs <- fmt.Errorf("token %d: got %q, want %q", id, got, want)
						return
					}
				}
			}()
		}

		wg.Wait()
		close(errs)
		for err := range errs {
			require.NoError(t, err)
		}
	})
}
