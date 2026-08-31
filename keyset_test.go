package collections

import (
	"bytes"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sortTimesAsc(times []time.Time) []time.Time {
	sortedTimes := append([]time.Time(nil), times...)
	sort.Slice(sortedTimes, func(i, j int) bool {
		return sortedTimes[i].Before(sortedTimes[j])
	})
	return sortedTimes
}

func sortTimesDesc(times []time.Time) []time.Time {
	sortedTimes := append([]time.Time(nil), times...)
	sort.Slice(sortedTimes, func(i, j int) bool {
		return sortedTimes[i].After(sortedTimes[j])
	})
	return sortedTimes
}

func TestKeySet(t *testing.T) {
	sk, ctx, _ := deps()
	keyset := NewKeySet[string](sk, 0, StringKeyEncoder)

	// test insert and get
	key := "hi"
	keyset.Insert(ctx, key)
	require.True(t, keyset.Has(ctx, key))

	// test delete and get error
	keyset.Delete(ctx, key)
	require.False(t, keyset.Has(ctx, key))
}

func TestKeySet_Iterate(t *testing.T) {
	sk, ctx, _ := deps()
	keyset := NewKeySet[string](sk, 0, StringKeyEncoder)
	keyset.Insert(ctx, "a")
	keyset.Insert(ctx, "aa")
	keyset.Insert(ctx, "b")
	keyset.Insert(ctx, "bb")

	expectedKeys := []string{"a", "aa", "b", "bb"}

	iter := keyset.Iterate(ctx, Range[string]{})
	defer iter.Close()
	for i, o := range iter.Keys() {
		require.Equal(t, expectedKeys[i], o)
	}
}

func TestKeysetIterator(t *testing.T) {
	sk, ctx, _ := deps()

	keyset := NewKeySet[string](sk, 0, StringKeyEncoder)
	keyset.Insert(ctx, "a")

	iter := keyset.Iterate(ctx, Range[string]{})
	defer iter.Close()

	assert.True(t, iter.Valid())
	assert.EqualValues(t, "a", iter.Key())

	iter.Next()
	assert.False(t, iter.Valid())
}

func TestTimeKeySet(t *testing.T) {
	sk, ctx, _ := deps()
	keyset := NewKeySet[time.Time](sk, 0, TimeKeyEncoder)

	// Use a fixed time
	now := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	keyset.Insert(ctx, now)
	require.True(t, keyset.Has(ctx, now))

	// Test delete and get
	keyset.Delete(ctx, now)
	require.False(t, keyset.Has(ctx, now))
}

func TestTimeKeySet_IterateAscending(t *testing.T) {
	sk, ctx, _ := deps()
	keyset := NewKeySet[time.Time](sk, 0, TimeKeyEncoder)

	now := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	times := []time.Time{
		now.Add(2 * time.Second),
		now.Add(1 * time.Second),
		now.Add(3 * time.Second),
		now,
	}

	// Insert times into the keyset
	for _, t := range times {
		keyset.Insert(ctx, t)
	}

	times = sortTimesAsc(times)

	// Iterate over the keyset in ascending order
	iter := keyset.Iterate(ctx, Range[time.Time]{})
	defer iter.Close()

	keys := iter.Keys()
	require.Equal(t, len(times), len(keys))

	for i, k := range keys {
		// Strip monotonic clock readings
		expectedTime := times[i].Round(0)
		actualTime := k.Round(0)

		// Compare UnixNano timestamps
		require.Equal(t, expectedTime.UnixNano(), actualTime.UnixNano())
	}
}

func TestTimeKeySet_IterateDescending(t *testing.T) {
	sk, ctx, _ := deps()
	keyset := NewKeySet[time.Time](sk, 0, TimeKeyEncoder)

	now := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	times := []time.Time{
		now.Add(2 * time.Second),
		now.Add(1 * time.Second),
		now.Add(3 * time.Second),
		now,
	}

	// Insert times into the keyset
	for _, t := range times {
		keyset.Insert(ctx, t)
	}

	times = sortTimesDesc(times)

	// Iterate over the keyset in descending order
	iter := keyset.Iterate(ctx, Range[time.Time]{}.Descending())
	defer iter.Close()

	keys := iter.Keys()
	require.Equal(t, len(times), len(keys))

	for i, k := range keys {
		expectedTime := times[i].Round(0)
		actualTime := k.Round(0)
		require.Equal(t, expectedTime.UnixNano(), actualTime.UnixNano())
	}
}

func TestTimeKeyEncoder_EncodeDecode(t *testing.T) {
	now := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	encoded := TimeKeyEncoder.Encode(now)
	_, decoded := TimeKeyEncoder.Decode(encoded)

	// Compare UnixNano timestamps
	require.Equal(t, now.UnixNano(), decoded.UnixNano())
}

func TestTimeKeyEncoder_BoundsAndOrdering(t *testing.T) {
	times := []time.Time{
		minTimeKey,
		time.Unix(-1, 999999999).UTC(),
		time.Unix(0, 0).UTC(),
		maxTimeKey,
	}

	for i, timestamp := range times {
		encoded := TimeKeyEncoder.Encode(timestamp)
		consumed, decoded := TimeKeyEncoder.Decode(encoded)
		require.Equal(t, 8, consumed)
		require.True(t, timestamp.Equal(decoded))
		if i > 0 {
			require.Less(t, bytes.Compare(TimeKeyEncoder.Encode(times[i-1]), encoded), 0)
		}
	}
}

func TestTimeKeyEncoder_RejectsOutOfRangeTimes(t *testing.T) {
	require.Panics(t, func() {
		TimeKeyEncoder.Encode(minTimeKey.Add(-time.Nanosecond))
	})
	require.Panics(t, func() {
		TimeKeyEncoder.Encode(maxTimeKey.Add(time.Nanosecond))
	})
}

func TestTimeKeySet_OrderConsistency(t *testing.T) {
	sk, ctx, _ := deps()
	keyset := NewKeySet[time.Time](sk, 0, TimeKeyEncoder)

	now := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	times := []time.Time{
		now.Add(-1 * time.Hour),
		now,
		now.Add(1 * time.Hour),
		now.Add(2 * time.Hour),
		now.Add(-2 * time.Hour),
	}

	// Insert times into the keyset
	for _, t := range times {
		keyset.Insert(ctx, t)
	}

	sortedTimesAsc := sortTimesAsc(times)

	// Iterate over the keyset in ascending order
	iterAsc := keyset.Iterate(ctx, Range[time.Time]{})
	defer iterAsc.Close()

	keysAsc := iterAsc.Keys()
	require.Equal(t, len(times), len(keysAsc))

	for i, k := range keysAsc {
		expectedTime := sortedTimesAsc[i].Round(0)
		actualTime := k.Round(0)
		require.Equal(t, expectedTime.UnixNano(), actualTime.UnixNano())
	}

	sortedTimesDesc := sortTimesDesc(times)

	// Iterate over the keyset in descending order
	iterDesc := keyset.Iterate(ctx, Range[time.Time]{}.Descending())
	defer iterDesc.Close()

	keysDesc := iterDesc.Keys()
	require.Equal(t, len(times), len(keysDesc))

	for i, k := range keysDesc {
		expectedTime := sortedTimesDesc[i].Round(0)
		actualTime := k.Round(0)
		require.Equal(t, expectedTime.UnixNano(), actualTime.UnixNano())
	}
}

func TestTimeKeySet_IterateRange(t *testing.T) {
	sk, ctx, _ := deps()
	keyset := NewKeySet[time.Time](sk, 0, TimeKeyEncoder)

	now := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	times := []time.Time{
		now.Add(1 * time.Second),
		now.Add(2 * time.Second),
		now.Add(3 * time.Second),
		now.Add(4 * time.Second),
		now.Add(5 * time.Second),
	}

	// Insert times into the keyset
	for _, t := range times {
		keyset.Insert(ctx, t)
	}

	// Define range from now.Add(2s) inclusive to now.Add(4s) exclusive
	iter := keyset.Iterate(ctx, Range[time.Time]{}.
		StartInclusive(now.Add(2*time.Second)).
		EndExclusive(now.Add(4*time.Second)))
	defer iter.Close()

	expectedTimes := []time.Time{
		now.Add(2 * time.Second),
		now.Add(3 * time.Second),
	}

	keys := iter.Keys()
	require.Equal(t, len(expectedTimes), len(keys))

	for i, k := range keys {
		expectedTime := expectedTimes[i].Round(0)
		actualTime := k.Round(0)
		require.Equal(t, expectedTime.UnixNano(), actualTime.UnixNano())
	}
}

func TestTimeKeySet_SameTimeKeys(t *testing.T) {
	sk, ctx, _ := deps()
	keyset := NewKeySet[time.Time](sk, 0, TimeKeyEncoder)

	now := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Insert the same time multiple times (should only be stored once in a set)
	keyset.Insert(ctx, now)
	keyset.Insert(ctx, now)
	keyset.Insert(ctx, now)

	iter := keyset.Iterate(ctx, Range[time.Time]{})
	defer iter.Close()

	keys := iter.Keys()
	require.Equal(t, 1, len(keys))

	expectedTime := now.Round(0)
	actualTime := keys[0].Round(0)
	require.Equal(t, expectedTime.UnixNano(), actualTime.UnixNano())
}
