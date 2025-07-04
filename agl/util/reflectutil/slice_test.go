package reflectutil_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leclecr04/go-tool/agl/util/must"
	"github.com/leclecr04/go-tool/agl/util/reflectutil"
)

func TestSliceBuilder(t *testing.T) {
	type helloT struct {
		Hello int
	}
	const size = 10

	{
		var slice []helloT
		sb := reflectutil.NewSliceBuilder(&slice)
		for i := 0; i < size; i++ {
			elem := sb.NewElemPtr()
			err := json.Unmarshal([]byte(fmt.Sprintf(`{"Hello": %d}`, i)), elem)
			must.Must(err)
			sb.Append(elem)
		}
		assert.Equal(t, size, len(slice))
		for i, e := range slice {
			assert.Equal(t, i, e.Hello)
		}
	}

	{
		var slice []*helloT
		sb := reflectutil.NewSliceBuilder(&slice)
		for i := 0; i < size; i++ {
			elem := sb.NewElemPtr()
			err := json.Unmarshal([]byte(fmt.Sprintf(`{"Hello": %d}`, i)), elem)
			must.Must(err)
			sb.Append(elem)
		}
		assert.Equal(t, size, len(slice))
		for i, e := range slice {
			assert.Equal(t, i, e.Hello)
		}
	}

	{
		var slice []int
		sb := reflectutil.NewSliceBuilder(&slice)
		for i := 0; i < size; i++ {
			elem := sb.NewElemPtr()
			err := json.Unmarshal([]byte(fmt.Sprintf(`%d`, i)), elem)
			must.Must(err)
			sb.Append(elem)
		}
		assert.Equal(t, size, len(slice))
		for i, e := range slice {
			assert.Equal(t, i, e)
		}
	}

	{
		var slice []map[string]interface{}
		sb := reflectutil.NewSliceBuilder(&slice)
		for i := 0; i < size; i++ {
			elem := sb.NewElemPtr()
			err := json.Unmarshal([]byte(fmt.Sprintf(`{"Hello": %d}`, i)), elem)
			must.Must(err)
			sb.Append(elem)
		}
		assert.Equal(t, size, len(slice))
		for i, e := range slice {
			assert.Equal(t, float64(i), e["Hello"].(float64))
		}
	}
}
