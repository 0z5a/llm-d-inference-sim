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

package simulator

import (
	"github.com/llm-d/llm-d-inference-sim/pkg/api"
	"github.com/llm-d/llm-d-inference-sim/pkg/common"
	"github.com/llm-d/llm-d-inference-sim/pkg/tokenizer"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type fixedResponseDataset struct {
	response *api.Tokenized
	err      error
}

func (d *fixedResponseDataset) Close() error { return nil }

func (d *fixedResponseDataset) GetResponseTokens(api.Request) (*api.Tokenized, string, error) {
	return d.response, common.StopFinishReason, d.err
}

var _ = Describe("response token mappings", func() {
	It("remembers dataset-provided strings for subsequent derendering", func() {
		response := &api.Tokenized{
			Tokens:  []uint32{900001, 900002},
			Strings: []string{"custom ", "answer"},
		}
		ctx := &SimContext{
			dataset:   &fixedResponseDataset{response: response},
			Tokenizer: tokenizer.NewSimpleTokenizer(),
		}

		got, finishReason, err := ctx.GetResponseTokens(&api.TextCompletionsRequest{})
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(response))
		Expect(finishReason).To(Equal(common.StopFinishReason))

		text, err := ctx.Tokenizer.Detokenize(got.Tokens)
		Expect(err).NotTo(HaveOccurred())
		Expect(text).To(Equal("custom answer"))
	})
})
