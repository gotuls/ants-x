// MIT License

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package ants

import "time"

const defaultTimeout = 5000 * time.Millisecond

// callable return the execute result and err.
type Callable func() (interface{}, error)

type Future interface {

	// get the async result with default timeout 5s.
	Get() (interface{}, error)

	// get the async result with custom timeout.
	GetTimeout(timeout time.Duration) (interface{}, error)
}

type futureTask struct {
	callable Callable
	chn      chan interface{}
}

func newFutureTask(callable Callable) *futureTask {
	return &futureTask{
		callable: callable,
		chn:      make(chan interface{}, 1),
	}
}

func (f *futureTask) Get() (interface{}, error) {
	return f.GetTimeout(defaultTimeout)
}

func (f *futureTask) GetTimeout(timeout time.Duration) (interface{}, error) {
	select {
	case res := <-f.chn:
		if err, ok := res.(error); ok {
			return nil, err
		}
		return res, nil
	case <-time.After(timeout):
		return nil, ErrTimeout
	}
}

// execute the future task's callable.
func (f *futureTask) Run() {
	if f.callable == nil {
		f.chn <- nil
		return
	}
	r, err := f.callable()
	if err != nil {
		f.chn <- err
		return
	}
	f.chn <- r
}
